import Ember from 'ember';
const { computed, Component, inject } = Ember;

export function linkParams({ mode, secret, queryParams }) {
  let params;
  const route = `vault.cluster.secrets.backend.${mode}`;

  if (!secret || secret === ' ') {
    params = [route + '-root'];
  } else {
    params = [route, secret];
  }

  if (queryParams) {
    params.push(queryParams);
  }

  return params;
}

export default Component.extend({
  tagName: '',
  namespace: inject.service(),
  mode: 'list',

  secret: null,
  queryParams: null,
  ariaLabel: null,

  linkParams: computed('namespace.path', 'mode', 'secret', 'queryParams', function() {
    let namespace = this.get('namespace.path');
    let data = this.getProperties('mode', 'secret', 'queryParams');
    if (namespace) {
      data.queryParams = { ...{ namespace }, ...(data.queryParams || {}) };
    }
    return linkParams(data);
  }),
});
