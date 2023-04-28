/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { withConfirmLeave } from 'core/decorators/confirm-leave';

@withConfirmLeave('model.config', ['model.urls', 'model.crl'])
export default class PkiConfigurationEditRoute extends Route {
  @service secretMountPath;

  model() {
    const { urls, crl, engine } = this.modelFor('configuration');
    return {
      engineId: engine.id,
      urls,
      crl,
    };
  }

  setupController(controller, resolvedModel) {
    super.setupController(controller, resolvedModel);
    controller.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: this.secretMountPath.currentPath, route: 'overview' },
      { label: 'configuration', route: 'configuration.index' },
      { label: 'edit' },
    ];
  }
}
