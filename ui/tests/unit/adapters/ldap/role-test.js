/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';
import { setupMirage } from 'ember-cli-mirage/test-support';

module('Unit | Adapter | ldap/role', function (hooks) {
  setupTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.store = this.owner.lookup('service:store');
    this.adapter = this.store.adapterFor('ldap/role');
    this.path = 'role';
  });

  test('it should make request to correct endpoints when listing records', async function (assert) {
    assert.expect(4);

    this.server.get('/ldap-test/:path', (schema, req) => {
      assert.ok(req.queryParams.list, 'list query param sent when listing roles');
      assert.strictEqual(
        req.params.path,
        this.path,
        'GET request made to correct endpoint when listing roles'
      );
      return { data: { keys: ['test-role'] } };
    });

    await this.store.query('ldap/role', { backend: 'ldap-test', type: 'dynamic' });
    this.path = 'static-role';
    await this.store.query('ldap/role', { backend: 'ldap-test', type: 'static' });
  });

  test('it should make request to correct endpoints when querying record', async function (assert) {
    assert.expect(2);

    this.server.get('/ldap-test/:path/test-role', () => {
      assert.ok('GET request made to correct endpoint when querying record');
      return { data: {} };
    });

    for (const type of ['dynamic', 'static']) {
      await this.store.queryRecord('ldap/role', { backend: 'ldap-test', type, name: 'test-role' });
      this.path = 'static-role';
    }
  });

  test('it should make request to correct endpoints when creating new record', async function (assert) {
    assert.expect(2);

    this.server.post('/ldap-test/:path/test-role', () => {
      assert.ok('POST request made to correct endpoint when creating new record');
    });

    const getModel = (type) => {
      return this.store.createRecord('ldap/role', {
        backend: 'ldap-test',
        name: 'test-role',
        type,
      });
    };

    for (const type of ['dynamic', 'static']) {
      const model = getModel(type);
      await model.save();
    }
  });

  test('it should make request to correct endpoints when updating record', async function (assert) {
    assert.expect(2);

    this.server.post('/ldap-test/:path/test-role', () => {
      assert.ok('POST request made to correct endpoint when updating record');
    });

    this.store.pushPayload('ldap/role', {
      modelName: 'ldap/role',
      backend: 'ldap-test',
      name: 'test-role',
    });
    const record = this.store.peekRecord('ldap/role', 'test-role');

    for (const type of ['dynamic', 'static']) {
      record.type = type;
      await record.save();
      this.path = 'static-role';
    }
  });

  test('it should make request to correct endpoints when deleting record', async function (assert) {
    assert.expect(2);

    this.server.delete('/ldap-test/:path/test-role', () => {
      assert.ok('DELETE request made to correct endpoint when deleting record');
    });

    const getModel = () => {
      this.store.pushPayload('ldap/role', {
        modelName: 'ldap/role',
        backend: 'ldap-test',
        name: 'test-role',
      });
      return this.store.peekRecord('ldap/role', 'test-role');
    };

    for (const type of ['dynamic', 'static']) {
      const record = getModel();
      record.type = type;
      await record.destroyRecord();
      this.path = 'static-role';
    }
  });

  test('it should make request to correct endpoints when fetching credentials', async function (assert) {
    assert.expect(2);

    this.path = 'creds';

    this.server.get('/ldap-test/:path/test-role', () => {
      assert.ok('GET request made to correct endpoint when fetching credentials');
    });

    for (const type of ['dynamic', 'static']) {
      await this.adapter.fetchCredentials('ldap-test', type, 'test-role');
      this.path = 'static-cred';
    }
  });

  test('it should make request to correct endpoint when rotating static role password', async function (assert) {
    assert.expect(1);

    this.server.post('/ldap-test/rotate-role/test-role', () => {
      assert.ok('GET request made to correct endpoint when rotating static role password');
    });

    await this.adapter.rotateStaticPassword('ldap-test', 'test-role');
  });
});
