'use strict';

// Declare app level module which depends on views, and components
angular.module('schedulerApp', [
  'ngRoute',
  'ngResource'
])
.factory('db', [ '$resource', '$q', '$rootScope', '$location', function($resource, $q, $rootScope, $location) {

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
.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/list', {templateUrl: 'task-list.html', controller: 'TaskList'});
  $routeProvider.when('/edit/:id', {templateUrl: 'task-edit.html', controller: 'TaskEdit'});
  $routeProvider.when('/create', {templateUrl: 'task-edit.html', controller: 'TaskEdit'});
  $routeProvider.otherwise({redirectTo: '/list'});
}])
.controller('ModelRefresh', ['$scope', 'db', function($scope, db) {
	db.refresh().then(
		function(things) {
			return things
		},
		function(m) {
			if (m && m.data) {
				$scope.message = m.data
			} else {
				$scope.message = "An error occurred."
			}
			console.debug("an error occurred while refreshing the model: ", m)
		}
	)
}])
.controller('TaskEdit', ['$scope', '$location', 'db', '$routeParams', function($scope, $location, db, $routeParams) {

	$scope.thingToModel = function(thing) {
		var result = {}
		result.id = thing.id
		result.description = thing.name
		var room = db.rooms[thing.location]
		if (!room) {
			return null
		}
		result.room = room.name
		result.on = "true"
		return result
	}

	$scope.formatTime = function (date) {
	    var h = date.getHours();
	    var m = date.getMinutes();
	    var s = date.getSeconds();
	    var pad = function(d) { return d <= 9 ? '0'+d : d }
	    return ''+pad(h)+":"+pad(m)+":"+pad(s)
	}

	$scope.actionToModel = function(action) {
	    // {
        //    "action" : "turnOn",
        //    "thingID" : "94dc7d1a-aa67-11e4-803a-7c669d02a706",
        //     "type" : "thing-action"
        // }
        //
        // ==>
        //
	    // {
	    //	  "id": thing-id
        //    "description" : "thing description",
        //	  "room": "room description"
        //	  "on": true
        // }
        //

		var result = {}
		switch (action["action"]) {
		case "turnOn":
			result.on = "true"
			break
		case "turnOff":
			result.on = "false"
			break
		default:
			console.debug("bad action ", action)
			return null
		}

		if (action.type != "thing-action") {
			console.debug("bad type ", action.type)
			return null
		}

		if (!action.thingID) {
			console.debug("missing thing id")
			return null
		}

		result.id = action.thingID
		var thing = db.things[result.id]
		if (!thing) {
			console.debug("invalid thing id", result.id)
			return null
		}
		result.description = thing.name
		var room = db.rooms[thing.location]
		if (!room) {
			console.debug("invalid room id", thing.location)
			return null
		}
		result.room = room.name

		return result
	}

	$scope.selectedThings = {}
	$scope.availableThings = {}

	$scope.task = function() {
		return db.tasks[$routeParams["id"]]
	}

	$scope.$watch('task()', function(task) {
		$scope.availableThings = {}
		angular.forEach(db.things, function(t) {
			var m = $scope.thingToModel(t)
			if (m) {
				$scope.availableThings[t["id"]] = m
				$scope.selected = t["id"]
			}
		})

		$scope.timeOfDay = $scope.formatTime(new Date(new Date().valueOf()+(60*1000)))
		$scope.description = '@ '+$scope.timeOfDay

		if (!task) {
			return
		}

		$scope.description = task.description
		$scope.timeOfDay = task.window.after.param


		angular.forEach(task.open, function(action) {
			var thing = $scope.actionToModel(action)
			if (thing) {
				$scope.selectedThings[thing.id] = thing
				delete $scope.availableThings[thing.id]
			} else {
				console.debug(action)
			}
		})
		$scope.updateSelected()
	})

	$scope.selected = null

	$scope.save = function() {
		$scope.message = ""

		var task = {
			"id": $routeParams["id"],
			"description": (($scope.description == '') ? '@ '+$scope.timeOfDay : $scope.description ),
			"open":
				(function() {
					var results = []
					angular.forEach($scope.selectedThings, function(m) {
						results.push(
							{
								"type": "thing-action",
								"action": (m.on == "true" ? "turnOn" : "turnOff"),
								"thingID": m.id
							}
						)
					})
					return results
				}()),
			"close": [],
			"window": {
				"after": {
					"rule": "time-of-day",
					"param": $scope.timeOfDay
				},
				"before": {
					"rule": "delay",
					"param": "00:01:00"
				}
			}
		}

		db.save(task).then(
			function() {
				$location.path('/list')
			}, function(m) {
				if (m && m.data) {
					$scope.message = m.data
				} else {
					$scope.message = "An error occurred."
				}
				console.debug(m)
			})
	}

	$scope.delete = function() {
		db.delete($routeParams["id"]).then(
			function() {
				$location.path('/list')
			}, function(m) {
				if (m && m.data) {
					$scope.message = m.data
				} else {
					$scope.message = "An error occurred."
				}
				console.debug(m)
			})
	}

	$scope.updateSelected = function() {
		$scope.selected = null
		angular.forEach($scope.availableThings, function(m) {
			$scope.selected = m.id
		})
	}
	$scope.add = function() {
		if ($scope.selected) {
			$scope.selectedThings[$scope.selected] = $scope.availableThings[$scope.selected]
			delete $scope.availableThings[$scope.selected]
			$scope.updateSelected()
		}

	}

	$scope.remove = function() {
		var tmp = []
		angular.forEach($scope.selectedThings, function (m) {
			if (m.selected) {
				tmp.push(m.id)
			}
		})
		angular.forEach(tmp, function(id) {
			var m = $scope.selectedThings[id]
			delete $scope.selectedThings[id]
			$scope.availableThings[id] = m
		})
		$scope.updateSelected()
	}

	$scope.cancel = function() {
		$location.path('/list')
	}

	$scope.isDescriptionFrozen = false
	$scope.freezeDescription = function() {
		$scope.isDescriptionFrozen = true
	}

	$scope.$watch('timeOfDay', function() {
		if (!$scope.isDescriptionFrozen) {
			$scope.description = '@ '+$scope.timeOfDay
		}
	})
}])
.controller('TaskList', ['$scope', 'db', '$rootScope', function($scope, db, $rootScope) {
	$scope.tasks = {}
	$scope.db = db
	$scope.$watch("db.tasks", function(tasks) {
		$scope.tasks = db.tasks
	})
}])

