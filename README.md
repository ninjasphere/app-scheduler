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
			"window": {
				"from": "...",  // dusk or dawn, local-time-of-day, local-timestamp
				"to": "..."	    // dawn or dawn, local-time-of-day, local-timestamp
			},
			"open": [
				{
					"type": "thing-action",
					"thing-id": "ed9f4064",
					"action": "on"
				}
			],
			"close": [
				{
					"type": "thing-action",
					"thing-id": "ed9f4064",
					"action": "off"
				}
			]
		]
	}
