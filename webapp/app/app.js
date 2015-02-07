'use strict';

// Declare app level module which depends on views, and components
angular.module('schedulerApp', [
  'ngRoute',
  'ngResource',
  'schedulerApp.db',
  'schedulerApp.controller.task-list',
  'schedulerApp.controller.task-edit',
  'schedulerApp.controller.db-refresh'
])
.config(['$routeProvider', function($routeProvider) {
  $routeProvider.otherwise({redirectTo: '/list'});
}])
