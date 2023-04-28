/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { create, visitable, fillable, clickable } from 'ember-cli-page-object';
import { settled } from '@ember/test-helpers';
import VAULT_KEYS from 'vault/tests/helpers/vault-keys';

const { rootToken } = VAULT_KEYS;

export default create({
  visit: visitable('/vault/auth'),
  logout: visitable('/vault/logout'),
  submit: clickable('[data-test-auth-submit]'),
  tokenInput: fillable('[data-test-token]'),
  usernameInput: fillable('[data-test-username]'),
  passwordInput: fillable('[data-test-password]'),
  namespaceInput: fillable('[data-test-auth-form-ns-input]'),
  login: async function (token) {
    // make sure we're always logged out and logged back in
    await this.logout();
    await settled();
    // clear session storage to ensure we have a clean state
    window.localStorage.clear();
    await this.visit({ with: 'token' });
    await settled();
    if (token) {
      await this.tokenInput(token).submit();
      return;
    }

    await this.tokenInput(rootToken).submit();
    return;
  },
  loginUsername: async function (username, password) {
    // make sure we're always logged out and logged back in
    await this.logout();
    await settled();
    // clear local storage to ensure we have a clean state
    window.localStorage.clear();
    await this.visit({ with: 'username' });
    await settled();
    await this.usernameInput(username);
    await this.passwordInput(password).submit();
    return;
  },
  loginNs: async function (ns) {
    // make sure we're always logged out and logged back in
    await this.logout();
    await settled();
    // clear session storage to ensure we have a clean state
    window.localStorage.clear();
    await this.visit({ with: 'token' });
    await settled();
    await this.namespaceInput(ns);
    await settled();
    await this.tokenInput(rootToken).submit();
    return;
  },
});
