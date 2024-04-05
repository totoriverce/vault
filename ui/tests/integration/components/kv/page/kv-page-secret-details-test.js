/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { click, find, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { kvDataPath, kvMetadataPath } from 'vault/utils/kv-path';
import { allowAllCapabilitiesStub } from 'vault/tests/helpers/stubs';
import { syncStatusResponse } from 'vault/mirage/handlers/sync';
import { KV_FORM, parseJsonEditor } from 'vault/tests/helpers/kv/kv-selectors';
import { KV_SECRET } from 'vault/tests/helpers/components/kv/page/secret/details-selectors';
import { GENERAL } from 'vault/tests/helpers/general-selectors';

module('Integration | Component | kv-v2 | Page::Secret::Details', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'kv');
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    this.store = this.owner.lookup('service:store');
    this.server.post('/sys/capabilities-self', allowAllCapabilitiesStub());
    this.backend = 'kv-engine';
    this.path = 'my-secret';
    this.pathComplex = 'my-secret-object';
    this.version = 2;
    this.dataId = kvDataPath(this.backend, this.path);
    this.dataIdComplex = kvDataPath(this.backend, this.pathComplex);
    this.metadataId = kvMetadataPath(this.backend, this.path);

    this.secretData = { foo: 'bar' };
    this.store.pushPayload('kv/data', {
      modelName: 'kv/data',
      id: this.dataId,
      secret_data: this.secretData,
      created_time: '2023-07-20T02:12:17.379762Z',
      custom_metadata: null,
      deletion_time: '',
      destroyed: false,
      version: this.version,
      backend: this.backend,
      path: this.path,
    });
    // nested secret
    this.secretDataComplex = {
      foo: {
        bar: 'baz',
      },
    };
    this.store.pushPayload('kv/data', {
      modelName: 'kv/data',
      id: this.dataIdComplex,
      secret_data: this.secretDataComplex,
      created_time: '2023-08-20T02:12:17.379762Z',
      custom_metadata: null,
      deletion_time: '',
      destroyed: false,
      version: this.version,
    });

    const metadata = this.server.create('kv-metadatum');
    metadata.id = this.metadataId;
    this.store.pushPayload('kv/metadata', {
      modelName: 'kv/metadata',
      ...metadata,
    });

    this.metadata = this.store.peekRecord('kv/metadata', this.metadataId);
    this.secret = this.store.peekRecord('kv/data', this.dataId);
    this.secretComplex = this.store.peekRecord('kv/data', this.dataIdComplex);

    // this is the route model, not an ember data model
    this.model = {
      backend: this.backend,
      path: this.path,
      secret: this.secret,
      metadata: this.metadata,
    };
    this.breadcrumbs = [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: this.model.backend, route: 'list' },
      { label: this.model.path },
    ];
    this.modelComplex = {
      backend: this.backend,
      path: this.pathComplex,
      secret: this.secretComplex,
      metadata: this.metadata,
    };
  });

  test('it renders secret details and toggles json view', async function (assert) {
    assert.expect(9);
    this.server.get(`sys/sync/associations/destinations`, (schema, req) => {
      assert.ok(true, 'request made to fetch sync status');
      assert.propEqual(
        req.queryParams,
        {
          mount: this.backend,
          secret_name: this.path,
        },
        'query params include mount and secret name'
      );
      // no records so response returns 404
      return syncStatusResponse(schema, req);
    });

    await render(
      hbs`
       <Page::Secret::Details
        @path={{this.model.path}}
        @secret={{this.model.secret}}
        @metadata={{this.model.metadata}}
        @breadcrumbs={{this.breadcrumbs}}
      />
      `,
      { owner: this.engine }
    );

    assert
      .dom(KV_SECRET.syncAlert())
      .doesNotExist('sync page alert banner does not render when sync status errors');
    assert.dom(GENERAL.title).includesText(this.model.path, 'renders secret path as page title');
    assert.dom(GENERAL.infoRowValue('foo')).exists('renders row for secret data');
    assert.dom(GENERAL.infoRowValue('foo')).hasText('***********');
    await click(KV_FORM.toggleMasked);
    assert.dom(GENERAL.infoRowValue('foo')).hasText('bar', 'renders secret value');
    await click(KV_FORM.toggleJson);
    await click(KV_FORM.toggleJsonValues);
    assert.propEqual(parseJsonEditor(find), this.secretData, 'json editor renders secret data');
    assert
      .dom(KV_SECRET.versionTimestamp)
      .includesText(`Version ${this.version} created`, 'renders version and time created');
  });

  test('it renders json view when secret is complex', async function (assert) {
    assert.expect(4);
    await render(
      hbs`
       <Page::Secret::Details
        @path={{this.modelComplex.path}}
        @secret={{this.modelComplex.secret}}
        @breadcrumbs={{this.breadcrumbs}}
      />
      `,
      { owner: this.engine }
    );
    assert.dom(GENERAL.infoRowValue('foo')).doesNotExist('does not render rows of secret data');
    assert.dom(KV_FORM.toggleJson).isChecked();
    assert.dom(KV_FORM.toggleJson).isNotDisabled();
    assert.dom('[data-test-component="code-mirror-modifier"]').exists('shows json editor');
  });

  test('it renders deleted empty state', async function (assert) {
    assert.expect(3);
    this.secret.deletionTime = '2023-07-23T02:12:17.379762Z';
    await render(
      hbs`
       <Page::Secret::Details
        @path={{this.model.path}}
        @secret={{this.model.secret}}
        @metadata={{this.model.metadata}}
        @breadcrumbs={{this.breadcrumbs}}
      />
      `,
      { owner: this.engine }
    );
    assert.dom(GENERAL.emptyStateTitle).hasText('Version 2 of this secret has been deleted');
    assert
      .dom(GENERAL.emptyStateMessage)
      .hasText(
        'This version has been deleted but can be undeleted. View other versions of this secret by clicking the Version History tab above.'
      );
    assert
      .dom(KV_SECRET.versionTimestamp)
      .includesText(`Version ${this.version} deleted`, 'renders version and time deleted');
  });

  test('it renders destroyed empty state', async function (assert) {
    assert.expect(2);
    this.secret.destroyed = true;
    await render(
      hbs`
       <Page::Secret::Details
        @path={{this.model.path}}
        @secret={{this.model.secret}}
        @metadata={{this.model.metadata}}
        @breadcrumbs={{this.breadcrumbs}}
      />
      `,
      { owner: this.engine }
    );
    assert.dom(GENERAL.emptyStateTitle).hasText('Version 2 of this secret has been permanently destroyed');
    assert
      .dom(GENERAL.emptyStateMessage)
      .hasText(
        'A version that has been permanently deleted cannot be restored. You can view other versions of this secret in the Version History tab above.'
      );
  });

  test('it renders secret version dropdown', async function (assert) {
    assert.expect(9);

    await render(
      hbs`
       <Page::Secret::Details
        @path={{this.model.path}}
        @secret={{this.model.secret}}
        @metadata={{this.model.metadata}}
        @breadcrumbs={{this.breadcrumbs}}
      />
      `,
      { owner: this.engine }
    );

    assert.dom(KV_SECRET.versionTimestamp).includesText(this.version, 'renders version');
    assert.dom(KV_SECRET.versionDropdown).hasText(`Version ${this.secret.version}`);
    await click(KV_SECRET.versionDropdown);

    for (const version in this.metadata.versions) {
      const data = this.metadata.versions[version];
      assert.dom(KV_SECRET.version(version)).exists(`renders ${version} in dropdown menu`);

      if (data.destroyed || data.deletion_time) {
        assert
          .dom(`${KV_SECRET.version(version)} [data-test-icon="x-square-fill"]`)
          .hasClass(`${data.destroyed ? 'has-text-danger' : 'has-text-grey'}`);
      }
    }

    assert
      .dom(`${KV_SECRET.version(this.metadata.currentVersion)} [data-test-icon="check-circle"]`)
      .exists('renders current version icon');
  });

  test('it renders sync status page alert and refreshes', async function (assert) {
    assert.expect(6); // assert count important because confirms request made to fetch sync status twice
    const destinationName = 'my-destination';
    this.server.create('sync-association', {
      type: 'aws-sm',
      name: destinationName,
      mount: this.backend,
      secret_name: this.path,
    });
    this.server.get(`sys/sync/associations/destinations`, (schema, req) => {
      // these assertions should be hit twice, once on init and again when the 'Refresh' button is clicked
      assert.ok(true, 'request made to fetch sync status');
      assert.propEqual(
        req.queryParams,
        {
          mount: this.backend,
          secret_name: this.path,
        },
        'query params include mount and secret name'
      );
      return syncStatusResponse(schema, req);
    });

    await render(
      hbs`
       <Page::Secret::Details
        @path={{this.model.path}}
        @secret={{this.model.secret}}
        @metadata={{this.model.metadata}}
        @breadcrumbs={{this.breadcrumbs}}
      />
      `,
      { owner: this.engine }
    );

    assert
      .dom(KV_SECRET.syncAlert(destinationName))
      .hasTextContaining(
        'Synced my-destination - last updated September',
        'renders sync status alert banner'
      );
    assert
      .dom(KV_SECRET.syncAlert())
      .hasTextContaining(
        'This secret has been synced from Vault to 1 destination. Updates to this secret will automatically sync to its destination.',
        'renders alert header referring to singular destination'
      );
    // sync status refresh button
    await click(`${KV_SECRET.syncAlert()} button`);
  });

  test('it renders sync status page alert for multiple destinations', async function (assert) {
    assert.expect(3); // assert count important because confirms request made to fetch sync status twice
    this.server.create('sync-association', {
      type: 'aws-sm',
      name: 'aws-dest',
      mount: this.backend,
      secret_name: this.path,
    });
    this.server.create('sync-association', {
      type: 'gh',
      name: 'gh-dest',
      mount: this.backend,
      secret_name: this.path,
    });
    this.server.get(`sys/sync/associations/destinations`, (schema, req) => {
      return syncStatusResponse(schema, req);
    });

    await render(
      hbs`
       <Page::Secret::Details
        @path={{this.model.path}}
        @secret={{this.model.secret}}
        @metadata={{this.model.metadata}}
        @breadcrumbs={{this.breadcrumbs}}
      />
      `,
      { owner: this.engine }
    );
    assert
      .dom(KV_SECRET.syncAlert('aws-dest'))
      .hasTextContaining('Synced aws-dest - last updated September', 'renders status for aws destination');
    assert
      .dom(KV_SECRET.syncAlert('gh-dest'))
      .hasTextContaining('Syncing gh-dest - last updated September', 'renders status for gh destination');
    assert
      .dom(KV_SECRET.syncAlert())
      .hasTextContaining(
        'This secret has been synced from Vault to 2 destinations. Updates to this secret will automatically sync to its destinations.',
        'renders alert title referring to plural destinations'
      );
  });
});
