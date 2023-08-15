/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import getStorage from 'vault/lib/token-storage';
import timestamp from 'core/utils/timestamp';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import { tracked } from '@glimmer/tracking';
import { inject as service } from '@ember/service';

/**
 * @module DashboardClientCountCard
 * DashboardClientCountCard component are used to display total and new client count information
 *
 * @example
 * ```js
 * <Dashboard::ClientCountCard @license={{@model.license}}  />
 * ```
 */

export default class DashboardClientCountCard extends Component {
  currentDate = timestamp.now().toISOString();

  @service store;

  @tracked startDate = null;
  @tracked activityData = null;
  @tracked clientConfig = null;
  @tracked updatedAt = null;

  constructor() {
    super(...arguments);
    this.fetchClientActivity.perform();
    this.clientConfig = this.store.queryRecord('clients/config', {}).catch(() => {});
  }

  get currentMonthActivityTotalCount() {
    return this.activityData?.byMonth?.lastObject?.new_clients.clients;
  }

  get licenseStartTime() {
    return this.args.license.startTime || getStorage().getItem('vault:ui-inputted-start-date') || null;
  }

  @task
  @waitFor
  *getActivity(start_time) {
    // on init ONLY make network request if we have a start_time
    return start_time
      ? yield this.store.queryRecord('clients/activity', {
          start_time: { timestamp: start_time },
          end_time: { timestamp: this.currentDate },
        })
      : {};
  }

  @task
  @waitFor
  *fetchClientActivity() {
    try {
      this.startDate = this.licenseStartTime;
      this.activityData = yield this.getActivity.perform(this.startDate);
      this.updatedAt = timestamp.now().toISOString();
      this.noActivityData = this.activityData.activity.id === 'no-data' ? true : false;
    } catch (error) {
      this.error = error;
    }
  }
}
