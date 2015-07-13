angular.module('trucker',['ngRoute'])
  .config(function($routeProvider) {
    $routeProvider
      .when('/', {
        controller:'listController',
        templateUrl:'assets/list.html'
      })
      .when('/new', {
        controller:'formController',
        templateUrl:'assets/form.html'
      })
      .otherwise({
        redirectTo:'/'
      });
  })
  .controller("listController", function($scope, $http){
    // Simple GET request example :
    $http.get('/api/entries')
      .success(function(data, status, headers, config) {
        // this callback will be called asynchronously
        // when the response is available
        $scope.columns=["Broker", "Trip No.", "Date"]
        $scope.items=data
      })
      .error(function(data, status, headers, config) {
        // called asynchronously if an error occurs
        // or server returns response with an error status.
      });

  })
  .controller("formController", function($scope, $http, $location){
    $scope.update = function(user) {
      $http.post('/api/entries', user)
        .success(function(data, status, headers, config) {
          console.log("SUCCESS!")
          $location.path('/')

        })
        .error(function(data, status, headers, config) {

        })

    };
  })