'use strict';

// Declare app level module which depends on views, and components
angular.module('schedulerApp.db', [
	'schedulerApp.db.rest',
	'schedulerApp.db.mock'
])
.factory('db', [ '$location', 'dbMock', 'dbRest', function($location, dbMock, dbRest) {

	if ($location.host() == 'localhost') {
		return dbMock
	} else {
		return dbRest
	}

}])