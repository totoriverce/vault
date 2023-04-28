/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Application from 'vault/app';
import config from 'vault/config/environment';
import * as QUnit from 'qunit';
import { setApplication } from '@ember/test-helpers';
import { setup } from 'qunit-dom';
import { start } from 'ember-qunit';
import './helpers/flash-message';
import preloadAssets from 'ember-asset-loader/test-support/preload-assets';
import manifest from 'vault/config/asset-manifest';

preloadAssets(manifest).then(() => {
  setApplication(Application.create(config.APP));
  // TODO CBS: Check what this is, upgrade added it
  setup(QUnit.assert);
  start({
    setupTestIsolationValidation: true,
  });
});
