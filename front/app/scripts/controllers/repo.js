'use strict';

/**
 * @ngdoc function
 * @name govcodeApp.controller:RepoCtrl
 * @description
 * # RepoCtrl
 * Show a repo detail
 */

function sortObj(arr) {
  // Setup Arrays
  var sortedKeys = new Array();
  var sortedObj = {};

  // Separate keys and sort them
  for (var i in arr){
      sortedKeys.push(i);
  }
  sortedKeys.sort();
  // Reconstruct sorted obj based on keys
  for (var i in sortedKeys){
      sortedObj[sortedKeys[i]] = arr[sortedKeys[i]];
  }

  return sortedObj;
}


angular.module('govcodeApp')
  .controller('RepoCtrl', ['$rootScope', 
                           '$scope',
                           '$http',
                           '$routeParams',
                           function ($rootScope, $scope, $http, $routeParams) {
    // Get a repo
    $http.get($rootScope.apiUrl + '/repos/' + $routeParams.repoName).success(function (data) {
      // Load the repos in the scope
      $scope.repo = data;

      // Create Graph info
      if (data && data.RepoStat) {
        var graphInput = data.RepoStat;

        graphInput.sort(function(a,b) {
          return (a.Week <= b.Week) ? -1 : ((a.Week > b.Week) ? 1 : 0);
        });

        var records = {};

        var start_date = new Date(graphInput[0].Week); 
        var now = Date.now();

        for (var d = new Date(start_date.valueOf()); d <= now; d.setDate(32)) {
          var date_str = "" + d.getFullYear() + "-" + (d.getMonth() + 1);
          records[date_str] = 0;
        }

        for (var i = 0; i < graphInput.length; i++) {
          var date = new Date(graphInput[i].Week);
          var date_str = "" + date.getFullYear() + "-" + (date.getMonth() + 1);
          records[date_str] += graphInput[i].Commits;
        };

        $scope.graphDataRecords = records;

        $scope.graphOptions = {
          scaleIntegersOnly: true
        }

        $scope.graphData = {
          labels: Object.keys(records),
          datasets: [{
            fillColor: "rgba(220,220,220,0.5)",
            strokeColor: "rgba(220,220,220,1)",
            pointColor: "rgba(220,220,220,1)",
            pointStrokeColor: "#fff",
            data: $.map(records, function (val) { return val })
          }]
        };
      }

      // Load latest activity
      var activity_url = 'https://api.github.com/repos/' + $scope.repo.Organization.Login + '/' + $scope.repo.Name + '/events';
      $http.get(activity_url).success(function (data) {
        $scope.latest_activity = $.map(data, function(el) {
          return githubSentences.convert(el);
        });
      });
    });

  }]);