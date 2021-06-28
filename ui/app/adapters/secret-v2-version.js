/* eslint-disable */
import AdapterError from '@ember-data/adapter/error';
import { isEmpty } from '@ember/utils';
import { get } from '@ember/object';
import ApplicationAdapter from './application';
import { encodePath } from 'vault/utils/path-encoding-helpers';
import ControlGroupError from 'vault/lib/control-group-error';

export default ApplicationAdapter.extend({
  namespace: 'v1',
  _url(backend, id, infix = 'data') {
    let url = `${this.buildURL()}/${encodePath(backend)}/${infix}/`;
    if (!isEmpty(id)) {
      url = url + encodePath(id);
    }
    return url;
  },

  urlForFindRecord(id) {
    let [backend, path, version] = JSON.parse(id);
    let base = this._url(backend, path);
    return version ? base + `?version=${version}` : base;
  },

  urlForQueryRecord(id) {
    return this.urlForFindRecord(id);
  },

  findRecord() {
    return this._super(...arguments).catch(errorOrModel => {
      // if the response is a real 404 or if the secret is gated by a control group this will be an error,
      // otherwise the response will be the body of a deleted / destroyed version
      if (errorOrModel instanceof AdapterError) {
        throw errorOrModel;
      }
      return errorOrModel;
    });
  },

  queryRecord(id, options) {
    return this.ajax(this.urlForQueryRecord(id), 'GET', options).then(resp => {
      if (options.wrapTTL) {
        return resp;
      }
      resp.id = id;
      resp.backend = backend;
      return resp;
    });
  },

  urlForCreateRecord(modelName, snapshot) {
    let backend = snapshot.belongsTo('secret').belongsTo('engine').id;
    let path = snapshot.attr('path');
    return this._url(backend, path);
  },

  createRecord(store, modelName, snapshot) {
    let backend = snapshot.belongsTo('secret').belongsTo('engine').id;
    let path = snapshot.attr('path');
    return this._super(...arguments).then(resp => {
      resp.id = JSON.stringify([backend, path, resp.version]);
      return resp;
    });
  },

  urlForUpdateRecord(id) {
    let [backend, path] = JSON.parse(id);
    return this._url(backend, path);
  },

  v2DeleteOperation(store, id, deleteType = 'delete') {
    let [backend, path, version] = JSON.parse(id);
    // deleteType should be 'delete', 'destroy', 'undelete', 'delete-latest-version', 'destroy-version'
    if ((!version && deleteType === 'delete') || deleteType === 'delete-latest-version') {
      return this.ajax(this._url(backend, path, 'data'), 'DELETE')
        .then(() => {
          let model = store.peekRecord('secret-v2-version', id);
          return model && model.rollbackAttributes() && model.reload();
        })
        .catch(e => {
          return e;
        });
    } else {
      return this.ajax(this._url(backend, path, deleteType), 'POST', { data: { versions: [version] } })
        .then(() => {
          let model = store.peekRecord('secret-v2-version', id);
          // potential that model.reload() is never called.
          return model && model.rollbackAttributes() && model.reload();
        })
        .catch(e => {
          return e;
        });
    }
  },

  handleResponse(status, headers, payload, requestData) {
    // the body of the 404 will have some relevant information
    if (status === 404 && get(payload, 'data.metadata')) {
      return this._super(200, headers, payload, requestData);
    }
    return this._super(...arguments);
  },
});
