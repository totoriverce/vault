import Ember from 'ember';
import ApplicationAdapter from './application';
import DS from 'ember-data';

export default ApplicationAdapter.extend({
  pathForType() {
    return 'namespace';
  },
  urlForFindAll() {
    return `/${this.urlPrefix()}/namespaces?list=true`;
  },
  urlForCreateRecord(modelName, snapshot) {
    let id = snapshot.attr('path');
    return this.buildURL(modelName, id);
  },

  createRecord(store, type, snapshot) {
    let id = snapshot.attr('path');
    return this._super(...arguments).then(() => {
      return { id };
    });
  },
});
