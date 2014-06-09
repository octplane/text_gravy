import DS from 'ember-data';

export default DS.Model.extend({
  title: DS.attr('string'),
  thumb: DS.attr('string'),
  large: DS.attr('string'),
  description : DS.attr('string'),
  tags: DS.hasMany('tag'),
});