/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { service } from '@ember/service';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import errorMessage from 'vault/utils/error-message';

import type FlashMessageService from 'vault/services/flash-messages';
import type StoreService from 'vault/services/store';
import type RouterService from '@ember/routing/router-service';

interface Args {
  onClose: () => void;
  onError: (errorMessage: string) => void;
}

export default class SyncActivationModal extends Component<Args> {
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly store: StoreService;
  @service declare readonly router: RouterService;

  @tracked hasConfirmedDocs = false;

  @task
  @waitFor
  *onFeatureConfirm() {
    try {
      yield this.store
        .adapterFor('application')
        .ajax('/v1/sys/activation-flags/secrets-sync/activate', 'POST');
      this.router.transitionTo('vault.cluster.sync.secrets.overview');
    } catch (error) {
      this.args.onError(errorMessage(error));
      this.flashMessages.danger(`Error enabling feature \n ${errorMessage(error)}`);
    } finally {
      this.args.onClose();
    }
  }
}
