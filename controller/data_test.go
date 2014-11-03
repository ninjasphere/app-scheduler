package controller

import (
	"github.com/ninjasphere/app-scheduler/model"
	"time"
)

var (
	muchEarlierTime     = time.Date(2014, 10, 29, 0, 3, 00, 0, time.Now().Location())
	earlierTime         = time.Date(2014, 10, 29, 0, 6, 00, 0, time.Now().Location())
	testTime            = time.Date(2014, 10, 29, 11, 22, 30, 0, time.Now().Location())
	futureTime          = time.Date(2014, 10, 29, 12, 00, 00, 0, time.Now().Location())
	futureTimeDelta1    = time.Date(2014, 10, 29, 12, 00, 00, 1, time.Now().Location())
	futureTimeDeltaNeg1 = time.Date(2014, 10, 29, 11, 59, 59, 999999999, time.Now().Location())

	muchEarlierTimeOfDayModel = &model.Event{
		Rule:  "time-of-day",
		Param: "03:00:00",
	}
	beforeNowTimeOfDayModel = &model.Event{
		Rule:  "time-of-day",
		Param: "09:00:00",
	}
	earlierTimeOfDayWindow = &model.Window{
		After:  muchEarlierTimeOfDayModel,
		Before: beforeNowTimeOfDayModel,
	}
	overlappingTimeOfDayWindow = &model.Window{
		After:  beforeNowTimeOfDayModel,
		Before: afterNowTimeOfDayModel,
	}
	laterTimeOfDayWindow = &model.Window{
		After:  shortlyAfterNowTimeOfDayModel,
		Before: afterNowTimeOfDayModel,
	}
	shortlyAfterNowTimeOfDayModel = &model.Event{
		Rule:  "time-of-day",
		Param: "11:37:30",
	}
	afterNowTimeOfDayModel = &model.Event{
		Rule:  "time-of-day",
		Param: "12:00:00",
	}
	muchEarlierTimestampModel = &model.Event{
		Rule:  "timestamp",
		Param: "2014-10-29 03:00:00",
	}
	beforeNowTimestampModel = &model.Event{
		Rule:  "timestamp",
		Param: "2014-10-29 09:00:00",
	}
	afterNowTimestampModel = &model.Event{
		Rule:  "timestamp",
		Param: "2014-10-29 12:00:00",
	}
	shortlyAfterNowTimestampModel = &model.Event{
		Rule:  "timestamp",
		Param: "2014-10-29 11:37:30",
	}
	earlierTimestampWindow = &model.Window{
		After:  muchEarlierTimestampModel,
		Before: beforeNowTimestampModel,
	}
	overlappingTimestampWindow = &model.Window{
		After:  beforeNowTimestampModel,
		Before: afterNowTimestampModel,
	}
	laterTimestampWindow = &model.Window{
		After:  shortlyAfterNowTimestampModel,
		Before: afterNowTimeOfDayModel,
	}
	shortDelayModel = &model.Event{
		Rule:  "delay",
		Param: "00:15:00",
	}
	delayModel = &model.Event{
		Rule:  "delay",
		Param: "00:45:00",
	}
	longDelayModel = &model.Event{
		Rule:  "delay",
		Param: "04:00:00",
	}
	oneMinuteDelay = &model.Event{
		Rule:  "delay",
		Param: "00:01:00",
	}
	fourMinuteDelay = &model.Event{
		Rule:  "delay",
		Param: "00:04:00",
	}
	earlierTimestampOpenDelayCloseWindow = &model.Window{
		After:  muchEarlierTimestampModel,
		Before: delayModel,
	}
	overlappingTimestampOpenDelayCloseWindow = &model.Window{
		After:  beforeNowTimestampModel,
		Before: longDelayModel,
	}
	laterTimestampOpenDelayCloseWindow = &model.Window{
		After:  shortlyAfterNowTimestampModel,
		Before: delayModel,
	}
	earlierTimeOfDayOpenDelayCloseWindow = &model.Window{
		After:  muchEarlierTimeOfDayModel,
		Before: delayModel,
	}
	sunriseSunsetWindow = &model.Window{
		After:  sunriseModel,
		Before: sunsetModel,
	}
	sunsetModel = &model.Event{
		Rule: "sunset",
	}
	sunriseModel = &model.Event{
		Rule: "sunrise",
	}
	dawnModel = &model.Event{
		Rule: "dawn",
	}
	duskModel = &model.Event{
		Rule: "dusk",
	}
	bogusModel = &model.Event{
		Rule: "bogus",
	}
	bogusTimestamp = &model.Event{
		Rule:  "timestamp",
		Param: "bogus",
	}
	taskWithEarlierTimeOfDayOpenDelayCloseWindow = &model.Task{
		ID:     "task",
		Window: earlierTimeOfDayOpenDelayCloseWindow,
		Open: []*model.Action{
			{
				ActionType: "thing-action",
				Action:     "turnOn",
				ThingID:    "some-thing",
			},
		},
		Close: []*model.Action{},
	}
	taskWithEarlierTimestampOpenDelayCloseWindow = &model.Task{
		ID:     "task",
		Window: earlierTimestampOpenDelayCloseWindow,
		Open: []*model.Action{
			{
				ActionType: "thing-action",
				Action:     "turnOn",
				ThingID:    "some-thing",
			},
		},
		Close: []*model.Action{},
	}
	windowWithTwoDelayEvents = &model.Window{
		After:  oneMinuteDelay,
		Before: fourMinuteDelay,
	}
	taskWithTwoDelayEvents = &model.Task{
		ID:     "task",
		Window: windowWithTwoDelayEvents,
		Open: []*model.Action{
			{
				ActionType: "thing-action",
				Action:     "turnOn",
				ThingID:    "some-thing",
			},
		},
		Close: []*model.Action{},
	}
)
