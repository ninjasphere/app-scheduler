package controller

import (
	"github.com/ninjasphere/app-scheduler/model"

	"testing"
	"time"
)

var (
	testTime                = time.Date(2014, 10, 29, 11, 22, 30, 0, time.Now().Location())
	futureTime              = time.Date(2014, 10, 29, 12, 00, 00, 0, time.Now().Location())
	beforeNowTimeOfDayModel = &model.Event{
		Rule:  "time-of-day",
		Param: "09:00:00",
	}
	afterNowTimeOfDayModel = &model.Event{
		Rule:  "time-of-day",
		Param: "12:00:00",
	}
	beforeNowTimestampModel = &model.Event{
		Rule:  "timestamp",
		Param: "2014-10-29 09:00:00",
	}
	afterNowTimestampModel = &model.Event{
		Rule:  "timestamp",
		Param: "2014-10-29 12:00:00",
	}
	shortDelayModel = &model.Event{
		Rule:  "delay",
		Param: "00:15:00",
	}
	delayModel = &model.Event{
		Rule:  "delay",
		Param: "00:45:00",
	}
	sunsetModel = &model.Event{
		Rule: "sunset",
	}
	sunriseModel = &model.Event{
		Rule: "sunrise",
	}
	bogusModel = &model.Event{
		Rule: "bogus",
	}
	bogusTimestamp = &model.Event{
		Rule:  "timestamp",
		Param: "bogus",
	}
)

func assertFired(t *testing.T, whenDone chan time.Time) {
	select {
	case tmp := <-whenDone:
		_ = tmp
	default:
		t.Fatalf("at %v event should have fired, but did not.", clock.Now())
	}
}

func assertNotFired(t *testing.T, whenDone chan time.Time) {
	select {
	case tmp := <-whenDone:
		t.Fatalf("at %v event should not have fired, but did %v", clock.Now(), tmp)
	default:
	}
}

func runBogus(t *testing.T, e *model.Event) {
	initMockClock(testTime)
	event, err := newEvent(e, false)
	if err == nil {
		t.Fatalf("expecting error but none found for %+v", *e)
	}
	if event != nil {
		t.Fatalf("expecting nil event, but found 1")
	}
}

func runBeforeNow(t *testing.T, e *model.Event, close bool, clockTime bool) {
	initMockClock(testTime)
	event, err := newEvent(e, close)
	if err != nil {
		t.Fatalf("unexpected error on newEvent %s", err)
	}
	if !event.hasTimestamp() {
		t.Fatalf("time of day event should have timestamp")
	}

	if !close {
		if event.asTimestamp(testTime).Sub(testTime) > 0 {
			t.Fatalf("test event (%v) is after the test time (%v)", event.asTimestamp(testTime), testTime)
		}
		whenDone := event.waiter(testTime)
		assertFired(t, whenDone)
	} else {
		diff := event.asTimestamp(testTime).Sub(testTime)
		if diff < 0 && clockTime {
			t.Fatalf("test event (%v) is before the test time (%v)", event.asTimestamp(testTime), testTime)
		}
		if diff > 0 && !clockTime {
			t.Fatalf("test event (%v) is after the test time (%v)", event.asTimestamp(testTime), testTime)
		}
	}
}

func runAfterNow(t *testing.T, e *model.Event, close bool, shouldFireInFuture bool) {
	mock := initMockClock(testTime)
	event, err := newEvent(e, close)
	if err != nil {
		t.Fatalf("unexpected error on newEvent %s", err)
	}
	if !event.hasTimestamp() {
		t.Fatalf("time of day event should have timestamp")
	}
	if event.asTimestamp(testTime).Sub(testTime) < 0 {
		t.Fatalf("test event (%v) is before the test time (%v)", event.asTimestamp(testTime), testTime)
	}
	whenDone := event.waiter(testTime)
	assertNotFired(t, whenDone)
	mock.SetNow(futureTime)
	if shouldFireInFuture {
		assertFired(t, whenDone)
	} else {
		assertNotFired(t, whenDone)
	}
}

func TestOpenTimeOfDayBeforeNow(t *testing.T) {
	runBeforeNow(t, beforeNowTimeOfDayModel, false, true)
}

func TestOpenTimeOfDayAfterNow(t *testing.T) {
	runAfterNow(t, afterNowTimeOfDayModel, false, true)
}

func TestCloseTimeOfDayBeforeNow(t *testing.T) {
	runBeforeNow(t, beforeNowTimeOfDayModel, true, true)
}

func TestCloseTimeOfDayAfterNow(t *testing.T) {
	runAfterNow(t, afterNowTimeOfDayModel, true, true)
}

func TestOpenTimestampBeforeNow(t *testing.T) {
	runBeforeNow(t, beforeNowTimestampModel, false, false)
}

func TestOpenTimestampAfterNow(t *testing.T) {
	runAfterNow(t, afterNowTimestampModel, false, true)
}

func TestCloseTimestampBeforeNow(t *testing.T) {
	runBeforeNow(t, beforeNowTimestampModel, true, false)
}

func TestCloseTimestampAfterNow(t *testing.T) {
	runAfterNow(t, afterNowTimestampModel, true, true)
}

func TestOpenDelayAfterNow(t *testing.T) {
	runAfterNow(t, delayModel, false, false)
}

func TestCloseDelayAfterNow(t *testing.T) {
	runAfterNow(t, delayModel, true, false)
}

func TestOpenShortDelayAfterNow(t *testing.T) {
	runAfterNow(t, shortDelayModel, false, true)
}

func TestCloseShortDelayAfterNow(t *testing.T) {
	runAfterNow(t, shortDelayModel, true, true)
}
func TestOpenSunsetAfterNow(t *testing.T) {
	runAfterNow(t, sunsetModel, false, false)
}
func TestCloseSunsetAfterNow(t *testing.T) {
	runAfterNow(t, sunsetModel, true, false)
}
func TestOpenSunriseBeforeNow(t *testing.T) {
	runBeforeNow(t, sunriseModel, false, true)
}
func TestCloseSunriseAfterNow(t *testing.T) {
	runAfterNow(t, sunriseModel, true, false)
}
func TestBogusModel(t *testing.T) {
	runBogus(t, bogusModel)
}
func TestBogusTimestamp(t *testing.T) {
	runBogus(t, bogusTimestamp)
}
