package controller

import (
	"fmt"

	"github.com/ninjasphere/app-scheduler/model"
	"time"
)

type task struct {
	model  *model.Task
	window *window
	quit   chan struct{}
}

func (t *task) init(m *model.Task) error {
	t.model = m
	t.window = &window{}
	err := t.window.init(m.Window)
	if err != nil {
		return err
	}
	t.quit = make(chan struct{}, 1)
	return nil
}

// A scheduler task
func (t *task) loop() {
	for {

		// FIXME: if the window is not recurrent, then we need to check that it is still valid.

		var openedAt time.Time
		now := clock.Now()

		if t.window.isPermanentlyClosed(now) {
			log.Debugf("at '%v' the window '%v' for task '%s' became permanently closed. the task will exit.", now, t.window, t.model.Uuid)
			// stop running when we can run no more
			return
		}

		if !t.window.isOpen(now) {
			var quit bool
			quit, openedAt = t.waitForOpenEvent(now)
			if quit {
				return
			}
		} else {
			openedAt = now
		}

		t.doOpenActions()

		if t.waitForCloseEvent(openedAt) {
			return
		}

		t.doCloseActions()

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

func (t *task) doOpenActions() {
	fmt.Printf("TBD: open actions to be performed now\n")
	//
}

func (t *task) doCloseActions() {
	fmt.Printf("TBD: close actions to be performed now\n")
}
