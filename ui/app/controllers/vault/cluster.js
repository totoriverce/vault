/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

/* eslint-disable ember/no-observers */
import { inject as service } from '@ember/service';
import { alias } from '@ember/object/computed';
import Controller from '@ember/controller';
import { observer, computed } from '@ember/object';
export default Controller.extend({
  auth: service(),
  store: service(),
  media: service(),
  router: service(),
  permissions: service(),
  namespaceService: service('namespace'),
  flashMessages: service(),

  vaultVersion: service('version'),
  console: service(),

  queryParams: [
    {
      namespaceQueryParam: {
        scope: 'controller',
        as: 'namespace',
      },
    },
  ],

  namespaceQueryParam: '',

  onQPChange: observer('namespaceQueryParam', function () {
    this.namespaceService.setNamespace(this.namespaceQueryParam);
  }),

  consoleOpen: alias('console.isOpen'),

  activeCluster: computed('auth.activeCluster', function () {
    return this.store.peekRecord('cluster', this.auth.activeCluster);
  }),

  activeClusterName: computed('activeCluster', function () {
    const activeCluster = this.activeCluster;
    return activeCluster ? activeCluster.get('name') : null;
  }),

  showNav: computed(
    'router.currentRouteName',
    'activeClusterName',
    'auth.currentToken',
    'activeCluster.{dr.isSecondary,needsInit,sealed}',
    function () {
      if (this.activeCluster.dr?.isSecondary || this.activeCluster.needsInit || this.activeCluster.sealed) {
        return false;
      }
      if (
        this.activeClusterName &&
        this.auth.currentToken &&
        this.router.currentRouteName !== 'vault.cluster.auth'
      ) {
        return true;
      }
      return;
    }
  ),

  actions: {
    toggleConsole() {
      this.toggleProperty('consoleOpen');
    },
  },
});
