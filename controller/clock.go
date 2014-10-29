package controller

import (
	"time"
)

type Clock interface {
	Now() time.Time
	AfterFunc(delay time.Duration, then func())
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
