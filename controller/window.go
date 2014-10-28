package controller

import (
	"time"

	"github.com/ninjasphere/app-scheduler/model"
)

type window struct {
	model *model.Window
	from  event
	until event
}

func (w *window) init(m *model.Window) error {
	w.model = m
	err := w.from.init(m.From)
	if err != nil {
		err = w.until.init(m.Until)
	}
	return err
}

// Answer true if the window is open with respect to the specified time.
func (w *window) isOpen(ref time.Time) bool {
	openTimestamp := w.openTimestamp(ref)
	closeTimestamp := w.closeTimestamp(openTimestamp)

	return openTimestamp.Sub(ref) < 0 &&
		ref.Sub(closeTimestamp) < 0 &&
		openTimestamp.Sub(closeTimestamp) > 0
}

// Answer timestamp of the current (or next) open event, given the current timestamp
func (w *window) openTimestamp(ref time.Time) time.Time {
	return w.from.asTimestamp(ref)
}

// Answer the timestamp of the next close event, given an open event with the specified timestamp.
func (w *window) closeTimestamp(ref time.Time) time.Time {
	openTimestamp := w.openTimestamp(ref)
	closeTimestamp := w.until.asTimestamp(openTimestamp)
	if closeTimestamp.Sub(openTimestamp) < 0 {
		if w.from.clockTime && w.until.clockTime {
			// when open and close times are specified with clock times, then
			// we must cope with open times that start before midnight and end after midnight
			closeTimestamp = closeTimestamp.AddDate(0, 0, 1)
		} else {
			log.Fatalf("confusing window specification - what should I do here? (from,until) == (%s, %s)", w.model.From, w.model.Until)
		}
	}
	return closeTimestamp
}

// Answer a channel that will signal when the specified deadline time has been reached.
func getWaiter(deadline time.Time) chan time.Time {
	now := time.Now()
	delay := deadline.Sub(now)
	waiter := make(chan time.Time, 1)
	if delay > 0 {
		time.AfterFunc(delay, func() {
			waiter <- time.Now()
		})
	} else {
		waiter <- now
	}
	return waiter
}

// Answer a channel that will receive an event when the next open event occurs.
func (w *window) whenOpen(ref time.Time) chan time.Time {
	return getWaiter(w.openTimestamp(ref))
}

// Answer a channel that will receive an event when the next close event after the specified open event occurs.
func (w *window) whenClosed(opened time.Time) chan time.Time {
	return getWaiter(w.closeTimestamp(opened))
}
