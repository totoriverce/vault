/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import buildRoutes from 'ember-engines/routes';

export default buildRoutes(function () {
  // There are two list routes because Ember won't let a query param (e.g. *path_to_secret) be blank.
  // We only have a query param when we're on a nested secret route like beep/boop/
  this.route('list');
  this.route('list-directory', { path: '/:path_to_secret/directory' });
  this.route('create');
  this.route('secret', { path: '/:name' }, function () {
    this.route('details');
    this.route('edit');
    this.route('metadata', function () {
      this.route('edit');
      this.route('versions');
      this.route('diff');
    });
  });
  this.route('configuration');
});
