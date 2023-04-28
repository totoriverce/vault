/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { task } from 'ember-concurrency';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';

/**
 * @module PkiRoleForm
 * PkiRoleForm components are used to create and update PKI roles.
 *
 * @example
 * ```js
 * <PkiRoleForm @model={{this.model}}/>
 * ```
 * @callback onCancel
 * @callback onSave
 * @param {Object} role - pki/role model.
 * @param {Array} issuers - pki/issuer model.
 * @param {onCancel} onCancel - Callback triggered when cancel button is clicked.
 * @param {onSave} onSave - Callback triggered on save success.
 */

export default class PkiRoleForm extends Component {
  @service store;
  @service flashMessages;
  @service secretMountPath;

  @tracked errorBanner;
  @tracked invalidFormAlert;
  @tracked modelValidations;
  @tracked showDefaultIssuer = true;

  constructor() {
    super(...arguments);

    this.showDefaultIssuer = this.args.role.issuerRef === 'default';
  }

  get issuers() {
    return this.args.issuers?.map((issuer) => {
      return { issuerDisplayName: issuer.issuerName || issuer.issuerId };
    });
  }

  get breadcrumbs() {
    const crumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: this.secretMountPath.currentPath, route: 'overview' },
      { label: 'roles', route: 'roles.index' },
    ];
    if (!this.args.role.isNew) {
      crumbs.push({ label: this.args.role.id, route: 'roles.role.details' }, { label: 'edit' });
    }
    return crumbs;
  }

  @task
  *save(event) {
    event.preventDefault();
    try {
      const { isValid, state, invalidFormMessage } = this.args.role.validate();
      this.modelValidations = isValid ? null : state;
      this.invalidFormAlert = invalidFormMessage;
      if (isValid) {
        const { isNew, name } = this.args.role;
        yield this.args.role.save();
        this.flashMessages.success(`Successfully ${isNew ? 'created' : 'updated'} the role ${name}.`);
        this.args.onSave();
      }
    } catch (error) {
      const message = error.errors ? error.errors.join('. ') : error.message;
      this.errorBanner = message;
      this.invalidFormAlert = 'There was an error submitting this form.';
    }
  }

  @action
  toggleShowDefaultIssuer() {
    this.showDefaultIssuer = !this.showDefaultIssuer;

    if (this.showDefaultIssuer) {
      this.args.role.issuerRef = 'default';
    }
  }
}
