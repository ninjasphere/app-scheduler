// Package model describes a collection of tasks to be executed at the opening and closing
// of schedule windows.
package model

// ActionSpec describes an action to be taken.
type ThingAction struct {
	ActionType string `json:"type,omitempty"`
	ThingId    string `json:"thing-id,omitempty"`
	Action     string `json:"action,omitempty"`
}

type Event struct {
	Rule  string `json:"rule,omitempty"`
	Param string `json:"param,omitempty"`
}

// WindowSpec describes a period of time.
type Window struct {
	From  *Event `json:"from,omitempty"`
	Until *Event `json:"until,omitempty"`
}

// ItemSpec describes the set of actions to be taken when a window opens
// and when a window closes.
type Task struct {
	Uuid        string         `json:"uuid,omitempty"`
	Description string         `json:"description,omitempty"`
	Window      *Window        `json:"window,omitempty"`
	Open        []*ThingAction `json:"open,omitempty"`
	Close       []*ThingAction `json:"close,omitempty"`
}

// A Schedule specifies a list items that the scheduler
type Schedule struct {
	version string  `json:"version,omitempty"`
	Tasks   []*Task `json:"schedule,omitempty"`
}
