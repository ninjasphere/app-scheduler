package controller

import (
	"github.com/ninjasphere/app-scheduler/model"
	"time"
)

var (
	testTime   = time.Date(2014, 10, 29, 11, 22, 30, 0, time.Now().Location())
	futureTime = time.Date(2014, 10, 29, 12, 00, 00, 0, time.Now().Location())

	beforeNowTimeOfDayModel = &model.Event{
		Rule:  "time-of-day",
		Param: "09:00:00",
	}
	shortlyAfterNowTimeOfDayModel = &model.Event{
		Rule:  "time-of-day",
		Param: "11:37:30",
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
