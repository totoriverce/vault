/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import { module, test } from 'qunit';
import { v4 as uuidv4 } from 'uuid';
import { click, currentURL, fillIn, typeIn, visit } from '@ember/test-helpers';

import { setupApplicationTest } from 'vault/tests/helpers';
import authPage from 'vault/tests/pages/auth';
import { deleteEngineCmd, mountEngineCmd, runCmd, tokenWithPolicyCmd } from 'vault/tests/helpers/commands';
import { personas } from 'vault/tests/helpers/kv/policy-generator';
import { clearRecords, writeVersionedSecret } from 'vault/tests/helpers/kv/kv-run-commands';
import { FORM, KV_WORKFLOW } from 'vault/tests/helpers/kv/kv-selectors';
import { grantAccessForWrite, setupControlGroup } from 'vault/tests/helpers/control-groups';
import { GENERAL } from 'vault/tests/helpers/general-selectors';
import { KV_METADATA_DETAILS } from 'vault/tests/helpers/components/kv/page/secret/metadata/details-selectors';

/**
 * This test set is for testing the flow for creating new secrets and versions.
 * Letter(s) in parenthesis at the end are shorthand for the persona,
 * for ease of tracking down specific tests failures from CI
 */
module('Acceptance | kv-v2 workflow | secret and version create', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    this.backend = `kv-create-${uuidv4()}`;
    this.store = this.owner.lookup('service:store');
    await authPage.login();
    await runCmd(mountEngineCmd('kv-v2', this.backend), false);
    await writeVersionedSecret(this.backend, 'app/first', 'foo', 'bar', 2);
  });

  hooks.afterEach(async function () {
    await authPage.login();
    return runCmd(deleteEngineCmd(this.backend));
  });

  module('admin persona', function (hooks) {
    hooks.beforeEach(async function () {
      const token = await runCmd(tokenWithPolicyCmd(`admin-${this.backend}`, personas.admin(this.backend)));
      await authPage.login(token);
      clearRecords(this.store);
      return;
    });
    test('cancel on create clears model (a)', async function (assert) {
      const backend = this.backend;
      await visit(`/vault/secrets/${backend}/kv/list`);
      assert.dom(KV_WORKFLOW.list.item()).exists({ count: 1 }, 'single secret exists on list');
      assert.dom(KV_WORKFLOW.list.item('app/')).hasText('app/', 'expected list item');
      await click(KV_WORKFLOW.list.createSecret);
      await fillIn(FORM.inputByAttr('path'), 'jk');
      await click(FORM.cancelBtn);
      assert.dom(KV_WORKFLOW.list.item()).exists({ count: 1 }, 'same amount of secrets');
      assert.dom(KV_WORKFLOW.list.item('app/')).hasText('app/', 'expected list item');
      await click(KV_WORKFLOW.list.createSecret);
      await fillIn(FORM.inputByAttr('path'), 'psych');
      await click(GENERAL.breadcrumbAtIdx(1));
      assert.dom(KV_WORKFLOW.list.item()).exists({ count: 1 }, 'same amount of secrets');
      assert.dom(KV_WORKFLOW.list.item('app/')).hasText('app/', 'expected list item');
    });
    test('cancel on new version rolls back model (a)', async function (assert) {
      const backend = this.backend;
      await visit(`/vault/secrets/${backend}/kv/${encodeURIComponent('app/first')}/details`);
      assert.dom(GENERAL.infoRowValue('foo')).exists('key has expected value');
      await click(KV_WORKFLOW.detail.createNewVersion);
      await fillIn(FORM.keyInput(), 'bar');
      await click(FORM.cancelBtn);
      assert.dom(GENERAL.infoRowValue('foo')).exists('secret is previous value');
      await click(KV_WORKFLOW.detail.createNewVersion);
      await fillIn(FORM.keyInput(), 'bar');
      await click(GENERAL.breadcrumbAtIdx(3));
      assert.dom(GENERAL.infoRowValue('foo')).exists('secret is previous value');
    });
    test('create & update root secret with default metadata (a)', async function (assert) {
      const backend = this.backend;
      const secretPath = 'some secret';
      await visit(`/vault/secrets/${backend}/kv/list`);
      await click(KV_WORKFLOW.list.createSecret);

      // Create secret form -- validations
      await click(FORM.saveBtn);
      assert.dom(FORM.invalidFormAlert).hasText('There is an error with this form.');
      assert.dom(FORM.validation('path')).hasText("Path can't be blank.");
      await typeIn(FORM.inputByAttr('path'), secretPath);
      assert
        .dom(FORM.validationWarning)
        .hasText(
          "Path contains whitespace. If this is desired, you'll need to encode it with %20 in API requests."
        );
      assert.dom(KV_WORKFLOW.create.metadataSection).doesNotExist('Hides metadata section by default');

      // Submit with API errors
      await click(FORM.saveBtn);
      assert.dom(FORM.messageError).hasText('Error no data provided', 'API error shows on form');

      await fillIn(FORM.keyInput(), 'api_key');
      await fillIn(FORM.maskedValueInput(), 'partyparty');
      await click(FORM.saveBtn);

      // Details page
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/${encodeURIComponent(secretPath)}/details?version=1`,
        'Goes to details page after save'
      );
      assert.dom(KV_WORKFLOW.detail.versionTimestamp).includesText('Version 1 created');
      assert.dom(KV_WORKFLOW.infoRow).exists({ count: 1 }, '1 row of data shows');
      assert.dom(GENERAL.infoRowValue('api_key')).hasText('***********');
      await click(KV_WORKFLOW.infoRowToggleMasked('api_key'));
      assert.dom(GENERAL.infoRowValue('api_key')).hasText('partyparty', 'secret value shows after toggle');

      // Metadata page
      await click(GENERAL.tab('Metadata'));
      assert
        .dom(`${KV_METADATA_DETAILS.customMetadataSection} ${GENERAL.emptyStateTitle}`)
        .hasText('No custom metadata', 'No custom metadata empty state');
      assert
        .dom(`${KV_METADATA_DETAILS.secretMetadataSection} ${KV_WORKFLOW.infoRow}`)
        .exists({ count: 4 }, '4 metadata rows show');
      assert.dom(GENERAL.infoRowValue('Maximum versions')).hasText('0', 'max versions shows 0');
      assert.dom(GENERAL.infoRowValue('Check-and-Set required')).hasText('No', 'cas not enforced');
      assert
        .dom(GENERAL.infoRowValue('Delete version after'))
        .hasText('Never delete', 'Delete version after has default 0s');

      // Add new version
      await click(GENERAL.tab('Secret'));
      await click(KV_WORKFLOW.detail.createNewVersion);
      assert.dom(FORM.inputByAttr('path')).isDisabled('path input is disabled');
      assert.dom(FORM.inputByAttr('path')).hasValue(secretPath);
      assert.dom(FORM.toggleMetadata).doesNotExist('Does not show metadata toggle when creating new version');
      assert.dom(FORM.keyInput()).hasValue('api_key');
      assert.dom(FORM.maskedValueInput()).hasValue('partyparty');
      await fillIn(FORM.keyInput(1), 'api_url');
      await fillIn(FORM.maskedValueInput(1), 'hashicorp.com');
      await click(FORM.saveBtn);

      // Back to details page
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/${encodeURIComponent(secretPath)}/details?version=2`
      );
      assert.dom(KV_WORKFLOW.detail.versionTimestamp).includesText('Version 2 created');
      assert.dom(KV_WORKFLOW.infoRow).exists({ count: 2 }, '2 rows of data shows');
      assert.dom(GENERAL.infoRowValue('api_key')).hasText('***********');
      assert.dom(GENERAL.infoRowValue('api_url')).hasText('***********');
      await click(KV_WORKFLOW.infoRowToggleMasked('api_key'));
      await click(KV_WORKFLOW.infoRowToggleMasked('api_url'));
      assert.dom(GENERAL.infoRowValue('api_key')).hasText('partyparty', 'secret value shows after toggle');
      assert.dom(GENERAL.infoRowValue('api_url')).hasText('hashicorp.com', 'secret value shows after toggle');
    });
    test('create nested secret with metadata (a)', async function (assert) {
      const backend = this.backend;
      await visit(`/vault/secrets/${backend}/kv/list`);
      await click(KV_WORKFLOW.list.createSecret);

      // Create secret
      await typeIn(FORM.inputByAttr('path'), 'my/');
      assert.dom(FORM.validation('path')).hasText("Path can't end in forward slash '/'.");
      await typeIn(FORM.inputByAttr('path'), 'secret');
      assert.dom(FORM.validation('path')).doesNotExist('form validation goes away');
      await fillIn(FORM.keyInput(), 'password');
      await fillIn(FORM.maskedValueInput(), 'kittens1234');

      await click(FORM.toggleMetadata);
      assert.dom(KV_WORKFLOW.create.metadataSection).exists('Shows metadata section after toggled');
      // Check initial values
      assert.dom(FORM.inputByAttr('maxVersions')).hasValue('0');
      assert.dom(FORM.inputByAttr('casRequired')).isNotChecked();
      assert.dom(FORM.toggleByLabel('Automate secret deletion')).isNotChecked();
      // MaxVersions validation
      await fillIn(FORM.inputByAttr('maxVersions'), 'seven');
      await click(FORM.saveBtn);
      assert.dom(FORM.validation('maxVersions')).hasText('Maximum versions must be a number.');
      await fillIn(FORM.inputByAttr('maxVersions'), '99999999999999999');
      await click(FORM.saveBtn);
      assert.dom(FORM.validation('maxVersions')).hasText('You cannot go over 16 characters.');
      await fillIn(FORM.inputByAttr('maxVersions'), '7');

      // Fill in other metadata
      await click(FORM.inputByAttr('casRequired'));
      await click(FORM.toggleByLabel('Automate secret deletion'));
      await fillIn(FORM.ttlValue('Automate secret deletion'), '1000');

      // Fill in custom metadata
      await fillIn(`${KV_WORKFLOW.create.metadataSection} ${FORM.keyInput()}`, 'team');
      await fillIn(`${KV_WORKFLOW.create.metadataSection} ${FORM.valueInput()}`, 'UI');
      // Fill in metadata
      await click(FORM.saveBtn);

      // Details
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/${encodeURIComponent('my/secret')}/details?version=1`
      );
      assert.dom(KV_WORKFLOW.detail.versionTimestamp).includesText('Version 1 created');
      assert.dom(KV_WORKFLOW.infoRow).exists({ count: 1 }, '1 row of data shows');
      assert.dom(GENERAL.infoRowValue('password')).hasText('***********');
      await click(KV_WORKFLOW.infoRowToggleMasked('password'));
      assert.dom(GENERAL.infoRowValue('password')).hasText('kittens1234', 'secret value shows after toggle');

      // Metadata
      await click(GENERAL.tab('Metadata'));
      assert
        .dom(`${KV_METADATA_DETAILS.customMetadataSection} ${KV_WORKFLOW.infoRow}`)
        .exists({ count: 1 }, 'One custom metadata row shows');
      assert
        .dom(`${KV_METADATA_DETAILS.customMetadataSection} ${GENERAL.infoRowValue('team')}`)
        .hasText('UI');

      assert
        .dom(`${KV_METADATA_DETAILS.secretMetadataSection} ${KV_WORKFLOW.infoRow}`)
        .exists({ count: 4 }, '4 metadata rows show');
      assert.dom(GENERAL.infoRowValue('Maximum versions')).hasText('7', 'max versions shows 0');
      assert.dom(GENERAL.infoRowValue('Check-and-Set required')).hasText('Yes', 'cas enforced');
      assert
        .dom(GENERAL.infoRowValue('Delete version after'))
        .hasText('16 minutes 40 seconds', 'Delete version after has custom value');
    });
    test('creates a secret at a sub-directory (a)', async function (assert) {
      const backend = this.backend;
      await visit(`/vault/secrets/${backend}/kv/list/app/`);
      assert.dom(KV_WORKFLOW.list.item('first')).exists('Lists first sub-secret');
      assert.dom(KV_WORKFLOW.list.item('new')).doesNotExist('Does not show new secret');
      await click(KV_WORKFLOW.list.createSecret);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/create?initialKey=app%2F`,
        'Goes to create page with initialKey'
      );
      await typeIn(FORM.inputByAttr('path'), 'new');
      await fillIn(FORM.keyInput(), 'api_key');
      await fillIn(FORM.maskedValueInput(), 'partyparty');
      await click(FORM.saveBtn);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/${encodeURIComponent('app/new')}/details?version=1`,
        'Redirects to detail after save'
      );
      await click(GENERAL.breadcrumbAtIdx(2));
      assert.strictEqual(currentURL(), `/vault/secrets/${backend}/kv/list/app/`, 'sub-dir page');
      assert.dom(KV_WORKFLOW.list.item('new')).exists('Lists new secret in sub-dir');
    });
    test('create new version of secret from older version (a)', async function (assert) {
      const backend = this.backend;
      await visit(`/vault/secrets/${backend}/kv/app%2Ffirst/details`);
      await click(KV_WORKFLOW.detail.versionDropdown);
      await click(`${KV_WORKFLOW.detail.version(1)} a`);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/app%2Ffirst/details?version=1`,
        'goes to version 1'
      );
      assert.dom(KV_WORKFLOW.detail.versionTimestamp).includesText('Version 1 created');
      await click(KV_WORKFLOW.detail.createNewVersion);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/app%2Ffirst/details/edit?version=1`,
        'Goes to new version page'
      );
      assert
        .dom(FORM.versionAlert)
        .hasText(
          'Warning You are creating a new version based on data from Version 1. The current version for app/first is Version 2.',
          'Shows version warning'
        );
      assert.dom(FORM.keyInput()).hasValue('key-1', 'Key input has old value');
      assert.dom(FORM.maskedValueInput()).hasValue('val-1', 'Val input has old value');

      await fillIn(FORM.keyInput(), 'my-key');
      await fillIn(FORM.maskedValueInput(), 'my-value');
      await click(FORM.saveBtn);

      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/app%2Ffirst/details?version=3`,
        'goes to latest version 3'
      );
      await click(KV_WORKFLOW.infoRowToggleMasked('my-key'));
      assert.dom(GENERAL.infoRowValue('my-key')).hasText('my-value', 'has new value');
    });
  });

  module('data-reader persona', function (hooks) {
    hooks.beforeEach(async function () {
      const token = await runCmd(
        tokenWithPolicyCmd(`data-reader-${this.backend}`, personas.dataReader(this.backend))
      );
      await authPage.login(token);
      clearRecords(this.store);
      return;
    });
    test('cancel on create clears model (dr)', async function (assert) {
      const backend = this.backend;
      await visit(`/vault/secrets/${backend}/kv/list`);
      assert.dom(KV_WORKFLOW.list.item()).doesNotExist('list view has no items');
      await click(KV_WORKFLOW.list.createSecret);
      await fillIn(FORM.inputByAttr('path'), 'jk');
      await click(FORM.cancelBtn);
      assert.dom(KV_WORKFLOW.list.item()).doesNotExist('list view still has no items');
      await click(KV_WORKFLOW.list.createSecret);
      await fillIn(FORM.inputByAttr('path'), 'psych');
      await click(GENERAL.breadcrumbAtIdx(1));
      assert.dom(KV_WORKFLOW.list.item()).doesNotExist('list view still has no items');
    });
    test('cancel on new version rolls back model (dr)', async function (assert) {
      const backend = this.backend;
      await visit(`/vault/secrets/${backend}/kv/${encodeURIComponent('app/first')}/details`);
      assert.dom(GENERAL.infoRowValue('foo')).exists('key has expected value');
      assert.dom(KV_WORKFLOW.detail.createNewVersion).doesNotExist();
    });
    test('create & update root secret with default metadata (dr)', async function (assert) {
      const backend = this.backend;
      const secretPath = 'some secret';
      await visit(`/vault/secrets/${backend}/kv/list`);
      await click(KV_WORKFLOW.list.createSecret);

      // Create secret form -- validations
      await click(FORM.saveBtn);
      assert.dom(FORM.invalidFormAlert).hasText('There is an error with this form.');
      assert.dom(FORM.validation('path')).hasText("Path can't be blank.");
      await typeIn(FORM.inputByAttr('path'), secretPath);
      assert
        .dom(FORM.validationWarning)
        .hasText(
          "Path contains whitespace. If this is desired, you'll need to encode it with %20 in API requests."
        );
      assert.dom(KV_WORKFLOW.create.metadataSection).doesNotExist('Hides metadata section by default');

      // Submit with API errors
      await click(FORM.saveBtn);
      assert
        .dom(FORM.messageError)
        .hasText('Error 1 error occurred: * permission denied', 'API error shows on form');

      // Since this persona can't create a new secret, test update with existing:
      await visit(`/vault/secrets/${backend}/kv/app%2Ffirst/details`);
      assert.dom(KV_WORKFLOW.detail.versionTimestamp).includesText('Version 2 created');
      assert.dom(KV_WORKFLOW.infoRow).exists({ count: 1 }, '1 row of data shows');
      assert.dom(GENERAL.infoRowValue('foo')).hasText('***********');
      await click(KV_WORKFLOW.infoRowToggleMasked('foo'));
      assert.dom(GENERAL.infoRowValue('foo')).hasText('bar', 'secret value shows after toggle');

      // Metadata page
      await click(GENERAL.tab('Metadata'));
      assert
        .dom(`${KV_METADATA_DETAILS.customMetadataSection} ${GENERAL.emptyStateTitle}`)
        .hasText('No custom metadata', 'No custom metadata empty state');
      assert
        .dom(`${KV_METADATA_DETAILS.secretMetadataSection} ${GENERAL.emptyStateTitle}`)
        .hasText('You do not have access to secret metadata', 'shows no access state on metadata');

      // Add new version
      await click(GENERAL.tab('Secret'));
      assert.dom(KV_WORKFLOW.detail.createNewVersion).doesNotExist('cannot create new version');
    });
    test('create nested secret with metadata (dr)', async function (assert) {
      const backend = this.backend;
      await visit(`/vault/secrets/${backend}/kv/list`);
      await click(KV_WORKFLOW.list.createSecret);

      // Create secret
      await typeIn(FORM.inputByAttr('path'), 'my/');
      assert.dom(FORM.validation('path')).hasText("Path can't end in forward slash '/'.");
      await typeIn(FORM.inputByAttr('path'), 'secret');
      assert.dom(FORM.validation('path')).doesNotExist('form validation goes away');
      await fillIn(FORM.keyInput(), 'password');
      await fillIn(FORM.maskedValueInput(), 'kittens1234');

      await click(FORM.toggleMetadata);
      assert.dom(KV_WORKFLOW.create.metadataSection).exists('Shows metadata section after toggled');
      // Check initial values
      assert.dom(FORM.inputByAttr('maxVersions')).hasValue('0');
      assert.dom(FORM.inputByAttr('casRequired')).isNotChecked();
      assert.dom(FORM.toggleByLabel('Automate secret deletion')).isNotChecked();
      // MaxVersions validation
      await fillIn(FORM.inputByAttr('maxVersions'), 'seven');
      await click(FORM.saveBtn);
      assert.dom(FORM.validation('maxVersions')).hasText('Maximum versions must be a number.');
      await fillIn(FORM.inputByAttr('maxVersions'), '99999999999999999');
      await click(FORM.saveBtn);
      assert.dom(FORM.validation('maxVersions')).hasText('You cannot go over 16 characters.');
      await fillIn(FORM.inputByAttr('maxVersions'), '7');

      // Fill in other metadata
      await click(FORM.inputByAttr('casRequired'));
      await click(FORM.toggleByLabel('Automate secret deletion'));
      await fillIn(FORM.ttlValue('Automate secret deletion'), '1000');

      // Fill in custom metadata
      await fillIn(`${KV_WORKFLOW.create.metadataSection} ${FORM.keyInput()}`, 'team');
      await fillIn(`${KV_WORKFLOW.create.metadataSection} ${FORM.valueInput()}`, 'UI');
      // Fill in metadata
      await click(FORM.saveBtn);
      assert
        .dom(FORM.messageError)
        .hasText('Error 1 error occurred: * permission denied', 'API error shows on form');
    });
    test('creates a secret at a sub-directory (dr)', async function (assert) {
      const backend = this.backend;
      await visit(`/vault/secrets/${backend}/kv/list/app/`);
      assert.dom(KV_WORKFLOW.list.item()).doesNotExist('Does not list any secrets');
      await click(KV_WORKFLOW.list.createSecret);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/create?initialKey=app%2F`,
        'Goes to create page with initialKey'
      );
      await typeIn(FORM.inputByAttr('path'), 'new');
      await fillIn(FORM.keyInput(), 'api_key');
      await fillIn(FORM.maskedValueInput(), 'partyparty');
      await click(FORM.saveBtn);
      assert
        .dom(FORM.messageError)
        .hasText('Error 1 error occurred: * permission denied', 'API error shows on form');
    });
    test('create new version of secret from older version (dr)', async function (assert) {
      const backend = this.backend;
      await visit(`/vault/secrets/${backend}/kv/app%2Ffirst/details?version=1`);
      assert.dom(KV_WORKFLOW.detail.versionDropdown).doesNotExist('version dropdown does not show');
      assert.dom(KV_WORKFLOW.detail.versionTimestamp).includesText('Version 1 created');
      assert.dom(KV_WORKFLOW.detail.createNewVersion).doesNotExist('cannot create new version');
    });
  });

  module('data-list-reader persona', function (hooks) {
    hooks.beforeEach(async function () {
      const token = await runCmd(
        tokenWithPolicyCmd(`data-list-reader-${this.backend}`, personas.dataListReader(this.backend))
      );
      await authPage.login(token);
      clearRecords(this.store);
      return;
    });
    test('cancel on create clears model (dlr)', async function (assert) {
      const backend = this.backend;
      await visit(`/vault/secrets/${backend}/kv/list`);
      assert.dom(KV_WORKFLOW.list.item()).exists({ count: 1 }, 'single secret exists on list');
      assert.dom(KV_WORKFLOW.list.item('app/')).hasText('app/', 'expected list item');
      await click(KV_WORKFLOW.list.createSecret);
      await fillIn(FORM.inputByAttr('path'), 'jk');
      await click(FORM.cancelBtn);
      assert.dom(KV_WORKFLOW.list.item()).exists({ count: 1 }, 'same amount of secrets');
      assert.dom(KV_WORKFLOW.list.item('app/')).hasText('app/', 'expected list item');
      await click(KV_WORKFLOW.list.createSecret);
      await fillIn(FORM.inputByAttr('path'), 'psych');
      await click(GENERAL.breadcrumbAtIdx(1));
      assert.dom(KV_WORKFLOW.list.item()).exists({ count: 1 }, 'same amount of secrets');
      assert.dom(KV_WORKFLOW.list.item('app/')).hasText('app/', 'expected list item');
    });
    test('cancel on new version rolls back model (dlr)', async function (assert) {
      const backend = this.backend;
      await visit(`/vault/secrets/${backend}/kv/${encodeURIComponent('app/first')}/details`);
      assert.dom(GENERAL.infoRowValue('foo')).exists('key has expected value');
      assert.dom(KV_WORKFLOW.detail.createNewVersion).doesNotExist('cannot create new version');
    });
    test('create & update root secret with default metadata (dlr)', async function (assert) {
      const backend = this.backend;
      const secretPath = 'some secret';
      await visit(`/vault/secrets/${backend}/kv/list`);
      await click(KV_WORKFLOW.list.createSecret);

      // Create secret form -- validations
      await click(FORM.saveBtn);
      assert.dom(FORM.invalidFormAlert).hasText('There is an error with this form.');
      assert.dom(FORM.validation('path')).hasText("Path can't be blank.");
      await typeIn(FORM.inputByAttr('path'), secretPath);
      assert
        .dom(FORM.validationWarning)
        .hasText(
          "Path contains whitespace. If this is desired, you'll need to encode it with %20 in API requests."
        );
      assert.dom(KV_WORKFLOW.create.metadataSection).doesNotExist('Hides metadata section by default');

      // Submit with API errors
      await click(FORM.saveBtn);
      assert
        .dom(FORM.messageError)
        .hasText('Error 1 error occurred: * permission denied', 'API error shows on form');

      // Since this persona can't create a new secret, test update with existing:
      await visit(`/vault/secrets/${backend}/kv/app%2Ffirst/details`);
      assert.dom(KV_WORKFLOW.detail.versionTimestamp).includesText('Version 2 created');
      assert.dom(KV_WORKFLOW.infoRow).exists({ count: 1 }, '1 row of data shows');
      assert.dom(GENERAL.infoRowValue('foo')).hasText('***********');
      await click(KV_WORKFLOW.infoRowToggleMasked('foo'));
      assert.dom(GENERAL.infoRowValue('foo')).hasText('bar', 'secret value shows after toggle');

      // Metadata page
      await click(GENERAL.tab('Metadata'));
      assert
        .dom(`${KV_METADATA_DETAILS.customMetadataSection} ${GENERAL.emptyStateTitle}`)
        .hasText('No custom metadata', 'No custom metadata empty state');
      assert
        .dom(`${KV_METADATA_DETAILS.secretMetadataSection} ${GENERAL.emptyStateTitle}`)
        .hasText('You do not have access to secret metadata', 'shows no access state on metadata');

      // Add new version
      await click(GENERAL.tab('Secret'));
      assert.dom(KV_WORKFLOW.detail.createNewVersion).doesNotExist('cannot create new version');
    });
    test('create nested secret with metadata (dlr)', async function (assert) {
      const backend = this.backend;
      await visit(`/vault/secrets/${backend}/kv/list`);
      await click(KV_WORKFLOW.list.createSecret);

      // Create secret
      await typeIn(FORM.inputByAttr('path'), 'my/');
      assert.dom(FORM.validation('path')).hasText("Path can't end in forward slash '/'.");
      await typeIn(FORM.inputByAttr('path'), 'secret');
      assert.dom(FORM.validation('path')).doesNotExist('form validation goes away');
      await fillIn(FORM.keyInput(), 'password');
      await fillIn(FORM.maskedValueInput(), 'kittens1234');

      await click(FORM.toggleMetadata);
      assert.dom(KV_WORKFLOW.create.metadataSection).exists('Shows metadata section after toggled');
      // Check initial values
      assert.dom(FORM.inputByAttr('maxVersions')).hasValue('0');
      assert.dom(FORM.inputByAttr('casRequired')).isNotChecked();
      assert.dom(FORM.toggleByLabel('Automate secret deletion')).isNotChecked();
      // MaxVersions validation
      await fillIn(FORM.inputByAttr('maxVersions'), 'seven');
      await click(FORM.saveBtn);
      assert.dom(FORM.validation('maxVersions')).hasText('Maximum versions must be a number.');
      await fillIn(FORM.inputByAttr('maxVersions'), '99999999999999999');
      await click(FORM.saveBtn);
      assert.dom(FORM.validation('maxVersions')).hasText('You cannot go over 16 characters.');
      await fillIn(FORM.inputByAttr('maxVersions'), '7');

      // Fill in other metadata
      await click(FORM.inputByAttr('casRequired'));
      await click(FORM.toggleByLabel('Automate secret deletion'));
      await fillIn(FORM.ttlValue('Automate secret deletion'), '1000');

      // Fill in custom metadata
      await fillIn(`${KV_WORKFLOW.create.metadataSection} ${FORM.keyInput()}`, 'team');
      await fillIn(`${KV_WORKFLOW.create.metadataSection} ${FORM.valueInput()}`, 'UI');
      // Fill in metadata
      await click(FORM.saveBtn);
      assert
        .dom(FORM.messageError)
        .hasText('Error 1 error occurred: * permission denied', 'API error shows on form');
    });
    test('creates a secret at a sub-directory (dlr)', async function (assert) {
      const backend = this.backend;
      await visit(`/vault/secrets/${backend}/kv/list/app/`);
      assert.dom(KV_WORKFLOW.list.item()).doesNotExist('Does not list any secrets');
      await click(KV_WORKFLOW.list.createSecret);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/create?initialKey=app%2F`,
        'Goes to create page with initialKey'
      );
      await typeIn(FORM.inputByAttr('path'), 'new');
      await fillIn(FORM.keyInput(), 'api_key');
      await fillIn(FORM.maskedValueInput(), 'partyparty');
      await click(FORM.saveBtn);
      assert
        .dom(FORM.messageError)
        .hasText('Error 1 error occurred: * permission denied', 'API error shows on form');
    });
    test('create new version of secret from older version (dlr)', async function (assert) {
      const backend = this.backend;
      await visit(`/vault/secrets/${backend}/kv/app%2Ffirst/details?version=1`);
      assert.dom(KV_WORKFLOW.detail.versionDropdown).doesNotExist('version dropdown does not show');
      assert.dom(KV_WORKFLOW.detail.versionTimestamp).includesText('Version 1 created');
      assert.dom(KV_WORKFLOW.detail.createNewVersion).doesNotExist('cannot create new version');
    });
  });

  module('metadata-maintainer persona', function (hooks) {
    hooks.beforeEach(async function () {
      const token = await runCmd(
        tokenWithPolicyCmd(`data-list-reader-${this.backend}`, personas.metadataMaintainer(this.backend))
      );
      await authPage.login(token);
      clearRecords(this.store);
      return;
    });
    test('cancel on create clears model (mm)', async function (assert) {
      const backend = this.backend;
      await visit(`/vault/secrets/${backend}/kv/list`);
      assert.dom(KV_WORKFLOW.list.item()).exists({ count: 1 }, 'single secret exists on list');
      assert.dom(KV_WORKFLOW.list.item('app/')).hasText('app/', 'expected list item');
      await click(KV_WORKFLOW.list.createSecret);
      await fillIn(FORM.inputByAttr('path'), 'jk');
      await click(FORM.cancelBtn);
      assert.dom(KV_WORKFLOW.list.item()).exists({ count: 1 }, 'same amount of secrets');
      assert.dom(KV_WORKFLOW.list.item('app/')).hasText('app/', 'expected list item');
      await click(KV_WORKFLOW.list.createSecret);
      await fillIn(FORM.inputByAttr('path'), 'psych');
      await click(GENERAL.breadcrumbAtIdx(1));
      assert.dom(KV_WORKFLOW.list.item()).exists({ count: 1 }, 'same amount of secrets');
      assert.dom(KV_WORKFLOW.list.item('app/')).hasText('app/', 'expected list item');
    });
    test('cancel on new version rolls back model (mm)', async function (assert) {
      const backend = this.backend;
      await visit(`/vault/secrets/${backend}/kv/${encodeURIComponent('app/first')}/details`);
      assert.dom(GENERAL.emptyStateTitle).hasText('You do not have permission to read this secret');
      assert
        .dom(KV_WORKFLOW.detail.createNewVersion)
        .doesNotExist('create new version button now allowed since user cannot read existing');
    });
    test('create & update root secret with default metadata (mm)', async function (assert) {
      const backend = this.backend;
      const secretPath = 'some secret';
      await visit(`/vault/secrets/${backend}/kv/list`);
      await click(KV_WORKFLOW.list.createSecret);

      // Create secret form -- validations
      await click(FORM.saveBtn);
      assert.dom(FORM.invalidFormAlert).hasText('There is an error with this form.');
      assert.dom(FORM.validation('path')).hasText("Path can't be blank.");
      await typeIn(FORM.inputByAttr('path'), secretPath);
      assert
        .dom(FORM.validationWarning)
        .hasText(
          "Path contains whitespace. If this is desired, you'll need to encode it with %20 in API requests."
        );
      assert.dom(KV_WORKFLOW.create.metadataSection).doesNotExist('Hides metadata section by default');

      // Submit with API errors
      await click(FORM.saveBtn);
      assert
        .dom(FORM.messageError)
        .hasText('Error 1 error occurred: * permission denied', 'API error shows on form');

      // Since this persona can't create a new secret, test update with existing:
      await visit(`/vault/secrets/${backend}/kv/app%2Ffirst/details`);
      assert.dom(KV_WORKFLOW.detail.versionTimestamp).doesNotExist('Version created tooltip does not show');
      assert.dom(KV_WORKFLOW.infoRow).doesNotExist('secret data not shown');
      assert.dom(GENERAL.emptyStateTitle).hasText('You do not have permission to read this secret');

      // Metadata page
      await click(GENERAL.tab('Metadata'));
      assert
        .dom(`${KV_METADATA_DETAILS.customMetadataSection} ${GENERAL.emptyStateTitle}`)
        .hasText('No custom metadata', 'No custom metadata empty state');
      assert
        .dom(`${KV_METADATA_DETAILS.secretMetadataSection} ${KV_WORKFLOW.infoRow}`)
        .exists({ count: 4 }, '4 metadata rows show');
      assert.dom(GENERAL.infoRowValue('Maximum versions')).hasText('0', 'max versions shows 0');
      assert.dom(GENERAL.infoRowValue('Check-and-Set required')).hasText('No', 'cas not enforced');
      assert
        .dom(GENERAL.infoRowValue('Delete version after'))
        .hasText('Never delete', 'Delete version after has default 0s');

      // Add new version
      await click(GENERAL.tab('Secret'));
      assert.dom(KV_WORKFLOW.detail.createNewVersion).doesNotExist('create new version button not rendered');
      await visit(`/vault/secrets/${backend}/kv/app%2Ffirst/details/edit?version=1`);
      assert
        .dom(FORM.noReadAlert)
        .hasText(
          'Warning You do not have read permissions for this secret data. Saving will overwrite the existing secret.',
          'shows alert for no read permissions'
        );

      assert.dom(FORM.inputByAttr('path')).isDisabled('path input is disabled');
      assert.dom(FORM.inputByAttr('path')).hasValue('app/first');
      assert.dom(FORM.toggleMetadata).doesNotExist('Does not show metadata toggle when creating new version');
      assert.dom(FORM.keyInput()).hasValue('', 'first row has no key');
      assert.dom(FORM.maskedValueInput()).hasValue('', 'first row has no value');
      await fillIn(FORM.keyInput(), 'api_url');
      await fillIn(FORM.maskedValueInput(), 'hashicorp.com');
      await click(FORM.saveBtn);
      assert
        .dom(FORM.messageError)
        .hasText('Error 1 error occurred: * permission denied', 'API error shows on form');
    });
    test('create nested secret with metadata (mm)', async function (assert) {
      const backend = this.backend;
      await visit(`/vault/secrets/${backend}/kv/list`);
      await click(KV_WORKFLOW.list.createSecret);

      // Create secret
      await typeIn(FORM.inputByAttr('path'), 'my/');
      assert.dom(FORM.validation('path')).hasText("Path can't end in forward slash '/'.");
      await typeIn(FORM.inputByAttr('path'), 'secret');
      assert.dom(FORM.validation('path')).doesNotExist('form validation goes away');
      await fillIn(FORM.keyInput(), 'password');
      await fillIn(FORM.maskedValueInput(), 'kittens1234');

      await click(FORM.toggleMetadata);
      assert.dom(KV_WORKFLOW.create.metadataSection).exists('Shows metadata section after toggled');
      // Check initial values
      assert.dom(FORM.inputByAttr('maxVersions')).hasValue('0');
      assert.dom(FORM.inputByAttr('casRequired')).isNotChecked();
      assert.dom(FORM.toggleByLabel('Automate secret deletion')).isNotChecked();
      // MaxVersions validation
      await fillIn(FORM.inputByAttr('maxVersions'), 'seven');
      await click(FORM.saveBtn);
      assert.dom(FORM.validation('maxVersions')).hasText('Maximum versions must be a number.');
      await fillIn(FORM.inputByAttr('maxVersions'), '99999999999999999');
      await click(FORM.saveBtn);
      assert.dom(FORM.validation('maxVersions')).hasText('You cannot go over 16 characters.');
      await fillIn(FORM.inputByAttr('maxVersions'), '7');

      // Fill in other metadata
      await click(FORM.inputByAttr('casRequired'));
      await click(FORM.toggleByLabel('Automate secret deletion'));
      await fillIn(FORM.ttlValue('Automate secret deletion'), '1000');

      // Fill in custom metadata
      await fillIn(`${KV_WORKFLOW.create.metadataSection} ${FORM.keyInput()}`, 'team');
      await fillIn(`${KV_WORKFLOW.create.metadataSection} ${FORM.valueInput()}`, 'UI');

      await click(FORM.saveBtn);
      assert
        .dom(FORM.messageError)
        .hasText('Error 1 error occurred: * permission denied', 'API error shows on form');
    });
    test('creates a secret at a sub-directory (mm)', async function (assert) {
      const backend = this.backend;
      await visit(`/vault/secrets/${backend}/kv/list/app/`);
      assert.dom(KV_WORKFLOW.list.item('first')).exists('Lists first sub-secret');
      assert.dom(KV_WORKFLOW.list.item('new')).doesNotExist('Does not show new secret');
      await click(KV_WORKFLOW.list.createSecret);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/create?initialKey=app%2F`,
        'Goes to create page with initialKey'
      );
      await typeIn(FORM.inputByAttr('path'), 'new');
      await fillIn(FORM.keyInput(), 'api_key');
      await fillIn(FORM.maskedValueInput(), 'partyparty');
      await click(FORM.saveBtn);
      assert
        .dom(FORM.messageError)
        .hasText('Error 1 error occurred: * permission denied', 'API error shows on form');
    });
    test('create new version of secret from older version (mm)', async function (assert) {
      const backend = this.backend;
      await visit(`/vault/secrets/${backend}/kv/app%2Ffirst/details`);
      assert.dom(KV_WORKFLOW.detail.versionDropdown).hasText('Version 2');
      await click(KV_WORKFLOW.detail.versionDropdown);
      await click(`${KV_WORKFLOW.detail.version(1)} a`);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/app%2Ffirst/details?version=1`,
        'goes to version 1'
      );
      assert.dom(KV_WORKFLOW.detail.versionDropdown).hasText('Version 1');
      assert.dom(KV_WORKFLOW.detail.versionTimestamp).doesNotExist('version timestamp not shown');
      assert.dom(KV_WORKFLOW.detail.createNewVersion).doesNotExist('create new version button not rendered');
      await visit(`/vault/secrets/${backend}/kv/app%2Ffirst/details/edit?version=1`);
      assert
        .dom(FORM.noReadAlert)
        .hasText(
          'Warning You do not have read permissions for this secret data. Saving will overwrite the existing secret.',
          'shows alert for no read permissions'
        );

      assert.dom(FORM.inputByAttr('path')).isDisabled('path input is disabled');
      assert.dom(FORM.inputByAttr('path')).hasValue('app/first');
      assert.dom(FORM.toggleMetadata).doesNotExist('Does not show metadata toggle when creating new version');
      assert.dom(FORM.keyInput()).hasValue('', 'first row has no key');
      assert.dom(FORM.maskedValueInput()).hasValue('', 'first row has no value');
      await fillIn(FORM.keyInput(), 'api_url');
      await fillIn(FORM.maskedValueInput(), 'hashicorp.com');
      await click(FORM.saveBtn);
      assert
        .dom(FORM.messageError)
        .hasText('Error 1 error occurred: * permission denied', 'API error shows on form');
    });
  });

  module('secret-creator persona', function (hooks) {
    hooks.beforeEach(async function () {
      const token = await runCmd(
        tokenWithPolicyCmd(`secret-creator-${this.backend}`, personas.secretCreator(this.backend))
      );
      await authPage.login(token);
      clearRecords(this.store);
      return;
    });
    test('cancel on create clears model (sc)', async function (assert) {
      const backend = this.backend;
      await visit(`/vault/secrets/${backend}/kv/list`);
      assert.dom(KV_WORKFLOW.list.item()).doesNotExist('list view has no items');
      await click(KV_WORKFLOW.list.createSecret);
      await fillIn(FORM.inputByAttr('path'), 'jk');
      await click(FORM.cancelBtn);
      assert.dom(KV_WORKFLOW.list.item()).doesNotExist('list view still has no items');
      await click(KV_WORKFLOW.list.createSecret);
      await fillIn(FORM.inputByAttr('path'), 'psych');
      await click(GENERAL.breadcrumbAtIdx(1));
      assert.dom(KV_WORKFLOW.list.item()).doesNotExist('list view still has no items');
    });
    test('cancel on new version rolls back model (sc)', async function (assert) {
      const backend = this.backend;
      await visit(`/vault/secrets/${backend}/kv/${encodeURIComponent('app/first')}/details`);
      assert
        .dom(GENERAL.emptyStateTitle)
        .hasText('You do not have permission to read this secret', 'no permissions state shows');
      await click(KV_WORKFLOW.detail.createNewVersion);
      await fillIn(FORM.keyInput(), 'bar');
      await click(FORM.cancelBtn);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/${encodeURIComponent('app/first')}/details`,
        'cancel goes to correct url'
      );
      assert.dom(KV_WORKFLOW.list.item()).doesNotExist('list view has no items');
      await click(KV_WORKFLOW.detail.createNewVersion);
      await fillIn(FORM.keyInput(), 'bar');
      await click(GENERAL.breadcrumbAtIdx(3));
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/${encodeURIComponent('app/first')}/details`,
        'breadcrumb goes to correct url'
      );
      assert.dom(KV_WORKFLOW.list.item()).doesNotExist('list view has no items');
    });
    test('create & update root secret with default metadata (sc)', async function (assert) {
      const backend = this.backend;
      const secretPath = 'some secret';
      await visit(`/vault/secrets/${backend}/kv/list`);
      await click(KV_WORKFLOW.list.createSecret);

      // Create secret form -- validations
      await click(FORM.saveBtn);
      assert.dom(FORM.invalidFormAlert).hasText('There is an error with this form.');
      assert.dom(FORM.validation('path')).hasText("Path can't be blank.");
      await typeIn(FORM.inputByAttr('path'), secretPath);
      assert
        .dom(FORM.validationWarning)
        .hasText(
          "Path contains whitespace. If this is desired, you'll need to encode it with %20 in API requests."
        );
      assert.dom(KV_WORKFLOW.create.metadataSection).doesNotExist('Hides metadata section by default');

      // Submit with API errors
      await click(FORM.saveBtn);
      assert.dom(FORM.messageError).hasText('Error no data provided', 'API error shows on form');

      await fillIn(FORM.keyInput(), 'api_key');
      await fillIn(FORM.maskedValueInput(), 'partyparty');
      await click(FORM.saveBtn);

      // Details page
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/${encodeURIComponent(secretPath)}/details`,
        'Goes to details page after save'
      );
      assert.dom(KV_WORKFLOW.detail.versionTimestamp).doesNotExist('Version created not shown');
      assert.dom(KV_WORKFLOW.infoRow).doesNotExist('does not show data contents');
      assert
        .dom(GENERAL.emptyStateTitle)
        .hasText('You do not have permission to read this secret', 'shows permissions empty state');

      // Metadata page
      await click(GENERAL.tab('Metadata'));
      assert
        .dom(`${KV_METADATA_DETAILS.customMetadataSection} ${GENERAL.emptyStateTitle}`)
        .hasText(
          'You do not have access to read custom metadata',
          'permissions empty state for custom metadata'
        );
      assert
        .dom(`${KV_METADATA_DETAILS.secretMetadataSection} ${GENERAL.emptyStateTitle}`)
        .hasText('You do not have access to secret metadata', 'permissions empty state for secret metadata');

      // Add new version
      await click(GENERAL.tab('Secret'));
      await click(KV_WORKFLOW.detail.createNewVersion);
      assert.dom(FORM.inputByAttr('path')).isDisabled('path input is disabled');
      assert.dom(FORM.inputByAttr('path')).hasValue(secretPath);
      assert.dom(FORM.toggleMetadata).doesNotExist('Does not show metadata toggle when creating new version');
      assert.dom(FORM.keyInput()).hasValue('', 'row 1 is empty key');
      assert.dom(FORM.maskedValueInput()).hasValue('', 'row 1 has empty value');
      await fillIn(FORM.keyInput(), 'api_url');
      await fillIn(FORM.maskedValueInput(), 'hashicorp.com');
      await click(FORM.saveBtn);

      // Back to details page
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/${encodeURIComponent(secretPath)}/details?version=2`,
        'goes back to details page'
      );
      assert.dom(KV_WORKFLOW.detail.versionTimestamp).doesNotExist('Version created does not show');
      assert.dom(KV_WORKFLOW.infoRow).doesNotExist('does not show data contents');
      assert
        .dom(GENERAL.emptyStateTitle)
        .hasText('You do not have permission to read this secret', 'shows permissions empty state');
    });
    test('create nested secret with metadata (sc)', async function (assert) {
      const backend = this.backend;
      await visit(`/vault/secrets/${backend}/kv/list`);
      await click(KV_WORKFLOW.list.createSecret);

      // Create secret
      await typeIn(FORM.inputByAttr('path'), 'my/');
      assert.dom(FORM.validation('path')).hasText("Path can't end in forward slash '/'.");
      await typeIn(FORM.inputByAttr('path'), 'secret');
      assert.dom(FORM.validation('path')).doesNotExist('form validation goes away');
      await fillIn(FORM.keyInput(), 'password');
      await fillIn(FORM.maskedValueInput(), 'kittens1234');

      await click(FORM.toggleMetadata);
      assert.dom(KV_WORKFLOW.create.metadataSection).exists('Shows metadata section after toggled');
      // Check initial values
      assert.dom(FORM.inputByAttr('maxVersions')).hasValue('0');
      assert.dom(FORM.inputByAttr('casRequired')).isNotChecked();
      assert.dom(FORM.toggleByLabel('Automate secret deletion')).isNotChecked();
      // MaxVersions validation
      await fillIn(FORM.inputByAttr('maxVersions'), 'seven');
      await click(FORM.saveBtn);
      assert.dom(FORM.validation('maxVersions')).hasText('Maximum versions must be a number.');
      await fillIn(FORM.inputByAttr('maxVersions'), '99999999999999999');
      await click(FORM.saveBtn);
      assert.dom(FORM.validation('maxVersions')).hasText('You cannot go over 16 characters.');
      await fillIn(FORM.inputByAttr('maxVersions'), '7');

      // Fill in other metadata
      await click(FORM.inputByAttr('casRequired'));
      await click(FORM.toggleByLabel('Automate secret deletion'));
      await fillIn(FORM.ttlValue('Automate secret deletion'), '1000');

      // Fill in custom metadata
      await fillIn(`${KV_WORKFLOW.create.metadataSection} ${FORM.keyInput()}`, 'team');
      await fillIn(`${KV_WORKFLOW.create.metadataSection} ${FORM.valueInput()}`, 'UI');
      // Fill in metadata
      await click(FORM.saveBtn);

      // Details
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/${encodeURIComponent('my/secret')}/details`,
        'goes back to details page'
      );
      assert.dom(KV_WORKFLOW.detail.versionTimestamp).doesNotExist('version created not shown');
      assert.dom(KV_WORKFLOW.infoRow).doesNotExist('does not show data contents');
      assert
        .dom(GENERAL.emptyStateTitle)
        .hasText('You do not have permission to read this secret', 'shows permissions empty state');

      // Metadata
      await click(GENERAL.tab('Metadata'));
      assert
        .dom(`${KV_METADATA_DETAILS.customMetadataSection} ${GENERAL.emptyStateTitle}`)
        .hasText(
          'You do not have access to read custom metadata',
          'permissions empty state for custom metadata'
        );
      assert
        .dom(`${KV_METADATA_DETAILS.secretMetadataSection} ${GENERAL.emptyStateTitle}`)
        .hasText('You do not have access to secret metadata', 'permissions empty state for secret metadata');
    });
    test('creates a secret at a sub-directory (sc)', async function (assert) {
      const backend = this.backend;
      await visit(`/vault/secrets/${backend}/kv/list/app/`);
      assert.dom(KV_WORKFLOW.list.item()).doesNotExist('Does not list any secrets');
      await click(KV_WORKFLOW.list.createSecret);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/create?initialKey=app%2F`,
        'Goes to create page with initialKey'
      );
      await typeIn(FORM.inputByAttr('path'), 'new');
      await fillIn(FORM.keyInput(), 'api_key');
      await fillIn(FORM.maskedValueInput(), 'partyparty');
      await click(FORM.saveBtn);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/${encodeURIComponent('app/new')}/details`,
        'Redirects to detail after save'
      );
      await click(GENERAL.breadcrumbAtIdx(2));
      assert.strictEqual(currentURL(), `/vault/secrets/${backend}/kv/list/app/`, 'sub-dir page');
      assert.dom(KV_WORKFLOW.list.item()).doesNotExist('Does not list any secrets');
    });
    test('create new version of secret from older version (sc)', async function (assert) {
      const backend = this.backend;
      await visit(`/vault/secrets/${backend}/kv/app%2Ffirst/details?version=1`);
      assert.dom(KV_WORKFLOW.detail.versionDropdown).doesNotExist('version dropdown does not show');
      assert.dom(KV_WORKFLOW.detail.versionTimestamp).doesNotExist('Version created not shown');
      await click(KV_WORKFLOW.detail.createNewVersion);
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/app%2Ffirst/details/edit?version=1`,
        'Goes to new version page'
      );
      assert
        .dom(FORM.noReadAlert)
        .hasText(
          'Warning You do not have read permissions for this secret data. Saving will overwrite the existing secret.',
          'shows alert for no read permissions'
        );
      assert.dom(FORM.keyInput()).hasValue('', 'Key input has empty value');
      assert.dom(FORM.maskedValueInput()).hasValue('', 'Val input has empty value');

      await fillIn(FORM.keyInput(), 'my-key');
      await fillIn(FORM.maskedValueInput(), 'my-value');
      await click(FORM.saveBtn);

      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/app%2Ffirst/details?version=3`,
        'redirects to details page'
      );
      assert.dom(KV_WORKFLOW.infoRow).doesNotExist('does not show data contents');
      assert
        .dom(GENERAL.emptyStateTitle)
        .hasText('You do not have permission to read this secret', 'shows permissions empty state');
    });
  });

  module('secret-nested-creator persona', function (hooks) {
    hooks.beforeEach(async function () {
      const token = await runCmd(
        tokenWithPolicyCmd(
          `secret-nested-creator-${this.backend}`,
          personas.secretNestedCreator(this.backend)
        )
      );
      await authPage.login(token);
      clearRecords(this.store);
      return;
    });
    test('can create a secret from the nested list view (snc)', async function (assert) {
      assert.expect(1);
      // go to nested secret directory list view
      await visit(`/vault/secrets/${this.backend}/kv/list/app/`);
      // correct popup menu items appear on list view
      const popupSelector = `${KV_WORKFLOW.list.item('first')} ${GENERAL.menuTrigger}`;
      await click(popupSelector);
      assert.dom(KV_WORKFLOW.list.listMenuCreate).exists('shows the option to create new version');
    });
  });

  module('enterprise controlled access persona', function (hooks) {
    hooks.beforeEach(async function () {
      this.controlGroup = this.owner.lookup('service:control-group');
      const userPolicy = `
path "${this.backend}/data/*" {
  capabilities = ["create", "read", "update"]
  control_group = {
    max_ttl = "24h"
    factor "authorizer" {
      controlled_capabilities = ["create", "update"]
      identity {
          group_names = ["managers"]
          approvals = 1
      }
    }
  }
}

path "${this.backend}/metadata" {
  capabilities = ["list", "read"]
}
path "${this.backend}/metadata/*" {
  capabilities = ["list", "read"]
}
`;

      const { userToken } = await setupControlGroup({
        userPolicy,
        backend: this.backend,
      });
      this.userToken = userToken;
      await authPage.login(userToken);
      clearRecords(this.store);
      return;
    });
    test('create & update root secret with default metadata (cg)', async function (assert) {
      const backend = this.backend;
      // Known issue: control groups do not work correctly in UI when encodable characters in path
      const secretPath = 'some-secret';
      await visit(`/vault/secrets/${backend}/kv/list`);
      await click(KV_WORKFLOW.list.createSecret);

      // Create secret form -- validations
      await click(FORM.saveBtn);
      assert.dom(FORM.invalidFormAlert).hasText('There is an error with this form.');
      assert.dom(FORM.validation('path')).hasText("Path can't be blank.");
      await typeIn(FORM.inputByAttr('path'), secretPath);
      assert.dom(KV_WORKFLOW.create.metadataSection).doesNotExist('Hides metadata section by default');

      await fillIn(FORM.keyInput(), 'api_key');
      await fillIn(FORM.maskedValueInput(), 'partyparty');
      await click(FORM.saveBtn);
      let tokenToUnwrap = this.controlGroup.tokenToUnwrap;
      assert.deepEqual(
        Object.keys(tokenToUnwrap),
        ['accessor', 'token', 'creation_path', 'creation_time', 'ttl'],
        'stored tokenToUnwrap includes correct keys'
      );
      assert.strictEqual(
        tokenToUnwrap.creation_path,
        `${backend}/data/${secretPath}`,
        'stored tokenToUnwrap includes correct creation path'
      );
      assert
        .dom(FORM.messageError)
        .includesText(
          `Error A Control Group was encountered at ${backend}/data/${secretPath}.`,
          'shows control group error'
        );
      await grantAccessForWrite({
        token: tokenToUnwrap.token,
        accessor: tokenToUnwrap.accessor,
        creation_path: `${backend}/data/${secretPath}`,
        originUrl: `/vault/secrets/${backend}/kv/create`,
        userToken: this.userToken,
        backend: this.backend,
      });
      // In a real scenario the user would stay on page, but in the test
      // we fill in the same info and try again
      await typeIn(FORM.inputByAttr('path'), secretPath);
      await fillIn(FORM.keyInput(), 'this can be anything');
      await fillIn(FORM.maskedValueInput(), 'this too, gonna use the wrapped data');
      await click(FORM.saveBtn);
      assert.strictEqual(this.controlGroup.tokenToUnwrap, null, 'clears tokenToUnwrap after successful save');
      // Details page
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/${secretPath}/details?version=1`,
        'Goes to details page after save'
      );
      assert.dom(KV_WORKFLOW.detail.versionTimestamp).includesText('Version 1 created');
      assert.dom(KV_WORKFLOW.infoRow).exists({ count: 1 }, '1 row of data shows');
      assert.dom(GENERAL.infoRowValue('api_key')).hasText('***********');
      await click(KV_WORKFLOW.infoRowToggleMasked('api_key'));
      assert.dom(GENERAL.infoRowValue('api_key')).hasText('partyparty', 'secret value shows after toggle');

      // Metadata page
      await click(GENERAL.tab('Metadata'));
      assert
        .dom(`${KV_METADATA_DETAILS.customMetadataSection} ${GENERAL.emptyStateTitle}`)
        .hasText('No custom metadata', 'No custom metadata empty state');
      assert
        .dom(`${KV_METADATA_DETAILS.secretMetadataSection} ${KV_WORKFLOW.infoRow}`)
        .exists({ count: 4 }, '4 metadata rows show');
      assert.dom(GENERAL.infoRowValue('Maximum versions')).hasText('0', 'max versions shows 0');
      assert.dom(GENERAL.infoRowValue('Check-and-Set required')).hasText('No', 'cas not enforced');
      assert
        .dom(GENERAL.infoRowValue('Delete version after'))
        .hasText('Never delete', 'Delete version after has default 0s');

      // Add new version
      await click(GENERAL.tab('Secret'));
      await click(KV_WORKFLOW.detail.createNewVersion);
      assert.dom(FORM.inputByAttr('path')).isDisabled('path input is disabled');
      assert.dom(FORM.inputByAttr('path')).hasValue(secretPath);
      assert.dom(FORM.toggleMetadata).doesNotExist('Does not show metadata toggle when creating new version');
      assert.dom(FORM.keyInput()).hasValue('api_key');
      assert.dom(FORM.maskedValueInput()).hasValue('partyparty');
      await fillIn(FORM.keyInput(1), 'api_url');
      await fillIn(FORM.maskedValueInput(1), 'hashicorp.com');
      await click(FORM.saveBtn);
      tokenToUnwrap = this.controlGroup.tokenToUnwrap;
      assert.strictEqual(
        tokenToUnwrap.creation_path,
        `${backend}/data/${secretPath}`,
        'stored tokenToUnwrap includes correct update path'
      );
      assert
        .dom(FORM.messageError)
        .includesText(
          `Error A Control Group was encountered at ${backend}/data/${secretPath}.`,
          'shows control group error'
        );
      // Normally the user stays on the page and tries again once approval is granted
      // unmark on test so it doesn't use the control group on read at the same path
      // when we return to the page after granting access below
      this.controlGroup.unmarkTokenForUnwrap();
      await grantAccessForWrite({
        accessor: tokenToUnwrap.accessor,
        token: tokenToUnwrap.token,
        creation_path: `${backend}/data/${secretPath}`,
        originUrl: `/vault/secrets/${backend}/kv/${secretPath}/details/edit`,
        userToken: this.userToken,
        backend: this.backend,
      });
      // Remark for unwrap as if we never left the page.
      this.controlGroup.markTokenForUnwrap(tokenToUnwrap.accessor);
      // No need to fill in data because we're using the stored wrapped request
      // and the path already exists
      await click(FORM.saveBtn);
      assert.strictEqual(
        this.controlGroup.tokenToUnwrap,
        null,
        'clears tokenToUnwrap after successful update'
      );

      // Back to details page
      assert.strictEqual(
        currentURL(),
        `/vault/secrets/${backend}/kv/${encodeURIComponent(secretPath)}/details?version=2`
      );
      assert.dom(KV_WORKFLOW.detail.versionTimestamp).includesText('Version 2 created');
      assert.dom(KV_WORKFLOW.infoRow).exists({ count: 2 }, '2 rows of data shows');
      assert.dom(GENERAL.infoRowValue('api_key')).hasText('***********');
      assert.dom(GENERAL.infoRowValue('api_url')).hasText('***********');
      await click(KV_WORKFLOW.infoRowToggleMasked('api_key'));
      await click(KV_WORKFLOW.infoRowToggleMasked('api_url'));
      assert.dom(GENERAL.infoRowValue('api_key')).hasText('partyparty', 'secret value shows after toggle');
      assert.dom(GENERAL.infoRowValue('api_url')).hasText('hashicorp.com', 'secret value shows after toggle');
    });
  });
});
