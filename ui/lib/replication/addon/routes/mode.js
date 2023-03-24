/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { inject as service } from '@ember/service';
import { hash } from 'rsvp';
import { setProperties } from '@ember/object';
import Route from '@ember/routing/route';

const SUPPORTED_REPLICATION_MODES = ['dr', 'performance'];

export default Route.extend({
  replicationMode: service(),
  store: service(),
  beforeModel() {
    const replicationMode = this.paramsFor(this.routeName).replication_mode;
    if (!SUPPORTED_REPLICATION_MODES.includes(replicationMode)) {
      return this.transitionTo('index');
    }
  },
  model() {
    const replicationMode = this.paramsFor(this.routeName).replication_mode;
    this.replicationMode.setMode(replicationMode);
    return this.modelFor('application');
  },
  afterModel(model) {
    return hash({
      // set new property on model to compare if the drMode changes when you are demoting the cluster
      drModeInit: model.drMode,
    }).then(({ drModeInit }) => {
      setProperties(model, {
        drModeInit,
      });
      return model;
    });
  },
});
