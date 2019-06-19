import BaseAdapter from './base';

export default BaseAdapter.extend({
  createRecord(store, type, snapshot) {
    let url = this._url(type.modelName, {
      backend: snapshot.record.backend,
      scope: snapshot.record.scope,
      role: snapshot.record.role,
    });
    url = `${url}/generate`;
    return this.ajax(url, 'POST', { data: snapshot.serialize() }).then(model => {
      // TODO change this to serial?
      return {
        ...model,
        id: model.serial || 'foo',
      };
    });
  },
});
