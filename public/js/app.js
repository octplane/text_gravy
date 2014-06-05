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

  App.Tag = DS.Model.extend({
    text: attr
  });

  App.PhotoInfo = DS.Model.extend({
    title: attr,
    description : attr,
    tags: DS.hasMany('tag')
  })

  App.Router.map(function() {
    this.route('index', { path: '/' });
    this.resource('index', {path: '/search/:query'})
    this.resource('photo', {path: '/photo/:photo_id'})

  });

  // App.Router.reopen({
  //   location: 'history'
  // });

  App.PhotoRoute = Ember.Route.extend({
    model: function(params) {
      return {
        photo: this.store.find('photo', params.photo_id),
        info: this.store.find('photo_info', params.photo_id)
      };
    }
  });

  App.IndexRoute = Ember.Route.extend({
    beforeModel: function(transition) {
      if (!transition.params.index.query) {
        this.transitionTo('/search/meditation');
      }
    },

    model: function(params) {
      var sParam = params.query;

      return {query: sParam, result: this.store.find('photo',{ query: sParam })};
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