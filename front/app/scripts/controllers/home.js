
angular.module('govcodeApp')
  .controller('HomeCtrl', ['$rootScope', 
                           '$scope',
                           '$http',
                           '$routeParams',
                           function ($rootScope, $scope, $http, $routeParams) {

    $scope.langs = [{ text: 'JavaScript', weight: 1052},
{ text: 'C++', weight: 1050},
{ text: 'Python', weight: 1035},
{ text: 'CSS', weight: 1035},
{ text: 'Java', weight: 402},
{ text: 'PHP', weight: 123},
{ text: 'Shell', weight: 72},
{ text: 'Objective-C', weight: 53}];
}]);