/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';
import { service } from '@ember/service';
import { task } from 'ember-concurrency';

export default Controller.extend({
  router: service(),
  showRoute: 'vault.cluster.access.identity.show',
  showTab: 'details',
  navAfterSave: task(function* ({ saveType, model }) {
    const isDelete = saveType === 'delete';
    const type = model.get('identityType');
    const listRoutes = {
      'entity-alias': 'vault.cluster.access.identity.aliases.index',
      'group-alias': 'vault.cluster.access.identity.aliases.index',
      group: 'vault.cluster.access.identity.index',
      entity: 'vault.cluster.access.identity.index',
    };
    const routeName = listRoutes[type];
    if (!isDelete) {
      yield this.router.transitionTo(this.showRoute, model.id, this.showTab);
      return;
    }
    yield this.router.transitionTo(routeName);
  }),
});
