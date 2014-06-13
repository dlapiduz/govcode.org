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
    'ui.unique'
  ])
  .config(function ($routeProvider) {
    $routeProvider
      .when('/', {
        templateUrl: 'views/repos.html',
        controller: 'ReposCtrl'
      })
      .when('/contributors', {
        templateUrl: 'views/contributors.html',
        controller: 'ContribCtrl'
      })
      .when('/stats', {
        templateUrl: 'views/stats.html',
        controller: 'StatsCtrl'
      })
      .when('/about', {
        templateUrl: 'views/about.html',
        controller: 'AboutCtrl'
      })
      .otherwise({
        redirectTo: '/'
      });
  })
  .filter('multifilter', function() {
    return function(items, options) {
      if (items === undefined) { return; }
      var filteredItems = [];
      var i;
      for (i = 0; i < items.length; i++) {
        // Filter by language
        if (options.Language !== null && items[i].Language !== options.Language) {
          continue;
        }
        // Filter by org
        if (!options.orgFilter[items[i].OrganizationId]) {
          continue;
        }

        filteredItems.push(items[i]);
      }

      return filteredItems;
    };
});
