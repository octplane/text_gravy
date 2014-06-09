import Ember from 'ember';

var Router = Ember.Router.extend({
  location: TextGravyENV.locationType
});

Router.map(function() {
  this.route('index', { path: '/' });
  this.resource('index', {path: '/search/:query'});
  this.resource('photo', {path: '/photo/:photo_id'});
});

export default Router;
