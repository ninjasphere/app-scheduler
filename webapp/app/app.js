'use strict';

// Declare app level module which depends on views, and components
angular.module('schedulerApp', [
  'ngRoute'
]).
config(['$routeProvider', function($routeProvider) {
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
.controller('TaskEdit', ['$scope', '$location', function($scope, $location) {
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

	$scope.save = function() {

	}

	$scope.delete = function() {

	}

	$scope.cancel = function() {
		$location.path('/list')
	}
}])
