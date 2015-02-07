'use strict';

// Declare app level module which depends on views, and components
angular.module('schedulerApp.db.mock', [
  'ngResource'
])
.factory('dbMock', [ '$resource', '$q', function($resource, $q) {

	var
		counter = 0,
		service = {
			rooms:  $resource("services/mocks/rooms.js").get(),
			things: $resource("services/mocks/things.js").get(),
			tasks: { },
			save: function(task) {
				if (! task.id || task.id == '') {
					task.id = ++counter;
				}
				service.tasks[task.id] = task
				var deferred = $q.defer()
				deferred.resolve(service)
				return deferred.promise
			},
			delete: function(id) {
				delete service.tasks[task.id]
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
	return service
}])