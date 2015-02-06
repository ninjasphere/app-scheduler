package controller

import (
	"fmt"
	"github.com/ninjasphere/app-scheduler/model"
	"time"
)

const (
	statusScheduled         = "scheduled"
	statusWaitingForOpen    = "waiting for open - %s"
	statusExecutingOpen     = "executing open actions"
	statusWaitingForClose   = "waiting for close - %s"
	statusExecutingClose    = "executing close actions"
	statusPermanentlyClosed = "permanently closed"
)

type task struct {
	model      *model.Task
	window     *window
	openers    []*action
	closers    []*action
	quit       chan struct{}
	actuations chan actuationRequest
	status     string
}

func (t *task) init(m *model.Task, actuations chan actuationRequest) error {
	t.model = m
	t.window = &window{}
	t.status = statusScheduled
	t.actuations = actuations
	err := t.window.init(m.Window)
	if err != nil {
		return err
	}
	t.openers = make([]*action, 0, 0)
	t.closers = make([]*action, 0, 0)
	for _, a := range m.Open {
		if actor, err := newAction(a); err == nil {
			t.openers = append(t.openers, actor)
		} else {
			return err
		}
	}
	for _, a := range m.Close {
		if actor, err := newAction(a); err == nil {
			t.closers = append(t.closers, actor)
		} else {
			return err
		}
	}

	t.quit = make(chan struct{}, 1)
	return nil
}

// A scheduler task
func (t *task) loop() bool {
	for {

		// FIXME: if the window is not recurrent, then we need to check that it is still valid.

		now := clock.Now()
		scheduledAt := now

		t.status = fmt.Sprintf(statusWaitingForOpen, t.window.after)

		for {
			if t.window.isPermanentlyClosed(now) {
				log.Debugf("At '%v' the window '%v' for task '%s' became permanently closed. The task will exit and then be cancelled.", now, t.window, t.model.ID)
				// stop running when we can run no more
				t.status = statusPermanentlyClosed
				return true
			}

			if t.window.isOpen(scheduledAt, now) {
				break
			} else {
				var quit bool
				if t.window.after.hasEventOccurred(scheduledAt, now) {
					t.window.after.setLastFired(now)
				}
				if t.window.before.hasEventOccurred(scheduledAt, now) {
					t.window.before.setLastFired(now)
				}
				quit, now = t.waitForOpenEvent(scheduledAt)
				if quit {
					return false
				}
			}
		}

		t.window.after.setLastFired(now)
		t.status = statusExecutingOpen

		t.doActions("open", t.openers)

		t.status = fmt.Sprintf(statusWaitingForClose, t.window.before)

		if t.waitForCloseEvent(now) {
			return false
		}

		t.window.before.setLastFired(now)
		t.status = statusExecutingClose

		t.doActions("close", t.closers)

	}
}

// wait for an open event or until I am told to quit
func (t *task) waitForOpenEvent(ref time.Time) (bool, time.Time) {

	done := t.window.whenOpen(ref)

	select {
	case quitSignal := <-t.quit:
		_ = quitSignal
		return true, clock.Now()
	case openSignal := <-done:
		return false, openSignal
	}
}

// wait for a close event or until I am told to quit
func (t *task) waitForCloseEvent(opened time.Time) bool {

	done := t.window.whenClosed(opened)

	select {
	case quitSignal := <-t.quit:
		_ = quitSignal
		return true
	case closeSignal := <-done:
		_ = closeSignal
		return false
	}
}

func (t *task) doActions(phase string, actions []*action) {
	reply := make(chan error)
	for _, o := range actions {
		t.actuations <- actuationRequest{
			action: o,
			reply:  reply,
		}
		err := <-reply
		if err != nil {
			log.Errorf("The '%s' action during '%s' event for task '%s' failed with error: %s", o.model.Action, phase, t.model.ID, err)
		}
	}
}
