package controller

import (
	"time"
)

// Clock provides an alternative interface to the system clock. Used to allow testing of time related functions.
type Clock interface {
	Now() time.Time
	AfterFunc(delay time.Duration, then func())
	Location() *time.Location
	SetLocation(location *time.Location)
}

type systemclock struct {
	location *time.Location
}

type callback func()

var clock Clock = &systemclock{
	location: time.Now().Location(),
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
