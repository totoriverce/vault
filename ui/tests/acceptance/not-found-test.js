/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { findAll, visit } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import authPage from 'vault/tests/pages/auth';
import logout from 'vault/tests/pages/logout';
import Ember from 'ember';

let adapterException;

module('Acceptance | not-found', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function () {
    adapterException = Ember.Test.adapter.exception;
    Ember.Test.adapter.exception = () => {};
    return authPage.login();
  });

  hooks.afterEach(function () {
    Ember.Test.adapter.exception = adapterException;
    return logout.visit();
  });

  test('top-level not-found', async function (assert) {
    await visit('/404');
    assert.ok(findAll('[data-test-not-found]').length, 'renders the not found component');
    assert.ok(
      findAll('[data-test-header-without-nav]').length,
      'renders the not found component with a header'
    );
  });

  test('vault route not-found', async function (assert) {
    await visit('/vault/404');
    assert.dom('[data-test-not-found]').exists('renders the not found component');
    assert.ok(findAll('[data-test-header-with-nav]').length, 'renders header with nav');
  });

  test('cluster route not-found', async function (assert) {
    await visit('/vault/secrets/secret/404/show');
    assert.dom('[data-test-not-found]').exists('renders the not found component');
    assert.ok(findAll('[data-test-header-with-nav]').length, 'renders header with nav');
  });

  test('secret not-found', async function (assert) {
    await visit('/vault/secrets/secret/show/404');
    assert.dom('[data-test-secret-not-found]').exists('renders the message about the secret not being found');
  });
});
