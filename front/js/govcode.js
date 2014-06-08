var gcControllers = angular.module('gcControllers', []);

gcControllers.controller('MainCtrl', ['$scope', '$location',
  function ($scope, $location) {
    $scope.getClass = function(path) {
      if ($location.path().substr(0, path.length) == path) {
        return "active"
      } else {
        return ""
      }
    };

    $scope.include = function(actual, expected) {
        $.each(expected, function(i, el) {
            if (actual == el) {
                console.log("True");
                return true;
            }
        });
        return false;
    }
  }]
);

gcControllers.controller('ReposCtrl', ['$scope', '$http',
  function ($scope, $http) {
    $http.get('http://localhost:3000/repos').success(function(data) {
      $scope.repos = data;
    });

    $http.get('http://localhost:3000/orgs').success(function(data) {
      $scope.orgs = data;
      var ids = $.map(data, function(el, i) { return el.Id });
      $scope.search = { OrganizationId: ids };
    });
  }]
);

var govcode = angular.module('govcode', ['ngRoute', 'gcControllers', 'truncate', 'ui.unique']);

govcode.filter('includefilter', function() {
    return function(items, options) {
        console.log(items);
        console.log(options)
      return items;
    };
});

govcode.config(['$routeProvider',
  function($routeProvider) {
    $routeProvider.
      when('/repos', {
        templateUrl: 'partials/repos.html',
        controller: 'ReposCtrl'
      }).
      otherwise({
        redirectTo: '/repos'
      });
}]);
