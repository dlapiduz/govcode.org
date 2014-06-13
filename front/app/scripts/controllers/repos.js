'use strict';

/**
 * @ngdoc function
 * @name govcodeApp.controller:ReposCtrl
 * @description
 * # ReposCtrl
 * Controller of the govcodeApp
 */
angular.module('govcodeApp')
  .controller('ReposCtrl', ['$scope', '$http',   function ($scope, $http) {
    $scope.search = {};
    $scope.search.orgFilter = {};
    $http.get('http://localhost:3000/repos').success(function (data) {
      $scope.repos = data;
      $scope.languages = $.unique($.map(data, function (el, i) { return el.Language }));
      $scope.search.Language = null;

    });

    $http.get('http://localhost:3000/orgs').success(function(data) {
      $scope.orgs = data;
      $.each(data, function(i, el) { 
        $scope.search.orgFilter[el.Id] = true;
      });
    });

  }]);
