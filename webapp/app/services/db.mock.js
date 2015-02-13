'use strict';

// Declare app level module which depends on views, and components
angular.module('schedulerApp.db.mock', [
  'ngResource'
])
.factory('dbMock', [ '$resource', '$q', function($resource, $q) {

	var
		counter = 0,
		scenes = {},
		service = {
			rooms:  $resource("services/mocks/rooms.js").get(),
			things: $resource("services/mocks/things.js").get(),
			scenes: scenes,
			tasks: { },
			save: function(task) {
				if (! task.id || task.id == '') {
					task.id = ++counter;
				}
				service.tasks[task.id] = task
				var deferred = $q.defer()
				deferred.resolve(service)
				console.debug("saved task", task)
				return deferred.promise
			},
			delete: function(id) {
				delete service.tasks[id]
				var deferred = $q.defer()
				deferred.resolve(service)
				return deferred.promise
			},
			refresh: function() {
				var deferred = $q.defer()
				deferred.resolve(service)
				return deferred.promise
			}
		}
	$resource("services/mocks/scenes.js").query({}, function(a){
		angular.forEach(a, function(e) {
			scenes[e.id] = e
		})
	})
	return service
}])