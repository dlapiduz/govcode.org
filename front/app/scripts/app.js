'use strict';

/**
 * @ngdoc overview
 * @name govcodeApp
 * @description
 * # govcodeApp
 *
 * Main module of the application.
 */
angular
  .module('govcodeApp', [
    'ngCookies',
    'ngResource',
    'ngRoute',
    'ngSanitize',
    'truncate', 
    'ui.unique',
    'angular-table',
    'chartjs',
    'angularMoment',
    'angular-google-analytics'
  ])
  .config(['$routeProvider',
           '$locationProvider',
           'AnalyticsProvider',
           function ($routeProvider, $locationProvider, AnalyticsProvider) {

    AnalyticsProvider.setAccount('UA-8480852-11');
    AnalyticsProvider.trackPages(true);

    $locationProvider.html5Mode(true);
    $routeProvider
      .when('/', {
        templateUrl: '/views/repos.html',
        controller: 'ReposCtrl'
      })
      .when('/repos/:repoName', {
        templateUrl: '/views/repo.html',
        controller: 'RepoCtrl'
      })
      .when('/contributors', {
        templateUrl: '/views/contributors.html',
        controller: 'ContribCtrl'
      })
      .when('/stats', {
        templateUrl: '/views/stats.html',
        controller: 'StatsCtrl'
      })
      .when('/about', {
        templateUrl: '/views/about.html',
        controller: 'AboutCtrl'
      })
      .otherwise({
        redirectTo: '/'
      });
  }])
  .run(['$rootScope', function($rootScope) {
    $rootScope.apiUrl = "http://localhost:3000";
  }])
  .filter('multifilter', function() {
    return function(items, options) {
      if (items === undefined) { return; }
      var filteredItems = [];
      var i;
      for (i = 0; i < items.length; i++) {
        // Filter by language
        if (options.Language !== '' && items[i].Language !== options.Language) {
          continue;
        }
        // Filter by org
        if (!options.orgFilter[items[i].OrganizationId]) {
          continue;
        }

        // Filter by lastActivity
        if (options.lastActivity < 12 && 
            options.lastActivity > 0 && 
            ( items[i].DaysSinceCommit > options.lastActivity * 30 || 
              items[i].DaysSinceCommit == 0)
            ) {
          continue;
        }


        filteredItems.push(items[i]);
      }

      return filteredItems;
    };
});
