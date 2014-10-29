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

func TestEarlierWindowIsNotOpen(t *testing.T) {
	initMockClock(testTime)
	stub := &window{}
	err := stub.init(earlierWindow)
	if err != nil {
		t.Fatalf("unexpected error while opening window: %v", err)
	}
	if stub.isOpen(testTime) {
		t.Fatalf("window should not be open now")
	}
}

func TestOverlappingWindowIsNotOpen(t *testing.T) {
	initMockClock(testTime)
	stub := &window{}
	err := stub.init(overlappingWindow)
	if err != nil {
		t.Fatalf("unexpected error while opening window: %v", err)
	}
	if !stub.isOpen(testTime) {
		t.Fatalf("window should be open now")
	}
}
