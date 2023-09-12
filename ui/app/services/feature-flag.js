/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Service from '@ember/service';

const FLAGS = {
  vaultCloudNamespace: 'VAULT_CLOUD_ADMIN_NAMESPACE',
};

export default Service.extend({
  featureFlags: null,
  setFeatureFlags(flags) {
    this.set('featureFlags', flags);
  },

  get managedNamespaceRoot() {
    if (this.featureFlags && this.featureFlags.includes(FLAGS.vaultCloudNamespace)) {
      return 'admin';
    }
    return null;
  },
});
