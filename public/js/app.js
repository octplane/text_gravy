(function() {

  var NS = 'api/v1';
  var App = Ember.Application.create();


  App.ApplicationAdapter = DS.RESTAdapter.extend({
    namespace: NS
  });

  var attr = DS.attr('string');

  App.Photo = DS.Model.extend({
    title: attr,
    thumb: attr,
    large: attr
  });

  App.Router.map(function() {
    this.resource('index', {path: '/search/:query'})
    this.resource('photo', {path: '/photo/:photo_id'})

  })

  App.PhotoRoute = Ember.Route.extend({
    model: function(params) {
      return {photo: this.store.find('photo', params.photo_id)};
    }
  });


  App.IndexRoute = Ember.Route.extend({
    model: function(params) {
      return {query: params.query, result: this.store.find('photo',{ query: params.query })};
    },
    actions: {
      search: function() {
        this.controller.transitionToRoute('/search/'+ this.controller.get("q"));
      }
    }
  });

  // Define child view
  App.PhotoInfoView = Ember.View.extend({
    templateName: 'photo-info'
  });


})();