import Ember from 'ember';
import DS from 'ember-data';
import Resolver from 'ember/resolver';
import loadInitializers from 'ember/load-initializers';

Ember.MODEL_FACTORY_INJECTIONS = true;

var App = Ember.Application.extend({
  modulePrefix: 'text-gravy', // TODO: loaded via config
  Resolver: Resolver,
  ApplicationAdapter: DS.RESTAdapter.extend({
    namespace: 'api/v1' })
});

loadInitializers(App, 'text-gravy');

export default App;
