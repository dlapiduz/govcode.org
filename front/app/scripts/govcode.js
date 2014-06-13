var gcControllers = angular.module('gcControllers', []);

gcControllers.controller('MainCtrl', ['$scope', '$location',
  function ($scope, $location) {
    $scope.getClass = function (path) {
      if (($location.path().substr(0, path.length) == path && path != "/") || $location.path() == path) {
        return "current";
      } else {
        return "";
      }
    };

  }]
);

gcControllers.controller('ReposCtrl', ['$scope', '$http',
  function ($scope, $http) {
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

  }]
);


gcControllers.controller('ContribCtrl', function () {});
gcControllers.controller('StatsCtrl', function () {});
gcControllers.controller('AboutCtrl', function () {});

var govcode = angular.module('govcode', ['ngRoute', 'gcControllers', 'truncate', 'ui.unique']);

govcode.filter('multifilter', function() {
    return function(items, options) {
      if (items === undefined) { return };
      var filtered_items = [];
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

        filtered_items.push(items[i]);
      };

      return filtered_items;
    };
});

govcode.config(['$routeProvider',
  function($routeProvider) {
    $routeProvider.
      when('/', {
        templateUrl: 'partials/repos.html',
        controller: 'ReposCtrl'
      }).
      when('/contributors', {
        templateUrl: 'partials/contributors.html',
        controller: 'ContribCtrl'
      }).
      when('/stats', {
        templateUrl: 'partials/stats.html',
        controller: 'StatsCtrl'
      }).
      when('/about', {
        templateUrl: 'partials/about.html',
        controller: 'AboutCtrl'
      }).
      otherwise({
        redirectTo: '/'
      });
}]);
