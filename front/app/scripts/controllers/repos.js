'use strict';

/**
 * @ngdoc function
 * @name govcodeApp.controller:ReposCtrl
 * @description
 * # ReposCtrl
 * Controller of the govcodeApp
 */

// Helper func to slugify langs
var mapLang = function (lang) {
  if (typeof lang === 'string') {
    return lang.toLowerCase().replace(/\+|\#/g, '');
  }
};
// Helper to get unique langs
var mangleLangs = function(data) {
  var langs = $.map(data, function (el, i) {
    if (el.Language !== "") {
      return el.Language;
    }
  });

  langs = langs.filter(function(el, index, arr) {
    return index == arr.indexOf(el);
  });

  return $.map(langs, function(el) { 
    return {
      name: el,
      slug: mapLang(el)
    }
  });
}


angular.module('govcodeApp')
  .controller('ReposCtrl', ['$scope', '$http', function ($scope, $http) {
    $scope.search = {};
    $scope.search.orgFilter = {};

    $scope.mapLang = mapLang;

    $scope.search.lastActivity = 6;

    $scope.sort = '-Forks';

    $scope.allOrgs = function(show) {
      $.each($scope.search.orgFilter, function(key, val) { 
        $scope.search.orgFilter[key] = show;
      });
    }
    
    // Get all repos
    $http.get('http://localhost:3000/repos').success(function (data) {
      // Load the repos in the scope
      $scope.repos = data;

      $scope.languages = mangleLangs(data);
      $scope.search.Language = "";

    });

    $http.get('http://localhost:3000/orgs').success(function(data) {
      $scope.orgs = data;
      $.each(data, function(i, el) { 
        $scope.search.orgFilter[el.Id] = true;
      });
    });

  }]);