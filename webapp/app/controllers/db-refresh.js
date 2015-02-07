angular.module('schedulerApp.controller.db-refresh', [
  'schedulerApp.db'
])
.config(['$routeProvider', function($routeProvider) {
  $routeProvider.when('/edit/:id', {templateUrl: 'views/task-edit.html', controller: 'TaskEdit'});
  $routeProvider.when('/create', {templateUrl: 'views/task-edit.html', controller: 'TaskEdit'});
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

