(function() {

  var NS = 'api/v1';
  var App = Ember.Application.create();

  App.ApplicationAdapter = DS.RESTAdapter.extend({
    namespace: NS
  });

  var attr = DS.attr('string');

  App.Kitten = DS.Model.extend({
    name: attr,
    picture: attr
  });

  App.Router.map(function() {
    this.resource('index', {path: '/search/:query'})
  })


  App.IndexRoute = Ember.Route.extend({
    model: function(params) {
      debugger;
      return {query: params.query, result: this.store.find('kitten',{ query: params.query })};
    },
    actions: {
      search: function() {
        debugger;
        this.controller.transitionToRoute('index', {query: this.controller.get("q")});
      }
    }
  });


})();