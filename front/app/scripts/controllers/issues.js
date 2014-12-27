
angular.module('govcodeApp')
  .controller('IssuesCtrl', ['$rootScope', 
                           '$scope',
                           '$http',
                           '$routeParams',
                           function ($rootScope, $scope, $http, $routeParams) {
    $scope.config = {
      itemsPerPage: 10,
      fillLastPage: false
    }

    $scope.helpWantedIssues = [];



    $http.get($rootScope.apiUrl + '/issues?label=help').success(function (data) {
      $scope.helpWantedIssues = data;
    });

    $http.get($rootScope.apiUrl + '/issues').success(function (data) {
      $scope.otherIssues = data;
    });

}]);