'use strict';

/**
 * @ngdoc function
 * @name govcodeApp.controller:MainCtrl
 * @description
 * # MainCtrl
 * Controller of the govcodeApp
 */
angular.module('govcodeApp')
  .controller('MainCtrl',['$scope', '$location', function ($scope, $location) {
    $scope.getClass = function (path) {
      if (($location.path().substr(0, path.length) === path && path !== '/') || $location.path() === path) {
        return 'current';
      } else {
        return '';
      }
    };
  }]);
