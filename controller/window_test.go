package controller

import (
	"github.com/ninjasphere/app-scheduler/model"
	"testing"
)

func TestBadInitNil(t *testing.T) {
	stub := &window{}
	err := stub.init(nil)
	if err == nil {
		t.Fatalf("expected error on init(nil)")
	}
}

func TestBadInitAfterNil(t *testing.T) {
	stub := &window{}
	err := stub.init(&model.Window{
		After:  nil,
		Before: afterNowTimestampModel,
	})
	if err == nil {
		t.Fatalf("expected error on init.from == nil")
	}
}

func TestBadInitBeforeNil(t *testing.T) {
	stub := &window{}
	err := stub.init(&model.Window{
		Before: nil,
		After:  afterNowTimestampModel,
	})
	if err == nil {
		t.Fatalf("expected error on init.until == nil")
	}
}

func runNonOverlappingWindow(t *testing.T, m *model.Window, permanentlyClosed bool) {
	initMockClock(testTime, defaultJitter)
	stub := &window{}
	err := stub.init(m)
	if err != nil {
		t.Fatalf("unexpected error while opening window: %v", err)
	}
	if stub.isOpen(testTime, testTime) {
		t.Fatalf("window should not be open now")
	}
	if permanentlyClosed != stub.isPermanentlyClosed(testTime) {
		t.Fatalf("expecting permanentlyClosed == %v, but found opposite %v", permanentlyClosed, stub)
	}
}

func runOverlappingWindow(t *testing.T, m *model.Window) {
	initMockClock(testTime, defaultJitter)
	stub := &window{}
	err := stub.init(m)
	if err != nil {
		t.Fatalf("unexpected error while opening window: %v", err)
	}
	if !stub.isOpen(testTime, testTime) {
		t.Fatalf("window should be open now")
	}
	if stub.isPermanentlyClosed(testTime) {
		t.Fatalf("expecting permanentlyClosed == false, but found opposite")
	}
}

func TestEarlierTimeOfDayWindowIsNotOpen(t *testing.T) {
	runNonOverlappingWindow(t, earlierTimeOfDayWindow, false)
}

func TestLaterTimeOfDayWindowIsNotOpen(t *testing.T) {
	runNonOverlappingWindow(t, laterTimeOfDayWindow, false)
}

func TestOverlappingTimeOfDayWindowIsNotOpen(t *testing.T) {
	runOverlappingWindow(t, overlappingTimeOfDayWindow)
}

func TestEarlierTimestampWindowIsNotOpen(t *testing.T) {
	runNonOverlappingWindow(t, earlierTimestampWindow, true)
}

func TestLaterTimestampWindowIsNotOpen(t *testing.T) {
	runNonOverlappingWindow(t, laterTimestampWindow, false)
}

func TestOverlappingTimestampWindowIsNotOpen(t *testing.T) {
	runOverlappingWindow(t, overlappingTimestampWindow)
}

func TestOverlappingTimestampOpenDelayCloseWindowIsNotOpen(t *testing.T) {
	runOverlappingWindow(t, overlappingTimestampOpenDelayCloseWindow)
}

func TestEarlierTimestampOpenDelayCloseWindowIsNotOpen(t *testing.T) {
	runNonOverlappingWindow(t, earlierTimestampOpenDelayCloseWindow, true)
}

func TestLaterTimestampOpenDelayCloseWindowIsNotOpen(t *testing.T) {
	runNonOverlappingWindow(t, laterTimestampOpenDelayCloseWindow, false)
}

func TestSunriseSunsetWindow(t *testing.T) {
	runNonOverlappingWindow(t, sunriseSunsetWindow, false)
}

func TestNowNeverWindow(t *testing.T) {
	initMockClock(testTime, defaultJitter)
	stub := &window{}
	if err := stub.init(nowNeverWindow); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !stub.isOpen(testTime, testTime) {
		t.Fatalf("nowNeverWindow.isOpen() was %v, wanted %v", false, true)
	}
	if stub.isPermanentlyClosed(testTime) {
		t.Fatalf("nowNeverWindow.isPermanentlyClosed() was %v, wanted %v", true, false)
	}
}

func TestNeverNowWindow(t *testing.T) {
	initMockClock(testTime, defaultJitter)
	stub := &window{}
	if err := stub.init(neverNowWindow); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stub.isOpen(testTime, testTime) {
		t.Fatalf("neverNowWindow.isOpen() was %v, wanted %v", true, false)
	}
	if stub.isPermanentlyClosed(testTime) {
		t.Fatalf("neverNowWindow.isPermanentlyClosed() was %v, wanted %v", true, false)
	}
}
