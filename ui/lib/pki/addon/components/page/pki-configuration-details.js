/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { action } from '@ember/object';
import Component from '@glimmer/component';
// import RouterService from '@ember/routing/router-service';
// import FlashMessageService from 'vault/services/flash-messages';
import { inject as service } from '@ember/service';
import errorMessage from 'vault/utils/error-message';
import { tracked } from '@glimmer/tracking';

export default class PkiConfigurationDetails extends Component {
  @service store;
  @service router;
  @service flashMessages;

  @tracked showDeleteAllIssuers = false;

  @action
  async deleteAllIssuers() {
    try {
      this.store.adapterFor('pki/issuer').deleteAllIssuers(this.args.currentPath);
      this.flashMessages.success('Issuers and keys deleted successfully.');
      this.showDeleteAllIssuers = false;
      this.router.transitionTo('vault.cluster.secrets.backend.pki.configuration');
    } catch (error) {
      this.flashMessages.danger(errorMessage(error));
    }
  }
}
