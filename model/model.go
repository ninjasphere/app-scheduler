// Package model describes a collection of tasks to be executed at the opening and closing
// of schedule windows.
package model

// A Action does something to a single thing.
type Action struct {
	ActionType string `json:"type,omitempty"`
	ThingID    string `json:"thingID,omitempty"`
	Action     string `json:"action,omitempty"`
	SubjectID  string `json:"subject,omitempty"`
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
// or the last Before event (if this is a 'After' event) or the last After event (if this is an
// 'Before' event).
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
// The After event specifies when the Window starts (usually a time of day). The Before event
// specifies when the Window closes (usually another time of day or a delay).
type Window struct {
	After  *Event `json:"after,omitempty"`
	Before *Event `json:"before,omitempty"`
}

// A Task waits until the After event occurs (unless it has already occurred) then performs the Open actions, waits for the Before event to occur
// then runs the Close actions. If the Window is a recurring window, then the cycle repeats, otherwise the Task ends.
type Task struct {
	ID          string    `json:"id,omitempty"`
	Tags        []string  `json:"tags,omitempty"`
	Description string    `json:"description,omitempty"`
	Window      *Window   `json:"window,omitempty"`
	Open        []*Action `json:"open,omitempty"`
	Close       []*Action `json:"close,omitempty"`
}

// A Location describes a particular geographical location. Used as input to resolution of "sunrise" and "sunset" event rules.
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Altitude  float64 `json:"altitude,omitempty"`
}

// A Schedule specifies a list of Tasks, a Location and a TimeZone
type Schedule struct {
	Version  string    `json:"version,omitempty"`
	TimeZone string    `json:"timezone"` // used during resolution of 'time-of-day' and 'timestamp' rules
	Location *Location `json:"location,omitempty"`
	Tasks    []*Task   `json:"schedule,omitempty"`
}
