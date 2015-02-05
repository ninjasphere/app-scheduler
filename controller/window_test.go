package controller

import (
	"github.com/ninjasphere/app-scheduler/model"
	"testing"
	"time"
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

func Test20150205(t *testing.T) {

	scheduledAt := time.Date(2015, 02, 05, 12, 47, 30, 337083531, time.Now().Location())
	// exactTime := time.Date(2015, 02, 05, 12, 48, 24, 0, time.Now().Location())
	ref := time.Date(2015, 02, 05, 12, 48, 24, 4585963, time.Now().Location())

	testWindow := &model.Window{
		After: &model.Event{
			Rule:  "time-of-day",
			Param: "12:48:24",
		},
		Before: &model.Event{
			Rule:  "delay",
			Param: "00:01:00",
		},
	}

	mockClock := initMockClock(scheduledAt, defaultJitter)
	stub := &window{}
	stub.init(testWindow)

	if stub.isPermanentlyClosed(scheduledAt) {
		t.Fatalf("was %v, wanted %v", true, false)
	}

	if stub.isOpen(scheduledAt, scheduledAt) {
		t.Fatalf("was %v, wanted %v", true, false)
	}

	wakeup := make(chan time.Time)

	go func() {
		done := stub.whenOpen(scheduledAt)

		select {
		case openSignal := <-done:
			wakeup <- openSignal
		}
	}()

	mockClock.SetNow(ref)

	now := <-wakeup

	if stub.isPermanentlyClosed(now) {
		t.Fatalf("was %v, wanted %v", true, false)
	}

	if !stub.isOpen(scheduledAt, now) {
		t.Fatalf("was %v, wanted %v", false, true)
	}

}
