package controller

import (
	"time"
)

// Clock provides an alternative interface to the system clock. Used to allow testing of time related functions.
type Clock interface {
	Now() time.Time
	AfterFunc(delay time.Duration, then func())
	Location() *time.Location
}

type systemclock struct {
}

type callback func()

var clock Clock = &systemclock{}

func (*systemclock) Now() time.Time {
	return time.Now()
}

func (*systemclock) AfterFunc(delay time.Duration, then func()) {
	time.AfterFunc(delay, then)
}

func (*systemclock) Location() *time.Location {
	return time.Now().Location()
}
