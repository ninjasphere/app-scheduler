package controller

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ninjasphere/app-scheduler/model"
)

type Event interface {
	// Time based events have timestamps, other types of events do not
	hasTimestamp() bool
	// Answers the timestamp of time-based event. Panics otherwise.
	asTimestamp(ref time.Time) time.Time
	// Answer a channel for the event which receives a timestamp when the event occurs
	waiter(ref time.Time) chan time.Time
}

type timeEvent struct {
	// the timestamp parsed from the time specification
	parsed *time.Time
	// the timestamp for an event that occurs near this time.
	asTimestamp func(ref time.Time) time.Time
}

// An event that occurs at a specified timestamp
type timestamp struct {
	timeEvent
}

// An event that occurs after a delay from the reference timestamp.
type delay struct {
	timeEvent
}

// An event that occurs, every day, at a specified time of day
type timeOfDay struct {
	timeEvent
	closeEvent bool
}

// Sunset each day.
type sunset struct {
	timeOfDay
}

// Sunrise each day.
type sunrise struct {
	timeOfDay
}

// Initialize the event from the specification.
func newEvent(m *model.Event, closeEvent bool) (Event, error) {
	var parsed time.Time
	var err error

	switch m.Rule {
	case "timestamp":
		parsed, err = time.ParseInLocation("2006-01-02 15:04:05", m.Param, clock.Location())
		if err == nil {
			result := &timestamp{
				timeEvent: timeEvent{
					parsed: &parsed,
				},
			}
			result.timeEvent.asTimestamp = result.asTimestamp
			return result, nil
		}
	case "time-of-day":
		parsed, err = time.Parse("15:04:05", m.Param)
		if err == nil {
			result := &timeOfDay{
				timeEvent: timeEvent{
					parsed: &parsed,
				},
				closeEvent: closeEvent,
			}
			result.timeEvent.asTimestamp = result.asTimestamp
			return result, nil
		}
	case "delay":
		parsed, err = time.Parse("15:04:05", m.Param)
		if err == nil {
			result := &delay{
				timeEvent: timeEvent{
					parsed: &parsed,
				},
			}
			result.timeEvent.asTimestamp = result.asTimestamp
			return result, nil
		}
	case "sunset":
		parsed, err = time.Parse("15:04:05", "18:00:00")
		if err == nil {
			result := &sunset{
				timeOfDay: timeOfDay{
					timeEvent: timeEvent{
						parsed: &parsed,
					},
					closeEvent: closeEvent,
				},
			}
			result.timeEvent.asTimestamp = result.asTimestamp
			return result, nil
		}
	case "sunrise":
		parsed, err = time.Parse("15:04:05", "06:00:00")
		if err == nil {
			result := &sunrise{
				timeOfDay: timeOfDay{
					timeEvent: timeEvent{
						parsed: &parsed,
					},
					closeEvent: closeEvent,
				},
			}
			result.timeEvent.asTimestamp = result.asTimestamp
			return result, nil
		}
	default:
		json, _ := json.Marshal(m)
		return nil, fmt.Errorf("bad time specification: '%s'", json)
	}

	return nil, err
}

func (t *timeEvent) hasTimestamp() bool {
	return true
}

func (t *timeEvent) waiter(ref time.Time) chan time.Time {
	now := clock.Now()
	delay := t.asTimestamp(ref).Sub(now)
	waiter := make(chan time.Time, 1)
	if delay > 0 {
		clock.AfterFunc(delay, func() {
			waiter <- clock.Now()
		})
	} else {
		waiter <- now
		log.Debugf("waiter fired event because time already passed %v", clock.Now())
	}
	return waiter
}

// Return the parsed timestamp.
func (t *timestamp) asTimestamp(ref time.Time) time.Time {
	return *t.parsed
}

// Return the specified time of day, relative to the reference timestamp.
func (t *timeOfDay) asTimestamp(ref time.Time) time.Time {
	tmp := time.Date(ref.Year(), ref.Month(), ref.Day(), (*t.parsed).Hour(), (*t.parsed).Minute(), (*t.parsed).Second(), 0, clock.Location())
	if tmp.Sub(ref) < 0 && t.closeEvent {
		tmp = tmp.AddDate(0, 0, 1)
	}
	return tmp
}

// Answer the timestamp after the delay specfied.
func (t *delay) asTimestamp(ref time.Time) time.Time {
	delay := time.Duration((*t.parsed).Hour())*time.Hour + time.Duration((*t.parsed).Minute())*time.Minute + time.Duration((*t.parsed).Second())*time.Second
	return ref.Add(delay)
}

// Answer the time of the next sunset in the current location.
func (t *sunset) asTimestamp(ref time.Time) time.Time {
	return t.timeOfDay.asTimestamp(ref)
}

// Answer the time of the next sunrise in the current location.
func (t *sunrise) asTimestamp(ref time.Time) time.Time {
	//FIXME: use location data, if available, to calculate sunrise
	return t.timeOfDay.asTimestamp(ref)
}
