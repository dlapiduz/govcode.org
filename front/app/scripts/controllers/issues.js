
angular.module('govcodeApp')
  .controller('IssuesCtrl', ['$rootScope', 
                           '$scope',
                           '$http',
                           '$routeParams',
                           '$location',
                           function ($rootScope, $scope, $http, $routeParams, $location) {
    $scope.config = {
      itemsPerPage: 10,
      fillLastPage: false
    }

    $scope.helpWantedIssues = [];

    var selectedLangs = ($location.search()['languages'] || '').split(',');
    var selectedOrgs = ($location.search()['organizations'] || '').split(',');

    function updateIssues(langs, orgs) {
      var filters = [];
      if (langs.length > 0) {
        filters.push("languages=" + langs.toString());
      }

      if (orgs.length > 0) {
        filters.push("organizations=" + orgs.toString());
      }
      
      $http.get($rootScope.apiUrl + '/issues?label=help&' + filters.join('&'), { cache: true }).success(function (data) {
        $scope.helpWantedIssues = data;
      });

      $http.get($rootScope.apiUrl + '/issues?' + filters.join('&'), { cache: true }).success(function (data) {
        $scope.otherIssues = data;
      });
    }


    $http.get($rootScope.apiUrl + '/orgs', { cache: true }).success(function(data) {
      $scope.orgs = $.map(
        data, 
        function(e) {
          var ticked = (selectedOrgs.indexOf(e) >= 0);
          return { name: e.Login, ticked: false };
      });

    });

    $http.get($rootScope.apiUrl + '/', { cache: true }).success(function(data) {
      $scope.langs = $.map(
        data.IssueLangs, 
        function(e) {
          var ticked = (selectedLangs.indexOf(e) >= 0);
          return { name: e, ticked: ticked };
      });
    });

    $scope.$watch('langs', function(newLangs, oldLangs) {
      var langs = $.map(
        newLangs,
        function(e) {
          if (e.ticked) {
            return e.name;
          };
      });

      if (langs.toString() != selectedLangs.toString()) {
        updateIssues(langs, selectedOrgs);
        selectedLangs = langs;
      }
    }, true);

    $scope.$watch('orgs', function(newOrgs, oldOrgs) {
      var orgs = $.map(
        newOrgs,
        function(e) {
          if (e.ticked) {
            return e.name;
          };
      });

      if (orgs.toString() != selectedOrgs.toString()) {
        updateIssues(selectedLangs, orgs);
        selectedOrgs = orgs;
      }
    }, true);

    updateIssues(selectedLangs, selectedOrgs);



}]);