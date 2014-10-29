package controller

import (
	"time"
)

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
