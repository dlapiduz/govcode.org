'use strict';

/**
 * @ngdoc function
 * @name govcodeApp.controller:StatsCtrl
 * @description
 * # StatsCtrl
 * Controller of the govcodeApp
 */
angular.module('govcodeApp')
  .controller('StatsCtrl', [
      "$rootScope",
      "$scope",
      "$http",
      function ($rootScope, $scope, $http) {

        var colors = ["#EA5172",
                      "#501367",
                      "#16E56D",
                      "#E49956",
                      "#AD51DE",
                      "#5160BB",
                      "#07DAF5",
                      "#E0E0E0",
                      "#F14B8E",
                      "#0887E5",
                      "#E49956",
                      "#F5F416",
                      "#F04919"];

        $scope.orgFilter = {};
        $scope.countGraphOptions = {
          scaleIntegersOnly: true,
          showTooltips: false,
          scaleBeginAtZero: true
        }

        $scope.countGraphData = {
          datasets: []
        }

        $scope.commitGraphData = {
          datasets: []
        }

        $scope.commitGraphOptions = {
          showTooltips: true,
          animation: false,
          datasetFill: false,
          multiTooltipTemplate: "<% if (datasetLabel){%><%=datasetLabel%>: <%}%><%= value %>"
        }

        $scope.allOrgs = function(show) {
          $.each($scope.orgFilter, function(key, val) {
            $scope.orgFilter[key] = show;
          });
        }

        $scope.$watch('orgFilter', function () {
          if (typeof $scope.stats != 'undefined') {
            var orgs = $.grep($scope.orgs, function (item) {
              if (typeof item.repoCount != 'undefined') {
                return $scope.orgFilter[item.Id];
              } else {
                return false
              }
            });

            if (orgs.length < 10) {
              $scope.commitGraphOptions.animation = true;
            } else {
              $scope.commitGraphOptions.animation = false;
            }

            $scope.countGraphData = {
              labels: $.map(orgs, function (r) { return r.Login }),
              datasets: [{
                fillColor: "#428BCA",
                strokeColor: "#2A6496",
                pointColor: "#428BCA",
                data: $.map(orgs, function (r) { return r.repoCount })
              }]
            };

            var org_stats = {};

            for (var i = $scope.stats.org_stats.length - 1; i >= 0; i--) {
              var stat = $scope.stats.org_stats[i];
              var org_logins = $.map(orgs, function (i) { return i.Login });
              if (org_logins.indexOf(stat.OrganizationLogin) > -1) {
                if (typeof org_stats[stat.OrganizationLogin] === 'undefined') {
                  org_stats[stat.OrganizationLogin] = [];
                  org_stats[stat.OrganizationLogin][11] = 0;
                }
                var mon_num = $scope.months.indexOf(stat.Month);
                org_stats[stat.OrganizationLogin][mon_num] = stat;
              }
            };


            var datasets = [];
            $.each(org_stats, function (org, stat) {
              var commits;
              commits = $.map(stat, function (el) {
                if (typeof el !== 'undefined') {
                  return el.Commits;
                }
              });
              var color = colors[datasets.length - (Math.floor(datasets.length / colors.length) * colors.length)];
              var obj = {
                label: org,
                strokeColor: color,
                pointColor: color,
                data: commits
              }
              datasets.push(obj);
            });

            $scope.commitGraphData = {
              labels: $scope.months,
              datasets: datasets
            }

          }

        }, true);


        $http.get($rootScope.apiUrl + '/orgs', { cache: true }).success(function (data) {
          $scope.orgs = data;
          $http.get($rootScope.apiUrl + '/stats').success(function (statsData) {
            for (var i = statsData.repo_counts.length - 1; i >= 0; i--) {
              var org;
              org = $.grep($scope.orgs, function (item) {
                return item.Login === statsData.repo_counts[i].OrganizationLogin;
              });
              org[0].repoCount = statsData.repo_counts[i].RepoCount;
              if (i < 8) {
                $scope.orgFilter[org[0].Id] = true;
              }
            };

            $scope.stats = statsData;
            $scope.months = $.map(statsData.org_stats, function (s) { return s.Month })
                            .filter(function(itm,i,a){
                              return i==a.indexOf(itm);
                            });
            $scope.months = $.unique($scope.months);
          });

        });

  }]);
