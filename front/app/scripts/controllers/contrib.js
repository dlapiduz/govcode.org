'use strict';

/**
 * @ngdoc function
 * @name govcodeApp.controller:ContribCtrl
 * @description
 * # ContribCtrl
 * Controller of the govcodeApp
 */

angular.module('govcodeApp')
  .controller('ContribCtrl', ['$scope', '$http', function ($scope, $http) {
    $http.get('http://localhost:3000/users').success(function (data) {
      $scope.users = data;
    });

  }]);
