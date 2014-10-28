// Package model describes a collection of tasks to be executed at the opening and closing
// of schedule windows.
package model

// ActionSpec describes an action to be taken.
type ThingAction struct {
	ActionType string `json:"type"`
	ThingId    string `json:"thing-id"`
	Action     string `json:"action"`
}

// WindowSpec describes a period of time.
type Window struct {
	From  string `json:"from"`
	Until string `json:"until"`
}

// ItemSpec describes the set of actions to be taken when a window opens
// and when a window closes.
type Task struct {
	Uuid        string         `json:"uuid"`
	Description string         `json:"description"`
	Window      *Window        `json:"window"`
	Open        []*ThingAction `json:"open"`
	Close       []*ThingAction `json:"close"`
}

// A Schedule specifies a list items that the scheduler
type Schedule struct {
	version string  `json:"version"`
	Tasks   []*Task `json:"schedule"`
}
