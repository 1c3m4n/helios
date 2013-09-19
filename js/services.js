angular.module('FeedService', ['ngResource']).factory('UserFeedBind', function($resource) { 
	return $resource('http://localhost\\:9090/UserFeed/:emailhash', {}, {
    query: {method:'GET', params:{emailhash:'emailhash'}, isArray:true}
	}); 
});


angular.module('UnreadService', ['ngResource']).factory('UserUnreadFeedBind', function($resource) { 
	return $resource('http://localhost\\:9090/UserFeed/UserUnreadFeed/:emailhash/:feedhash', {}, {
    query: {method:'GET', params:{emailhash:'emailhash'}, isArray:true},
    queryfeed: {method:'GET', params:{emailhash:'emailhash',feedhash:'feedhash'}, isArray:true}
	}); 
});






