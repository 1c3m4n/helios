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

$scope.loadFeed= function(userfeed) {
        $scope.Title = userfeed.FeedName
        $scope.Content = "Selected User Feed Items"
    }
    $scope.loadAll = function() {
        $scope.Title = "All"
        $scope.Content = "All User Feed Items"
    }
    $scope.loadToday = function() {
        $scope.Title = "Today's Feed"
        $scope.Content = "Today's user feed items"
    }
    $scope.loadFriends = function() {
        $scope.Title = "Friend's Shared Items"
        $scope.Content = "Friend's shared feed items"
    }
    $scope.loadLiked = function() {
        $scope.Title = "Liked Items"
        $scope.Content = "Liked user feed items"
    }
    $scope.loadShared = function() {
        $scope.Title = "Your Shared Items"
        $scope.Content = "User shared items"
    }
    $scope.loadPopular = function() {
        $scope.Title = "Popular Feeds"
        $scope.Content = "popular feed items"
    }
}


