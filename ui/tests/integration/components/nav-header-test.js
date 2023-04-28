/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { create } from 'ember-cli-page-object';
import hbs from 'htmlbars-inline-precompile';
import navHeader from 'vault/tests/pages/components/nav-header';

const component = create(navHeader);

module('Integration | Component | nav header', function (hooks) {
  setupRenderingTest(hooks);

  test('it renders', async function (assert) {
    await render(hbs`
        <div id="modal-wormhole"></div>
        {{#nav-header as |h|}}
          {{#h.home}}
            Home!
          {{/h.home}}
          {{#h.items}}
            Some Items
          {{/h.items}}
          {{#h.main}}
            Main stuff here
          {{/h.main}}
        {{/nav-header}}
      `);

    assert.ok(component.ele, 'renders the outer element');
    assert.strictEqual(component.homeText.trim(), 'Home!', 'renders home content');
    assert.strictEqual(component.itemsText.trim(), 'Some Items', 'renders items content');
    assert.strictEqual(component.mainText.trim(), 'Main stuff here', 'renders items content');
  });
});
