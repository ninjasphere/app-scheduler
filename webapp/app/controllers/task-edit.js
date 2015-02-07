'use strict';

angular.module('schedulerApp.controller.task-edit', [
  'ngRoute',
  'ngResource',
  'schedulerApp.db',
])
.controller('TaskEdit', ['$scope', '$location', 'db', '$routeParams', function($scope, $location, db, $routeParams) {

	$scope.actionModels = {};
	$scope.isDescriptionFrozen = false
	$scope.repeatDaily = true

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

	$scope.formatDate = function (date) {
	    var y = date.getFullYear();
	    var m = date.getMonth()+1;
	    var d = date.getDate();
	    var pad = function(d) { return d <= 9 ? '0'+d : d }
	    return ''+y+"-"+pad(m)+"-"+pad(d)
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
		result.selected = true

		return result
	}

	$scope.task = function() {
		return db.tasks[$routeParams["id"]]
	}

	$scope.timestamp = function() {

	}

	$scope.save = function() {
		$scope.message = ""
		var rule, param

		switch ($scope.timeOfDay) {
		case "dawn":
		case "dusk":
		case "sunrise":
		case "sunset":
			rule = $scope.timeOfDay
			param = ""
			break
		default:
			var
				now = new Date(),
				ts = new Date($scope.formatDate(now)+" "+$scope.timeOfDay)

			if (!ts.getFullYear()) {
				$scope.message = "enter a time of the form hh:mm:dd"
				return
			} else {
				if (ts < now) {
					ts.setDate(ts.getDate()+1)
				}
				if ($scope.repeatDaily) {
					rule = "time-of-day"
					param = $scope.formatTime(ts)
				} else {
					rule = "timestamp"
					param = $scope.formatDate(ts) + " " + $scope.formatTime(ts)
				}
			}
			break
		}

		var desc = (($scope.description == '') ? '@ '+$scope.timeOfDay : $scope.description )

		var open = (function() {
			var results = []
			angular.forEach($scope.actionModels, function(m) {
				if (m.selected) {
					results.push(
						{
							"type": "thing-action",
							"action": (m.on == "true" ? "turnOn" : "turnOff"),
							"thingID": m.id
						}
					)
				}
			})
			return results
		}())

		var task = {
			"id": $routeParams["id"],
			"description": desc,
			"tags": ["simple-ui"],
			"open": open,
			"close": [],
			"window": {
				"after": {
					"rule": rule,
					"param": param
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

	$scope.cancel = function() {
		$location.path('/list')
	}

	$scope.freezeDescription = function() {
		$scope.isDescriptionFrozen = true
	}

	$scope.toggleSelect = function(model) {
		model.selected = !model.selected
	}

	$scope.toggleActionState = function(model) {
		model.on = ! model.on
	}

	$scope.setRepeatDaily = function(value) {
		$scope.repeatDaily = value;
	}

	$scope.$watch('task()', function(task) {
		angular.forEach(db.things, function(t) {
			var m = $scope.thingToModel(t)
			if (m) {
				$scope.actionModels[t["id"]] = m
			}
		})

		$scope.isDescriptionFrozen = false
		$scope.timeOfDay = $scope.formatTime(new Date(new Date().valueOf()+(60*1000)))
		$scope.description = '@ '+$scope.timeOfDay
		$scope.repeatDaily = true

		if (!task) {
			// console.debug("new task", $scope)
			return
		}

		$scope.isDescriptionFrozen = true
		$scope.description = task.description

		angular.forEach(task.open, function(action) {
			var model = $scope.actionToModel(action)
			if (model && $scope.actionModels[model.id]) {
				$scope.actionModels[model.id] = model
			} else {
				console.debug("found an action for a thing that no longer exists: ", action)
			}
		})

		switch (task.window.after.rule) {
		case "sunrise":
		case "sunset":
		case "dawn":
		case "dusk":
			$scope.timeOfDay = task.window.after.rule
			break;
		case "time-of-day":
			$scope.timeOfDay = task.window.after.param
			break;
		case "timestamp":
			var ts = new Date(task.window.after.param)
			if (ts.getFullYear()) {
				$scope.timeOfDay = $scope.formatTime(ts)
				$scope.repeatDaily = false
			}
			break;
		default:
			console.debug("can't edit rule of type: ", task.window.after.rule)
			return
		}

		// the description is frozen iff the generated description does not match the saved description
		$scope.isDescriptionFrozen = ($scope.description != '') && ($scope.description != $scope.generatedDescription())

		// console.debug("loaded task", $scope)

	})

	$scope.generatedDescription = function() {
		return '@ '+$scope.timeOfDay
	}

	$scope.$watch('timeOfDay', function() {
		if (!$scope.isDescriptionFrozen) {
			$scope.description = $scope.generatedDescription()
		}
	});

}])

