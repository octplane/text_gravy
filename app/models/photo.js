import DS from 'ember-data';

export default DS.Model.extend({
  title: DS.attr(),
  thumb: DS.attr(),
  large: DS.attr(),
});