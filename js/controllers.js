function readerController($scope,$http,UserFeedBind)
{
    $http({ 
    	method: 'GET', 
    	url: 'http://localhost:9090/UserFeed/_NozVwtdndMkCeguF6a70ao6mI7GKbrF8mtHc-8r7bI=',
    	isArray:true
    }).success(function(data, status, headers, config) {
    	$scope.UserFeeds = data;
		  $scope.rawr = status;
  	}).error(function(data, status, headers, config) {
    	$scope.rawr = status;
  	});
}
