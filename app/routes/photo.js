import Ember from 'ember';

var PhotoRoute = Ember.Route.extend({
  model: function(params) {
    return {
      photo: this.store.find('photo', params.photo_id),
      info: this.store.find('photoinfo', params.photo_id)
    };
  }
});

export default PhotoRoute;
