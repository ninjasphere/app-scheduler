'use strict';

angular.module('schedulerApp.controller.task-list', [
  'ngRoute',
  'ngResource',
  'schedulerApp.db',
])
.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/list', {templateUrl: 'views/task-list.html', controller: 'TaskList'});
}])
.controller('TaskList', ['$scope', 'db', '$rootScope', function($scope, db, $rootScope) {
	$scope.tasks = {}
	$scope.db = db
	$scope.$watch("db.tasks", function(tasks) {
		$scope.tasks = db.tasks
	})
}])

