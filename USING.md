
Introducing the Ninja Sphere Scheduler App
==========================================

The Ninja Sphere Scheduler App is a new Ninja Labs application that runs on the Ninja Sphere and allows you to turn your things on and off at different times of day.
You can interact with the application using a web browser in your phone, tablet or PC.

Open The Scheduler App
----------------------
Type the following URL into a web browser on your phone, PC or web browser.

	http://ninjasphere:8100

If this doesn't work for you, you might need to replace the word 'ninjasphere' with the actual IP address of your Ninja Sphere, which you can discover
using the procedure "Determine The IP Address Of Your Sphere" below.

Note that because the application user interface is presented as a web app, your browser needs to be running on a device which is connected to the same WiFi network
as your Ninja Sphere. Later versions of the Ninja Sphere Scheduler App will be integrated into the Ninja Sphere Phone App and this restriction will not apply.

Scheduling A Task
-----------------
A task is an action that can be scheduled at some later a time.

To create a new task, open the web application and "schedule a new task" button at the top of the page. This action will present you with a form that allows
you to configure the details of the task.

For each task, you can configure the following details:

* a description of the task
* the time at which the action will be taken
* the list of things to act upon
* for each thing to be acted upon, the "Turn On" or "Turn Off" action for each
* whether the task should be run once or at the same time each day
* whether the reverse action should be taken after a delay.

Adjust the time for the scheduled action by supplying the hour and minute that you want the action to take place. In addition to a time of day, you can also specify times like: dawn, dusk, sunrise or sunset. These times depend on the physical location of your Ninja Sphere and are calculated based on the geographic coordinates you supplied when you setup the Ninja Sphere for the first time.

To select things to act upon, click the check box next to each thing you would like to act upon. Tap the "Turn It On/Off" button to specify whether the thing should be turned on or off at the specified time.

If the task is to be performed only once, click the "once off" button. Otherwise, click the "daily" button to have the task execute every day. A "once off" task will be automatically removed from the
schedule once it has been run.

If you would like to have the opposite actions to be performed after a delay, enter the duration of the delay into the "duration" field.

Once you have entered the details of the task, press the "Save" button to add the task to the schedule. The task will then appear in the list of scheduled tasks.

The list of tasks allows you edit select an existing task to edit or delete it.

Determine The IP Address Of Your Ninja Sphere
---------------------------------------------

To access the app, you first need to determine the IP address of your Ninja Sphere using the following steps.

1. Open the Ninja Sphere phone app
2. Press the menu button in the top left hand corner
3. Scroll down to the bottom of the pane
4. Take note of the route string which should read something like:

	Local via 10.0.1.14

	If the route string reads "Cloud" make sure that your phone is connected to the same WiFi network that the Sphere is on.

