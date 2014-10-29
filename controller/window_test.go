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

func TestBadInitFromNil(t *testing.T) {
	stub := &window{}
	err := stub.init(&model.Window{
		From:  nil,
		Until: afterNowTimestampModel,
	})
	if err == nil {
		t.Fatalf("expected error on init.from == nil")
	}
}

func TestBadInitUntilNil(t *testing.T) {
	stub := &window{}
	err := stub.init(&model.Window{
		Until: nil,
		From:  afterNowTimestampModel,
	})
	if err == nil {
		t.Fatalf("expected error on init.until == nil")
	}
}

func runNonOverlappingWindow(t *testing.T, m *model.Window, permanentlyClosed bool) {
	initMockClock(testTime)
	stub := &window{}
	err := stub.init(m)
	if err != nil {
		t.Fatalf("unexpected error while opening window: %v", err)
	}
	if stub.isOpen(testTime) {
		t.Fatalf("window should not be open now")
	}
	if permanentlyClosed != stub.isPermanentlyClosed(testTime) {
		t.Fatalf("expecting permanentlyClosed == %v, but found opposite", permanentlyClosed)
	}
}

func runOverlappingWindow(t *testing.T, m *model.Window) {
	initMockClock(testTime)
	stub := &window{}
	err := stub.init(m)
	if err != nil {
		t.Fatalf("unexpected error while opening window: %v", err)
	}
	if !stub.isOpen(testTime) {
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
