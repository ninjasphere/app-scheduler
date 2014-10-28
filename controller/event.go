package controller

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ninjasphere/app-scheduler/model"
)

type event struct {
	// given a reference time, yields a timestamp relative to the reference time
	asTimestamp func(reference time.Time) time.Time
	// the timestamp parsed from the time specification
	parsed *time.Time
	// whether the parsed event is a clock time (true) or an absolute (local) timestamp
	clockTime bool
}

// Initialize the event from the specification.
func (t *event) init(m *model.Event) error {
	var err error

	switch m.Rule {
	case "timestamp":
		parsed, err := time.Parse("20060102T150405", m.Param)
		if err == nil {
			t.asTimestamp = t.timestamp
			t.parsed = &parsed
		}
		fallthrough
	case "time-of-day":
		parsed, err := time.Parse("15:04:05", m.Param)
		if err == nil {
			t.asTimestamp = t.timeofday
			t.parsed = &parsed
			t.clockTime = true
		}
		fallthrough
	case "delay":
		parsed, err := time.Parse("15:04:05", m.Param)
		if err == nil {
			t.asTimestamp = t.timeofday
			t.parsed = &parsed
		}
		fallthrough
	case "dusk":
		parsed, err := time.Parse("15:04:05", "18:00:00")
		if err == nil {
			t.asTimestamp = t.dusk
			t.parsed = &parsed
			t.clockTime = true
		}
		fallthrough
	case "dawn":
		parsed, err := time.Parse("15:04:05", "06:00:00")
		if err == nil {
			t.asTimestamp = t.dawn
			t.parsed = &parsed
			t.clockTime = true
		}
		fallthrough
	default:
		json, _ := json.Marshal(m)
		return fmt.Errorf("bad time specification: '%s'", json)
	}
	return err
}

// Return the parsed timestamp.
func (t *event) timestamp(ref time.Time) time.Time {
	return *t.parsed
}

// Return the specified time of day, relative to the reference timestamp.
func (t *event) timeofday(ref time.Time) time.Time {
	return time.Date(ref.Year(), ref.Month(), ref.Day(), (*t.parsed).Hour(), (*t.parsed).Minute(), (*t.parsed).Second(), 0, nil)
}

// Answer the timestamp after the delay specfied.
func (t *event) delay(ref time.Time) time.Time {
	delay := time.Duration((*t.parsed).Hour())*time.Hour + time.Duration((*t.parsed).Minute())*time.Minute + time.Duration((*t.parsed).Second())*time.Second
	return ref.Add(delay)
}

// Answer the time of the next dusk in the current location.
func (t *event) dusk(ref time.Time) time.Time {
	//FIXME: use location data, if available, to calculate sunset
	return t.timeofday(ref)
}

// Answer the time of the next dawn in the current location.
func (t *event) dawn(ref time.Time) time.Time {
	//FIXME: use location data, if available, to calculate sunrise
	return t.timeofday(ref)
}
