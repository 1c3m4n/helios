angular.module('readerService', ['ngResource'])
.factory('UserFeedBind', function($resource) { 
	return $resource('http://localhost\\:9090/UserFeed/:emailhash', { emailhash: '@emailhash' }); 
});
