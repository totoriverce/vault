/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { inject as service } from '@ember/service';
import { computed, set } from '@ember/object';
import Component from '@ember/component';

const MODEL_TYPES = {
  'ssh-sign': {
    model: 'ssh-sign',
  },
  'ssh-creds': {
    model: 'ssh-otp-credential',
    title: 'Generate SSH Credentials',
  },
  'aws-creds': {
    model: 'aws-credential',
    title: 'Generate AWS Credentials',
    backIsListLink: true,
  },
  'pki-issue': {
    model: 'pki/cert',
    title: 'Issue Certificate',
  },
  'pki-sign': {
    model: 'pki-certificate-sign',
    title: 'Sign Certificate',
  },
};

export default Component.extend({
  store: service(),
  router: service(),
  // set on the component
  backendType: null,
  backendPath: null,
  roleName: null,
  action: null,

  model: null,
  loading: false,
  emptyData: '{\n}',

  modelForType() {
    const type = this.options;
    if (type) {
      return type.model;
    }
    // if we don't have a mode for that type then redirect them back to the backend list
    this.router.transitionTo('vault.cluster.secrets.backend.list-root', this.backendPath);
  },

  options: computed('action', 'backendType', function () {
    const action = this.action || 'creds';
    return MODEL_TYPES[`${this.backendType}-${action}`];
  }),

  init() {
    this._super(...arguments);
    this.createOrReplaceModel();
  },

  willDestroy() {
    if (!this.model.isDestroyed && !this.model.isDestroying) {
      this.model.unloadRecord();
    }
    this._super(...arguments);
  },

  createOrReplaceModel() {
    const modelType = this.modelForType();
    const model = this.model;
    const roleName = this.roleName;
    const backendPath = this.backendPath;
    if (!modelType) {
      return;
    }
    if (model) {
      model.unloadRecord();
    }
    const attrs = {
      role: {
        backend: backendPath,
        name: roleName,
      },
      id: `${backendPath}-${roleName}`,
    };
    const newModel = this.store.createRecord(modelType, attrs);
    this.set('model', newModel);
  },

  actions: {
    create() {
      const model = this.model;
      this.set('loading', true);
      this.model.save().finally(() => {
        model.set('hasGenerated', true);
        this.set('loading', false);
      });
    },

    codemirrorUpdated(attr, val, codemirror) {
      codemirror.performLint();
      const hasErrors = codemirror.state.lint.marked.length > 0;

      if (!hasErrors) {
        set(this.model, attr, JSON.parse(val));
      }
    },

    newModel() {
      this.createOrReplaceModel();
    },
  },
});
