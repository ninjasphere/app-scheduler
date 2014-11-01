package controller

import (
	// "github.com/ninjasphere/app-scheduler/model"
	"testing"
	"time"
)

func TestTaskRespectsClosedWindows(t *testing.T) {
	initMockClock(testTime)
	task := &task{}
	actuations := make(chan actuationRequest, 2)
	if err := task.init(taskWithEarlierTimeOfDayOpenDelayCloseWindow, actuations); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	go func() {
		task.loop()
	}()
	time.Sleep(time.Millisecond * time.Duration(500))
	select {
	case actuation := <-actuations:
		actuation.reply <- nil
		t.Fatalf("unexpected actuation")
	default:
	}
	task.quit <- struct{}{}
}

func TestTaskRespectsPermanentlyClosedWindows(t *testing.T) {
	initMockClock(testTime)
	task := &task{}
	actuations := make(chan actuationRequest, 2)
	if err := task.init(taskWithEarlierTimestampOpenDelayCloseWindow, actuations); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	done := make(chan struct{})
	go func() {
		task.loop()
		done <- struct{}{}
	}()
	time.Sleep(time.Millisecond * time.Duration(500))
	select {
	case actuation := <-actuations:
		actuation.reply <- nil
		t.Fatalf("unexpected actuation")
	default:
	}
	select {
	case done := <-done:
		_ = done
	default:
		t.Fatalf("expected task to exit")
	}
}
