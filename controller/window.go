package controller

import (
	"fmt"
	"time"

	"github.com/ninjasphere/app-scheduler/model"
)

type window struct {
	after  Event
	before Event
}

func (w *window) init(m *model.Window) error {
	var err error
	if m == nil {
		return fmt.Errorf("illegal argument: m == nil")
	}
	w.after, err = newEvent(m.After, false)
	if err == nil {
		w.before, err = newEvent(m.Before, true)
	}
	return err
}

// Answer true if the window is permanently closed with respect to the specified time.
// That is: the close event is a non-recurring event has already occurred.
//
// This method is never true for recurring close events provided the open
// event has not yet occurred or itself recurring
func (w *window) isPermanentlyClosed(ref time.Time) bool {

	return !w.isOpen(ref, ref) &&
		(w.after.hasFinalEventOccurred(ref) || (w.after.hasEventOccurred(ref, ref) && w.before.hasFinalEventOccurred(ref)))
}

// Answer true if the window is open with respect to the specified time.
func (w *window) isOpen(scheduledAt time.Time, ref time.Time) bool {
	openWaitsForTime := w.after.hasTimestamp()

	if openWaitsForTime {

		// when both events are timestamp based, check
		// that the reference timestamp is within the boundaries
		// of those timestamp

		return w.after.hasEventOccurred(scheduledAt, ref) &&
			!w.before.hasEventOccurred(w.after.asTimestamp(scheduledAt), ref)
	} else {

		// when neither events are timestamp based, we have to
		// wait to wait for the open event to know we are open

		return w.after.hasEventOccurred(scheduledAt, ref) &&
			!w.before.hasEventOccurred(scheduledAt, ref)
	}
}

// Answer a channel that will receive an event when the next open event occurs.
func (w *window) whenOpen(ref time.Time) chan time.Time {
	return w.after.waiter(ref)
}

// Answer a channel that will receive an event when the next close event after the specified open event occurs.
func (w *window) whenClosed(opened time.Time) chan time.Time {
	return w.before.waiter(opened)
}

func (w *window) StringAt(ref time.Time) string {
	return fmt.Sprintf("window[after: %s, before:%s]", dump(w.after, ref), dump(w.before, ref))
}

func (w *window) String() string {
	return fmt.Sprintf("[%s, %s]", w.after.String(), w.before.String())
}
