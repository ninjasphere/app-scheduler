package controller

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/ninjasphere/app-scheduler/model"
)

//Event objects know how to wait for some kind of event to occur. Most time-based
//events can answer the timestamp when they are next expected to occur, relative to a
//specified reference time. Some events are recurring, some are not. Non-recurring events
//may be final.
type Event interface {
	String() string                      // A description of the event.
	hasTimestamp() bool                  // Time based events have timestamps, other types of events do not
	asTimestamp(ref time.Time) time.Time // Answers the timestamp of time-based event. Panics otherwise.
	waiter(ref time.Time) chan time.Time // Answer a channel for the event which receives a timestamp when the event occurs
	isRecurring() bool                   // Answer true if the event is a recurring event but false if the event can only happen once
	hasEventOccurred(                    // Answer true if the event scheduled at scheduledAt has occurred by ref
		scheduledAt time.Time,
		ref time.Time) bool
	// Answer true if the final event of this type has occurred. Not true for
	// recurring events or for non-recurring events whose timestamp is less than the reference timestamp.
	hasFinalEventOccurred(ref time.Time) bool
}

type timeEvent struct {
	// the event model
	model *model.Event
	// true when the event has fired
	lastFired time.Time
	// the timestamp parsed from the time specification
	parsed *time.Time
	// the timestamp for an event that occurs near this time.
	polyAsTimestamp func(ref time.Time) time.Time
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

	if m == nil {
		return nil, fmt.Errorf("illegal argument: model == nil")
	}

	switch m.Rule {
	case "timestamp":
		parsed, err = time.ParseInLocation("2006-01-02 15:04:05", m.Param, clock.Location())
		if err == nil {
			result := &timestamp{
				timeEvent: timeEvent{
					model:  m,
					parsed: &parsed,
				},
			}
			result.timeEvent.polyAsTimestamp = result.asTimestamp
			return result, nil
		}
	case "time-of-day":
		parsed, err = time.Parse("15:04:05", m.Param)
		if err == nil {
			result := &timeOfDay{
				timeEvent: timeEvent{
					model:  m,
					parsed: &parsed,
				},
				closeEvent: closeEvent,
			}
			result.timeEvent.polyAsTimestamp = result.asTimestamp
			return result, nil
		}
	case "delay":
		parsed, err = time.Parse("15:04:05", m.Param)
		if err == nil {
			result := &delay{
				timeEvent: timeEvent{
					model:  m,
					parsed: &parsed,
				},
			}
			result.timeEvent.polyAsTimestamp = result.asTimestamp
			return result, nil
		}
	case "sunset":
		parsed, err = time.Parse("15:04:05", "18:00:00")
		if err == nil {
			result := &sunset{
				timeOfDay: timeOfDay{
					timeEvent: timeEvent{
						model:  m,
						parsed: &parsed,
					},
					closeEvent: closeEvent,
				},
			}
			result.timeEvent.polyAsTimestamp = result.asTimestamp
			return result, nil
		}
	case "sunrise":
		parsed, err = time.Parse("15:04:05", "06:00:00")
		if err == nil {
			result := &sunrise{
				timeOfDay: timeOfDay{
					timeEvent: timeEvent{
						model:  m,
						parsed: &parsed,
					},
					closeEvent: closeEvent,
				},
			}
			result.timeEvent.polyAsTimestamp = result.asTimestamp
			return result, nil
		}
	default:
		json, _ := json.Marshal(m)
		return nil, fmt.Errorf("bad time specification: '%s'", json)
	}

	return nil, err
}

func (t *timeEvent) String() string {
	return fmt.Sprintf("%s %s", t.model.Rule, t.model.Param)
}

func (t *timeEvent) hasTimestamp() bool {
	return true
}

func (t *timeEvent) isRecurring() bool {
	return true
}

func (t *timeEvent) hasEventOccurred(scheduledAt time.Time, ref time.Time) bool {
	return t.polyAsTimestamp(scheduledAt).Sub(ref) <= 0
}

func (t *timeEvent) hasFinalEventOccurred(ref time.Time) bool {
	return false
}

func (t *timestamp) isRecurring() bool {
	return false
}

func (t *timestamp) hasFinalEventOccurred(ref time.Time) bool {
	return t.polyAsTimestamp(ref).Sub(ref) <= 0
}

func (t *timeEvent) waiter(ref time.Time) chan time.Time {
	now := clock.Now()
	when := t.polyAsTimestamp(ref)
	delay := when.Sub(now)
	waiter := make(chan time.Time, 1)
	if when.Sub(t.lastFired) >= 0 {
		if delay > 0 {
			log.Debugf("waiting now (%v) for an event at (%v)", clock.Now(), when)
			clock.AfterFunc(delay, func() {
				t.lastFired = clock.Now()
				waiter <- t.lastFired
			})
		} else {
			t.lastFired = now
			waiter <- now
			log.Debugf("waiter fired event because time already passed %v", clock.Now())
		}
	} else {
		log.Warningf("this waiter will block forever because the event @ %v was already fired at %v", delay, t.lastFired)
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
	for t.lastFired.Sub(tmp) >= 0 {
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

func dump(e Event, ref time.Time) string {
	dump := fmt.Sprintf("event @ %v, type=%s", ref, reflect.ValueOf(e).Type())
	if e.hasTimestamp() {
		dump = fmt.Sprintf("%s, asTimestamp(.)=%v", dump, e.asTimestamp(ref))
	}
	dump = fmt.Sprintf("%s, isRecurring=%v, hasEventOccurred(., .)=%v, hasFinalEventOccurred(.)=%v", dump, e.isRecurring(), e.hasEventOccurred(ref, ref), e.hasFinalEventOccurred(ref))
	return dump
}
