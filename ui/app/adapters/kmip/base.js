import ApplicationAdapater from '../application';
import { encodePath } from 'vault/utils/path-encoding-helpers';

export default ApplicationAdapater.extend({
  namespace: 'v1',
  pathForType(type) {
    return type.replace('kmip/', '');
  },

  _url(modelType, meta = {}, id) {
    let { backend, scope, role } = meta;
    let type = this.pathForType(modelType);
    let base;
    switch (type) {
      case 'scope':
        base = `${encodePath(backend)}/scope`;
        break;
      case 'role':
        base = `${encodePath(backend)}/scope/${encodePath(scope)}/role`;
        break;
      case 'credential':
        base = `${encodePath(backend)}/scope/${encodePath(scope)}/role/${encodePath(role)}/credential`;
        break;
    }

    if (id) {
      return `/v1/${base}/${encodePath(id)}`;
    }
    return `/v1/${base}`;
  },

  urlForQuery(query, modelType) {
    let base = this._url(modelType, query);
    return base + '?list=true';
  },

  query(store, type, query) {
    return this.ajax(this.urlForQuery(query, type.modelName), 'GET');
  },

  queryRecord(store, type, query) {
    let id = query.id;
    delete query.id;
    return this.ajax(this._url(type.modelName, query, id), 'GET').then(resp => {
      resp.id = id;
      resp = { ...resp, ...query };
      return resp;
    });
  },
  buildURL(modelName, id, snapshot, requestType, query) {
    if (requestType === 'createRecord') {
      return this._super(...arguments);
    }
    return this._super(`${modelName}`, id, snapshot, requestType, query);
  },
});
