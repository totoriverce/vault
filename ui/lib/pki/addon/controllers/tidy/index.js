/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */
import Ember from 'ember';
import Controller from '@ember/controller';
import { task, timeout } from 'ember-concurrency';
import { tracked } from '@glimmer/tracking';
import { inject as service } from '@ember/service';

const POLL_INTERVAL_MS = 2000;

export default class PkiTidyIndexController extends Controller {
  @service store;
  @service secretMountPath;

  @tracked tidyStatus = null;

  @task
  *pollTidyStatus() {
    while (true) {
      // when testing, the polling loop causes promises to never settle so acceptance tests hang
      // to get around that, we just disable the poll in tests
      if (Ember.testing) {
        return;
      }
      yield timeout(POLL_INTERVAL_MS);
      try {
        const tidyStatusResponse = yield this.fetchTidyStatus();
        this.tidyStatus = tidyStatusResponse;
      } catch (e) {
        // we want to keep polling here
      }
    }
  }
}
