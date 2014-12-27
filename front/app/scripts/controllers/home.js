
angular.module('govcodeApp')
  .controller('HomeCtrl', ['$rootScope', 
                           '$scope',
                           '$http',
                           '$routeParams',
                           '$location',
                           function ($rootScope, $scope, $http, $routeParams, $location) {

    $scope.mapLang = mapLang;

    // Get home stats
    $http.get($rootScope.apiUrl + '/', { cache: true }).success(function (data) {
      $scope.home_stats = data;
    });

    // Get latest repos
    $http.get($rootScope.apiUrl + '/repos?perPage=5', { cache: true }).success(function (data) {
      $scope.latestRepos = data.slice(0,5);
    });

    // Get latest issues
    $http.get($rootScope.apiUrl + '/issues?perPage=5', { cache: true }).success(function (data) {
      $scope.latestIssues = data;
    });

    // Get latest issues
    $http.get($rootScope.apiUrl + '/issues?perPage=5&label=help', { cache: true }).success(function (data) {
      $scope.helpWantedIssues = data;
    });

    $scope.selLang = function(event) {
      var obj = $(event.currentTarget);
      if (obj.hasClass("active")) {
        obj.removeClass("active");
      } else {
        obj.addClass("active");
      }
    };

    $scope.seeIssues = function() {
      var langs = $(".selector .opt.active").map(function() {
        return $(this).attr("rel");
      }).get().join();

      $location.url("/issues/?languages=" + langs);
    };

}]);