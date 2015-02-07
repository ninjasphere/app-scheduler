'use strict';

// Declare app level module which depends on views, and components
angular.module('schedulerApp.db.rest', [
  'ngResource'
])
.factory('dbRest', [ '$resource', '$q', '$rootScope', '$location', function($resource, $q, $rootScope, $location) {

	var
		Things = $resource("http://"+$location.host()+":8000/rest/v1/things", {}),
		Rooms = $resource("http://"+$location.host()+":8000/rest/v1/rooms", {}),
		Tasks = $resource("http://"+$location.host()+":"+$location.port()+"/rest/v1/tasks", {}),
		Task = $resource("http://"+$location.host()+":"+$location.port()+"/rest/v1/tasks/:id", {}),
		refreshThings = function() {
			return Things.get({}).$promise.then(
				function(things) {
					service.things = {}
					angular.forEach(things.data,
						function(thing) {
							var found = false
							if (!thing.promoted || !thing.name || !thing.location) {
								// console.debug("skipping thing because not promoted, described or assigned to room", thing)
								return
							}
							angular.forEach(thing.device.channels, function(channel) {
								if (channel.supportedMethods && channel.supportedMethods.indexOf("turnOff") >= 0) {
									found = true
								}
							})
							if (found) {
								service.things[thing["id"]] = thing
							}
						}
					)
					return service
				}
			)
		},
		refreshRooms = function() {
			return Rooms.get({}).$promise.then(
				function(rooms) {
					service.rooms = {}
					angular.forEach(rooms.data, function(room) {
						service.rooms[room["id"]] = room
					})
					return service
				}
			)
		},
		refreshTasks = function() {
			return Tasks.get({}).$promise.then(
				function(tasks) {
					service.tasks = {}
					angular.forEach(tasks.schedule, function(task) {

						if (!task.tags || task.tags.indexOf("simple-ui") < 0) {
							return;
						}

						if (task.window.after.rule == 'timestamp') {
							var
								now = new Date(),
								ts = new Date(task.window.after.param)
							if (! ts.getFullYear() || ts < now) {
								// console.debug("skipping expired or invalid task")
								return
							}
						}

						if (task.description && task.description != '') {
							service.tasks[task["id"]] = task
						}
					})
					return service
				}
			)
		},
		service = {
			things: { } ,
			rooms: { } ,
			tasks: { } ,
			save: function(task) {
				var deferred = $q.defer()
				Tasks.save(
					{},
					task,
					function(r) {
						refreshTasks().then(
							function() { deferred.resolve(r) },
							deferred.reject
						)},
					deferred.reject)
				return deferred.promise
			},
			delete: function(id) {
				var deferred = $q.defer()
				Task.delete({"id": id},
					function(r) {
						refreshTasks().then(function() { deferred.resolve(r)}, deferred.reject)
					},
					deferred.reject
				)
				return deferred.promise
			},
			refresh: function() {
				return $q.all([ refreshRooms(),	refreshThings(), refreshTasks() ]).then(
					function() {
						return service
					}
				)
			}
		}
	return service
}])