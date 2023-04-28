/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { run } from '@ember/runloop';
import hbs from 'htmlbars-inline-precompile';
import Service from '@ember/service';

const Permissions = Service.extend({
  globPaths: null,
  hasNavPermission() {
    return this.globPaths ? true : false;
  },
});

module('Integration | Helper | has-permission', function (hooks) {
  setupRenderingTest(hooks);
  hooks.beforeEach(function () {
    this.owner.register('service:permissions', Permissions);
    this.permissions = this.owner.lookup('service:permissions');
  });

  test('it renders', async function (assert) {
    await render(hbs`{{#if (has-permission)}}Yes{{else}}No{{/if}}`);

    assert.dom(this.element).hasText('No');
    await run(() => {
      this.permissions.set('globPaths', { 'test/': { capabilities: ['update'] } });
    });
    assert.dom(this.element).hasText('Yes', 'the helper re-computes when globPaths changes');
  });
});
