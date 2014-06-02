(function() {

  var App = Ember.Application.create();

  App.ApplicationAdapter = DS.RESTAdapter.extend({
    namespace: 'api/v1'
  });

  var attr = DS.attr('string');

  App.Kitten = DS.Model.extend({
    name: attr,
    picture: attr
  });

  App.IndexRoute = Ember.Route.extend({
    model: function() {
      return this.store.find('kitten');
    }
  });


  // App.IndexRoute = Ember.Route.extend({
  //   model: function() {
  //     return ['red', 'yellow', 'blue'];
  //   }
  // });


})();