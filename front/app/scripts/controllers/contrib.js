'use strict';

/**
 * @ngdoc function
 * @name govcodeApp.controller:ContribCtrl
 * @description
 * # ContribCtrl
 * Controller of the govcodeApp
 */

angular.module('govcodeApp')
  .controller('ContribCtrl', ['$rootScope',
                              '$scope',
                              '$http',
                              '$filter',
                              function ($rootScope, $scope, $http, $filter) {

    $scope.config = {
      itemsPerPage: 100,
      fillLastPage: false
    }

    $scope.filteredUsers = [];



    $http.get($rootScope.apiUrl + '/users').success(function (data) {
      $scope.users = data;

      $.map($scope.users, function(el) {
        el.OrgList = el.OrgList.replace(/[\{\}]/g, "").split(',');
        return el;
      })

      $scope.updateFilteredList();
    });

    $scope.updateFilteredList = function() {
      $scope.filteredUsers = $filter("filter")($scope.users, $scope.query);
    };

  }]);
