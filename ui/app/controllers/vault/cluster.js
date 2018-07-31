import Ember from 'ember';

const { Controller, computed, inject } = Ember;
export default Controller.extend({
  auth: inject.service(),
  store: inject.service(),
  media: inject.service(),

  vaultVersion: inject.service('version'),
  console: inject.service(),

  queryParams: [{ namespaceQueryParam: 'namespace' }],
  namespaceQueryParam: '',

  consoleOpen: computed.alias('console.isOpen'),
  activeCluster: computed('auth.activeCluster', function() {
    return this.get('store').peekRecord('cluster', this.get('auth.activeCluster'));
  }),
  activeClusterName: computed('activeCluster', function() {
    const activeCluster = this.get('activeCluster');
    return activeCluster ? activeCluster.get('name') : null;
  }),
  showNav: computed(
    'activeClusterName',
    'auth.currentToken',
    'activeCluster.dr.isSecondary',
    'activeCluster.{needsInit,sealed}',
    function() {
      if (
        this.get('activeCluster.dr.isSecondary') ||
        this.get('activeCluster.needsInit') ||
        this.get('activeCluster.sealed')
      ) {
        return false;
      }
      if (this.get('activeClusterName') && this.get('auth.currentToken')) {
        return true;
      }
    }
  ),
  actions: {
    toggleConsole() {
      this.toggleProperty('consoleOpen');
    },
  },
});
