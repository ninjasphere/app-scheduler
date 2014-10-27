NAME
====
app-scheduler

DESCRIPTION
===========
Implements a simple application that uses a time based schedule to perform simple actions in response to the progress of time.

The schedule is passed to the application on startup.

{
	"schedule" : [
	{
		"uuid": "ed9f4064-4f2a-4d1a-8583-c05f19af0b58",
		"description": "Turn on lounge lights at night",
		"open": {
		      "at": "22:00", // or a fuzzier time specification, resolved by what?
		      "action": [
			      {
			      	"action-type": "call",
			      	"channel": "$device/{device-id}/channel/{channel-id}",
			      	"method": "turnOn",
			      	"params": {}
			      }
		      ]

		},
		"close": {
		      "at": "06:00",
		      "action": [
			      {
			      	"action-type": "call",
			      	"channel": "$device/{device-id}/channel/{channel-id}",
			      	"method": "turnOff",
			      	"params": {}
			      }
		      ]

		},
	}
	]
}
