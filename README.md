NAME
====
app-scheduler

DESCRIPTION
===========
app-scheduler is an application that manages the scheduled execution of tasks. A task delays until its 'after' event occurs,
executes its 'open' actions, delays until the 'before' event occurs, then executes the 'close' actions. If the 'after'
and 'before' events are both recurring events, then repeats the cycle again, otherwise the task ends.

'after' and 'before' events can be specified with different kinds of event generation rules. The following rules exist:

<dl>
	<dt>timestamp</dt>
	<dd>Specifies an absolute timestamp which can occur at most once.</dd>
	<dt>time-of-day</dt>
	<dd>Specifies a time of day that will recur once each day.</dd>
	<dt>delay</dt>
	<dd>Specifies a delay after the previous event. In the case of a 'after' event this delay will be measured
	from the start of the task (for the first iteration of the task loop) or from the last 'before' event otherwise. In the
	case of an 'before' event, this delay will be measured from the time of the last 'after' event.</dd>
	<dt>sunset</dt>
	<dd>Specifies sunset in the local time zone.</dd>
	<dt>sunrise</dt>
	<dd>Specifies sunrise in the local time zone.</dd>
	<dt>dawn</dt>
	<dd>Specifies civil dawn (sun 6 degrees below horizon) in the local time zone.</dd>
	<dt>dusk</dt>
	<dd>Specifies civil dusk (sun 6 degrees below horizon) in the local time zone.</dd>
</dl>

The schedule is passed to the application on startup as its configuration object. The schedule looks like this:

	{
	    "schedule": [
	        {
	            "id": "ed9f4064-4f2a-4d1a-8583-c05f19af0b58",
	            "description": "Turn on lounge lights for 30 minutes at 19:00",
	            "window": {
	                "after": {
	                    "rule": "time-of-day",
	                    "param": "19:00:00"
	                },
	                "before": {
	                    "rule": "delay",
	                    "param": "00:30:00"
	                }
	            },
	            "open": [
	                {
	                    "type": "thing-action",
	                    "thingID": "ed9f4064",
	                    "action": "turnOn"
	                }
	            ],
	            "close": [
	                {
	                    "type": "thing-action",
	                    "thingID": "ed9f4064",
	                    "action": "turnOff"
	                }
	            ]
	        }
	    ]
	}

A command line utility, nscheduler, can be used to schedule and cancel individual tasks. For example, the task with the above definition would be generated with a nscheduler invocation that looked like this:

	nscheduler \
		   --after time-of-day 19:00:00 \
		   --before delay 00:30:00 \
		   --on-open turnOn \
		   --on-close turnOff \
		   --thing ed9f4064 \
		   -- schedule

For more details about the nscheduler command, run nscheduler without any arguments to read some help text.

TOPICS
======

The scheduler service listens on a topic called:

	$node/{serial}/app/com.ninjablocks.scheduler/service/scheduler

It supports the following methods:

	schedule {task-model}
	cancel {task-id}

The scheduler also listens on topics of the form:

	$device/:deviceId/channel/user-agent

for 'schedule-task' and 'cancel-task' events.