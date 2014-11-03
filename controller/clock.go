package controller

import (
	"github.com/ninjasphere/astrotime"
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
	Dawn(ref time.Time) time.Time                                          // Next civil dawn after the specified time
	Dusk(ref time.Time) time.Time                                          // Next civil dusk after the specified time.
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
	if !clk.useCoordinates {
		sunset := time.Date(ref.Year(), ref.Month(), ref.Day(), 18, 0, 0, 0, ref.Location())
		if sunset.Sub(ref) < 0 {
			sunset = sunset.AddDate(0, 0, 1)
		}
		return sunset
	}
	return astrotime.NextSunset(ref, clk.latitude, clk.longtitude)
}

func (clk *systemclock) Sunrise(ref time.Time) time.Time {
	if !clk.useCoordinates {
		sunrise := time.Date(ref.Year(), ref.Month(), ref.Day(), 6, 0, 0, 0, ref.Location())
		if sunrise.Sub(ref) < 0 {
			sunrise = sunrise.AddDate(0, 0, 1)
		}
		return sunrise
	}
	return astrotime.NextSunrise(ref, clk.latitude, clk.longtitude)
}

func (clk *systemclock) Dawn(ref time.Time) time.Time {
	if !clk.useCoordinates {
		dawn := time.Date(ref.Year(), ref.Month(), ref.Day(), 18, 0, 0, 0, ref.Location())
		if dawn.Sub(ref) < 0 {
			dawn = dawn.AddDate(0, 0, 1)
		}
		return dawn
	}
	return astrotime.NextDawn(ref, clk.latitude, clk.longtitude, astrotime.CIVIL_DAWN)
}

func (clk *systemclock) Dusk(ref time.Time) time.Time {
	if !clk.useCoordinates {
		dusk := time.Date(ref.Year(), ref.Month(), ref.Day(), 6, 0, 0, 0, ref.Location())
		if dusk.Sub(ref) < 0 {
			dusk = dusk.AddDate(0, 0, 1)
		}
		return dusk
	}
	return astrotime.NextDusk(ref, clk.latitude, clk.longtitude, astrotime.CIVIL_DUSK)
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
