/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { withConfirmLeave } from 'core/decorators/confirm-leave';

@withConfirmLeave('model.tidy')
export default class PkiConfigurationTidyRoute extends Route {
  @service store;
  @service secretMountPath;

  model() {
    return this.store.createRecord('pki/tidy', { backend: this.secretMountPath.currentPath });
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: this.secretMountPath.currentPath, route: 'overview' },
      { label: 'configuration', route: 'configuration.index' },
      { label: 'tidy' },
    ];
  }
}
