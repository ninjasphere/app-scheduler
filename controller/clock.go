package controller

import (
	"time"
)

// Clock provides an alternative interface to the system clock. Used to allow testing of time related functions.
type Clock interface {
	Now() time.Time                                                        // The current time.
	AfterFunc(delay time.Duration, then func())                            // Do the specified thing after the specified time
	Location() *time.Location                                              // Get the time zone location
	SetLocation(location *time.Location)                                   // Set the time zone location
	Sunset(ref time.Time) time.Time                                        // Next sunset after the specified time
	Sunrise(ref time.Time) time.Time                                       // Next sunrise after the specified time.
	SetCoordinates(latitude float64, longtitude float64, altitude float64) // Set the coordinates to be used for sunrise/sunset calculations
	ResetCoordinates()                                                     // Clear the coordinates.
}

type systemclock struct {
	location       *time.Location
	useCoordinates bool
	latitude       float64
	longtitude     float64
	altitude       float64
}

type callback func()

var clock Clock = &systemclock{
	location:       time.Now().Location(),
	useCoordinates: false,
}

func (*systemclock) Now() time.Time {
	return time.Now()
}

func (*systemclock) AfterFunc(delay time.Duration, then func()) {
	time.AfterFunc(delay, then)
}

func (clk *systemclock) Location() *time.Location {
	return clk.location
}

func (clk *systemclock) SetLocation(l *time.Location) {
	if l == nil {
		clk.location = time.Now().Location()
	} else {
		clk.location = l
	}
}

func (clk *systemclock) Sunset(ref time.Time) time.Time {
	if true || !clk.useCoordinates {
		sunset := time.Date(ref.Year(), ref.Month(), ref.Day(), 18, 0, 0, 0, ref.Location())
		if sunset.Sub(ref) < 0 {
			sunset = sunset.AddDate(0, 0, 1)
		}
		return sunset
	}
	// FIXME: support use of the coordinates
	return ref
}

func (clk *systemclock) Sunrise(ref time.Time) time.Time {
	if true || !clk.useCoordinates {
		sunrise := time.Date(ref.Year(), ref.Month(), ref.Day(), 6, 0, 0, 0, ref.Location())
		if sunrise.Sub(ref) < 0 {
			sunrise = sunrise.AddDate(0, 0, 1)
		}
		return sunrise
	}
	// FIXME: support use of the coordinates
	return ref
}

func (clk *systemclock) SetCoordinates(latitude float64, longtitude float64, altitude float64) {
	clk.useCoordinates = true
	clk.latitude = latitude
	clk.longtitude = longtitude
	clk.altitude = altitude
}

func (clk *systemclock) ResetCoordinates() {
	clk.useCoordinates = false
	clk.latitude = 0
	clk.longtitude = 0
	clk.altitude = 0
}
