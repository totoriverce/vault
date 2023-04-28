/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { inject as service } from '@ember/service';
import { alias } from '@ember/object/computed';
import Controller from '@ember/controller';
export default Controller.extend({
  namespaceService: service('namespace'),
  accessibleNamespaces: alias('namespaceService.accessibleNamespaces'),
  currentNamespace: alias('namespaceService.path'),
  actions: {
    refreshNamespaceList() {
      // fetch new namespaces for the namespace picker
      this.namespaceService.findNamespacesForUser.perform();
      this.send('reload');
    },
  },
});
