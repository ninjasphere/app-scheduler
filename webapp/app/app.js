'use strict';

// Declare app level module which depends on views, and components
angular.module('schedulerApp', [
  'ngRoute',
  'ngResource'
])
.factory('scheduler', [ '$resource', function($resource) {

	var
		Things = $resource("http://10.0.1.168:8000/rest/v1/things", {})
		service = {
			refresh: function() {
				return $resource.get({}).$promise.then(
					function(things) {
						return things
					}
				).$promise
			}

		}
	return service
}])
.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/list', {templateUrl: 'task-list.html', controller: 'TaskList'});
  $routeProvider.when('/edit/:task', {templateUrl: 'task-edit.html', controller: 'TaskEdit'});
  $routeProvider.when('/create', {templateUrl: 'task-edit.html', controller: 'TaskEdit'});
  $routeProvider.otherwise({redirectTo: '/list'});
}])
.controller('TaskList', ['$scope', function($scope) {
	$scope.tasks = [ {
			"id": "1",
			"description": "task 1"
	}]
}])
.controller('TaskEdit', ['$scope', '$location', 'scheduler', function($scope, $location, scheduler) {
	$scope.description = "schedule item"
	$scope.timeOfDay = "09:30:00"
	$scope.selectedThings = [
		{
			"description": "Downlight",
			"room": "Lounge",
			"on": true,
		}
	]

	$scope.availableThings = [
		{
			"description": "Lamp",
			"room": "Bedroom"
		}
	]

	scheduler.refresh().done(
		function(things) {
			console.debug(things)
		},
		function(err) {
			$scope.message = err
		}
	)

	$scope.save = function() {

	}

	$scope.delete = function() {

	}

	$scope.cancel = function() {
		$location.path('/list')
	}
}])
