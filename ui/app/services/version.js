/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Service, { inject as service } from '@ember/service';
import { keepLatestTask, task } from 'ember-concurrency';
import { tracked } from '@glimmer/tracking';
import { buildWaiter } from '@ember/test-waiters';

const waiter = buildWaiter('version');

export default class VersionService extends Service {
  @service store;
  @tracked features = [];
  @tracked version = null;
  @tracked type = null;

  get isEnterprise() {
    return this.type === 'enterprise';
  }

  get isCommunity() {
    return !this.isEnterprise;
  }

  /* Features */
  get hasPerfReplication() {
    return this.features.includes('Performance Replication');
  }

  get hasDRReplication() {
    return this.features.includes('DR Replication');
  }

  get hasSentinel() {
    return this.features.includes('Sentinel');
  }

  get hasNamespaces() {
    return this.features.includes('Namespaces');
  }

  get hasControlGroups() {
    return this.features.includes('Control Groups');
  }

  get hasSecretsSync() {
    return this.features.includes('Secrets Sync');
  }

  get versionDisplay() {
    if (!this.version) {
      return '';
    }
    return this.isEnterprise ? `v${this.version.slice(0, this.version.indexOf('+'))}` : `v${this.version}`;
  }

  @task({ drop: true })
  *getVersion() {
    if (this.version) return;
    const response = yield this.store.adapterFor('cluster').fetchVersion();
    this.version = response.data?.version;
  }

  @task
  *getType() {
    if (this.type !== null) return;
    const response = yield this.store.adapterFor('cluster').health();
    if (response.has_chroot_namespace) {
      // chroot_namespace feature is only available in enterprise
      this.type = 'enterprise';
      return;
    }
    this.type = response.enterprise ? 'enterprise' : 'community';
  }

  @keepLatestTask
  *getFeatures() {
    const waiterToken = waiter.beginAsync();
    if (this.features?.length || this.isCommunity) {
      waiter.endAsync(waiterToken);
      return;
    }
    try {
      const response = yield this.store.adapterFor('cluster').features();
      this.features = response.features;
      return;
    } catch (err) {
      // if we fail here, we're likely in DR Secondary mode and don't need to worry about it
    } finally {
      waiter.endAsync(waiterToken);
    }
  }

  fetchVersion() {
    return this.getVersion.perform();
  }

  fetchType() {
    return this.getType.perform();
  }

  fetchFeatures() {
    return this.getFeatures.perform();
  }
}
