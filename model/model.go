// Package model describes a collection of tasks to be executed at the opening and closing
// of schedule windows.
package model

// A ThingAction does something to a single thing.
type ThingAction struct {
	ActionType string `json:"type,omitempty"`
	ThingId    string `json:"thing-id,omitempty"`
	Action     string `json:"action,omitempty"`
}

// An Event specifies a point in time and is used to delimit the open
// and close of a task window.
//
// Rule determines how Param is interpreted. The following rules are supported:
//
// timestamp - param is an absolute timestamp of the form "yyyy-mm-dd HH:MM:SS"
//
// time-of-day - specifies a clock time between 00:00:00 and 23:59:59
//
// delay - specifies an event that occurs HH:MM:SS after the start of the scheduler
// or the last Until event (if this is a 'From' event) or the last From event (if this is an
// 'Until' event).
//
// sunrise - specifies an event that occurs at the local sunrise
//
// sunset - speciies an event that occurs at the local sunset
//
// Other events that might be used in future (but not currently supported) are:
//
// once - occurs once after the scheduler starts and never again
//
// never - an event that never occurs
//
type Event struct {
	Rule  string `json:"rule,omitempty"`
	Param string `json:"param,omitempty"`
}

// A Window describes a period of time during which a Scheduler Task runs.
// The From event specifies when the Window starts (usually a time of day). The Until event
// specifies when the Window closes (usually another time of day or a delay).
type Window struct {
	From  *Event `json:"from,omitempty"`
	Until *Event `json:"until,omitempty"`
}

// A Task runs in a Window and performs the Open actions, waits for the Until event to occur
// then runs the Close actions. If the Window is a recurring window, then the cycle repeats, otherwise the Task ends.
type Task struct {
	Uuid        string         `json:"uuid,omitempty"`
	Description string         `json:"description,omitempty"`
	Window      *Window        `json:"window,omitempty"`
	Open        []*ThingAction `json:"open,omitempty"`
	Close       []*ThingAction `json:"close,omitempty"`
}

// A Schedule specifies a list of Tasks.
type Schedule struct {
	Version string  `json:"version,omitempty"`
	Tasks   []*Task `json:"schedule,omitempty"`
}
