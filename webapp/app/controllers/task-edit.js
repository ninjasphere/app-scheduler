'use strict';

angular.module('schedulerApp.controller.task-edit', [
  'ngRoute',
  'ngResource',
  'schedulerApp.db',
])
.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/edit/:id', {templateUrl: 'views/task-edit.html', controller: 'TaskEdit'});
  $routeProvider.when('/create', {templateUrl: 'views/task-edit.html', controller: 'TaskEdit'});
}])
.directive('actionLabel', [function() {
	return {
		'restrict': 'A',
		'scope': {
			"model": "=actionLabel"
		},
		'link': function(scope, elt, attr) {
			scope.$watch("model.on", function() {
				switch (scope.model.actionType) {
				case "thing-action":
					elt.html("turn it <strong>"+(scope.model.on?"on":"off")+"</strong>")
					break
				case "presets-action":
					elt.html("<strong>apply</strong>")
					break
				}
			})
		}
	}
}])
.controller('TaskEdit', ['$scope', '$location', 'db', '$routeParams', function($scope, $location, db, $routeParams) {

	$scope.actionModels = {};
	$scope.isDescriptionFrozen = false
	$scope.repeatDaily = true

	$scope.duration = ''

	var model = function() {
	}
	model.prototype.actionLabel = function() {
	}

	$scope.thingToModel = function(thing) {
		var result = new model()
		result.id = "thing:"+thing.id
		result.description = thing.name
		result.actionType = 'thing-action'
		var room = db.rooms[thing.location]
		if (!room) {
			return null
		}
		result.room = room.name

		result.on = !(thing.onOffChannel
			&& thing.onOffChannel.lastState
			&& thing.onOffChannel.lastState.payload == true)
		return result
	}

	$scope.sceneToModel = function(scene) {
		var result = new model(),
			location=''

		result.id = "scene:"+scene.id
		result.description = scene.label
		result.actionType = 'presets-action'
		var scopeParts = /(.*):(.*)/.exec(scene.scope)
		if (scopeParts.length < 3) {
			console.debug("invalid scope", action)
			return null
		}
		if (scopeParts[1] != "site") {
			location = scopeParts[2]
		}

		if (location != '') {
			var room = db.rooms[location]
			if (!room) {
				return null
			}
			result.room = room.name
		} else {
			result.room = 'Home'
		}

		result.on =  true
		return result
	}

	$scope.formatTime = function (date, ui) {
	    var h = date.getHours();
	    var m = date.getMinutes();
	    var s = date.getSeconds();
	    var pad = function(d) { return d <= 9 ? '0'+d : d }
	    return ''+pad(h)+":"+pad(m)+(s != 0 || !ui ?  ":"+pad(s) : "")
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

		var result = new model(),
			location = ''

		if (!action.subject) {
			console.debug("missing thing id", action)
			return null
		}

		result.id = action.subject
		var parts = /(.*):(.*)/.exec(action.subject)
		if (parts.length < 3) {
			console.debug("invalid subject id", action)
			return null
		}


		switch (action.type) {
		case "thing-action":
			switch (action["action"]) {
			case "turnOn":
				result.on = true
				break
			case "turnOff":
				result.on = false
				break
			default:
				console.debug("bad action ", action)
				return null
			}
			if (parts[1] != 'thing') {
				console.debug("thing-action with something other than a thing id", action)
			}
			var thing = db.things[parts[2]]
			if (!thing) {
				console.debug("invalid thing id", action)
				return null
			}
			result.description = thing.name
			location = thing.location
			if (location == '') {
				console.debug("invalid location", action)
				return null
			}
			break;
		case "presets-action":
			switch (action["action"]) {
			case "applyScene":
				result.on = true
				break
			default:
				console.debug("bad action ", action)
				return null
			}
			if (parts[1] != 'scene') {
				console.debug("presets-action with something other than a scene id", action)
				return null
			}
			var scene = db.scenes[parts[2]]
			if (!scene) {
				console.debug("invalid scene id", action)
				return null
			}
			result.description = scene.label
			var scopeParts = /(.*):(.*)/.exec(scene.scope)
			if (scopeParts.length < 3) {
				console.debug("invalid scope", action)
				return null
			}
			if (scopeParts[1] == "site") {
				result.room = "Home"
				location = ''
			} else {
				location = scopeParts[2]
			}
			break;
		default:
			console.debug("bad type ", action)
			return null
		}

		result.actionType = action.type

		if (location != '') {
			var room = db.rooms[location]
			if (!room) {
				console.debug("invalid room id", location)
				return null
			}
			result.room = room.name
		}

		result.selected = true

		return result
	}

	$scope.task = function() {
		return db.tasks[$routeParams["id"]]
	}

	$scope.timestamp = function(ts) {
		var
			now = new Date(),
			ts = new Date($scope.formatDate(now).replace(/-/g, '/')+" "+(ts ? ts : $scope.timeOfDay))
			if (ts.getFullYear() && ts < now) {
				ts.setDate(ts.getDate()+1)
			}
			return ts
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
				ts = $scope.timestamp()

			if (!ts.getFullYear()) {
				$scope.message = "enter a time of the form hh:mm:dd"
				return
			}

			if ($scope.repeatDaily) {
				rule = "time-of-day"
				param = $scope.formatTime(ts, false)
			} else {
				rule = "timestamp"
				param = $scope.formatDate(ts) + " " + $scope.formatTime(ts, false)
			}
			break
		}

		var desc = (($scope.description == '') ? '@ '+$scope.timeOfDay : $scope.description )

		var actions = (function() {
			var
				open=[],
				close=[]

			angular.forEach($scope.actionModels, function(m) {
				if (m.selected) {

					var openAction, closeAction

					switch(m.actionType) {
					case "thing-action":
						openAction =  (m.on ? "turnOn" : "turnOff")
						closeAction =  (!m.on ? "turnOn" : "turnOff")
						break
					case "presets-action":
						openAction = "applyScene"
						closeAction = "undoScene"
						break
					default:
						console.debug("logic error: actionType not supported", m)
						return
					}

					var obj = {
						"type": m.actionType,
						"action": openAction,
						"subject": m.id
					}
					open.push(obj)
					obj = angular.copy(obj)
					obj.action = closeAction
					close.push(obj)
				}
			})
			return [open,close]
		}())

		var
			beforeParam

		if ($scope.duration == '') {
			beforeParam = "00:01:00"
			actions[1] = []
		} else {
			beforeParam = $scope.duration
		}

		var task = {
			"id": $routeParams["id"],
			"description": desc,
			"tags": ["simple-ui"],
			"open": actions[0],
			"close": actions[1],
			"window": {
				"after": {
					"rule": rule,
					"param": param
				},
				"before": {
					"rule": "delay",
					"param": beforeParam
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
				$scope.actionModels["thing:"+t["id"]] = m
			}
		})

		angular.forEach(db.scenes, function(s) {
			var m = $scope.sceneToModel(s)
			if (m) {
				$scope.actionModels["scene:"+s["id"]] = m
			}
		})

		$scope.isDescriptionFrozen = false
		$scope.timeOfDay = $scope.formatTime(new Date(new Date().valueOf()+(60*1000))).substring(0,5)
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
				console.debug("found an action for a subject that no longer exists: ", action)
			}
		})

		if (!task.window || !task.window.before || task.window.before.rule != 'delay') {
			console.debug("invalid before rule", task)
			return
		}

		console.debug(task)
		if (task.close && task.close.length != 0) {
			$scope.duration = task.window.before.param
		}


		switch (task.window.after.rule) {
		case "sunrise":
		case "sunset":
		case "dawn":
		case "dusk":
			$scope.timeOfDay = task.window.after.rule
			break;
		case "time-of-day":
			$scope.timeOfDay = $scope.formatTime($scope.timestamp(task.window.after.param), true)
			break;
		case "timestamp":
			var ts = new Date(task.window.after.param)
			if (ts.getFullYear()) {
				$scope.timeOfDay = $scope.formatTime(ts, true)
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
		var
			base = '@ '+$scope.timeOfDay,
			ts = $scope.timestamp(),
			now = new Date()

		if (!$scope.repeatDaily && ts.getFullYear()) {
			if ($scope.formatDate(now) == $scope.formatDate($scope.timestamp())) {
				return base + ' today'
			} else {
				return base + ' tomorrow'
			}
		} else {
			return base
		}
	}

	$scope.$watch('timeOfDay+repeatDaily', function() {
		if (!$scope.isDescriptionFrozen) {
			$scope.description = $scope.generatedDescription()
		}
	});

}])

