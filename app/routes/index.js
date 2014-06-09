import Ember from 'ember';

var IndexRoute = Ember.Route.extend({
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

export default IndexRoute;