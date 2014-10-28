package controller

import (
	"fmt"
	"strings"
	"time"
)

type timespec struct {
	// given a reference time, yields a timestamp relative to the reference time
	asTimestamp func(reference time.Time) time.Time
	// the timestamp parsed from the time specification
	parsed *time.Time
	// whether the parsed timespec is a clock time (true) or an absolute (local) timestamp
	clockTime bool
}

// Initialize the timespec from the specification.
func (t *timespec) init(spec string) error {
	var err error
	words := strings.Split(spec, " ")
	if len(words) > 0 {
		var arg string
		if len(words) > 1 {
			arg = words[1]
		} else {
			arg = ""
		}

		switch words[0] {
		case "timestamp":
			if len(words) > 1 {
				parsed, err := time.Parse("20060102T150405", arg)
				if err == nil {
					t.asTimestamp = t.timestamp
					t.parsed = &parsed
				}
			}
			fallthrough
		case "time-of-day":
			if len(words) > 1 {
				parsed, err := time.Parse("15:04:05", arg)
				if err == nil {
					t.asTimestamp = t.timeofday
					t.parsed = &parsed
					t.clockTime = true
				}
			}
			fallthrough
		case "delay":
			if len(words) > 1 {
				parsed, err := time.Parse("15:04:05", arg)
				if err == nil {
					t.asTimestamp = t.timeofday
					t.parsed = &parsed
				}
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
			return fmt.Errorf("bad time specification: '%s'", spec)
		}
	}
	return err
}

// Return the parsed timestamp.
func (t *timespec) timestamp(ref time.Time) time.Time {
	return *t.parsed
}

// Return the specified time of day, relative to the reference timestamp.
func (t *timespec) timeofday(ref time.Time) time.Time {
	return time.Date(ref.Year(), ref.Month(), ref.Day(), (*t.parsed).Hour(), (*t.parsed).Minute(), (*t.parsed).Second(), 0, nil)
}

// Answer the timestamp after the delay specfied.
func (t *timespec) delay(ref time.Time) time.Time {
	delay := time.Duration((*t.parsed).Hour())*time.Hour + time.Duration((*t.parsed).Minute())*time.Minute + time.Duration((*t.parsed).Second())*time.Second
	return ref.Add(delay)
}

// Answer the time of the next dusk in the current location.
func (t *timespec) dusk(ref time.Time) time.Time {
	//FIXME: use location data, if available, to calculate sunset
	return t.timeofday(ref)
}

// Answer the time of the next dawn in the current location.
func (t *timespec) dawn(ref time.Time) time.Time {
	//FIXME: use location data, if available, to calculate sunrise
	return t.timeofday(ref)
}
