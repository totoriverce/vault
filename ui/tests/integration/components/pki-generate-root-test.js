import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { click, fillIn, render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import Sinon from 'sinon';
import { allowAllCapabilitiesStub } from 'vault/tests/helpers/stubs';

const SELECTORS = {
  mainSectionTitle: '[data-test-generate-root-title="Root parameters"]',
  urlSectionTitle: '[data-test-generate-root-title="Issuer URLs"]',
  keyParamsGroupToggle: '[data-test-toggle-group="Key parameters"]',
  sanGroupToggle: '[data-test-toggle-group="Subject Alternative Name (SAN) Options"]',
  additionalGroupToggle: '[data-test-toggle-group="Additional subject fields"]',
  toggleGroupDescription: '[data-test-toggle-group-description]',
  formField: '[data-test-field]',
  typeField: '[data-test-input="type"]',
  fieldByName: (name) => `[data-test-field="${name}"]`,
  saveButton: '[data-test-pki-generate-root-save]',
  cancelButton: '[data-test-pki-generate-root-cancel]',
  formInvalidError: '[data-test-pki-generate-root-validation-error]',
};

module('Integration | Component | pki-generate-root', function (hooks) {
  setupRenderingTest(hooks);
  setupMirage(hooks);
  setupEngine(hooks, 'pki');

  hooks.beforeEach(async function () {
    this.server.post('/sys/capabilities-self', allowAllCapabilitiesStub());
    this.store = this.owner.lookup('service:store');
    this.secretMountPath = this.owner.lookup('service:secret-mount-path');
    this.secretMountPath.currentPath = 'pki-test';
    this.model = this.store.createRecord('pki/config', {
      role: 'my-role',
    });
    this.onSave = Sinon.spy();
    this.onCancel = Sinon.spy();
  });

  test('it renders with correct sections', async function (assert) {
    await render(
      hbs`<PkiGenerateRoot @model={{this.model}} @onSave={{this.onSave}} @onCancel={{this.onCancel}} />`,
      {
        owner: this.engine,
      }
    );

    // Titles
    assert.dom('h2').exists({ count: 2 }, 'two H2 titles are visible on page load');
    assert.dom(SELECTORS.mainSectionTitle).hasText('Root parameters');
    assert.dom(SELECTORS.urlSectionTitle).hasText('Issuer URLs');

    assert.dom('[data-test-toggle-group]').exists({ count: 3 }, '3 toggle groups shown');
  });

  test('it shows the appropriate fields under the toggles', async function (assert) {
    await render(
      hbs`<PkiGenerateRoot @model={{this.model}} @onSave={{this.onSave}} @onCancel={{this.onCancel}} />`,
      {
        owner: this.engine,
      }
    );

    await click(SELECTORS.additionalGroupToggle);
    assert
      .dom(SELECTORS.toggleGroupDescription)
      .hasText('These fields provide more information about the client to which the certificate belongs.');
    assert
      .dom(`[data-test-group="Additional subject fields"] ${SELECTORS.formField}`)
      .exists({ count: 7 }, '7 form fields under Additional Fields toggle');

    await click(SELECTORS.sanGroupToggle);
    assert
      .dom(SELECTORS.toggleGroupDescription)
      .hasText(
        'SAN fields are an extension that allow you specify additional host names (sites, IP addresses, common names, etc.) to be protected by a single certificate.'
      );
    assert
      .dom(`[data-test-group="Subject Alternative Name (SAN) Options"] ${SELECTORS.formField}`)
      .exists({ count: 6 }, '7 form fields under SANs toggle');

    await click(SELECTORS.keyParamsGroupToggle);
    assert
      .dom(SELECTORS.toggleGroupDescription)
      .hasText(
        'Please choose a type to see key parameter options.',
        'Shows empty state description before type is selected'
      );
    assert
      .dom(`[data-test-group="Key parameters"] ${SELECTORS.formField}`)
      .exists({ count: 0 }, '0 form fields under keyParams toggle');
  });

  test('it renders the correct form fields in key params', async function (assert) {
    this.set('type', '');
    await render(
      hbs`<PkiGenerateRoot @model={{this.model}} @onSave={{this.onSave}} @onCancel={{this.onCancel}} />`,
      {
        owner: this.engine,
      }
    );
    await click(SELECTORS.keyParamsGroupToggle);
    assert
      .dom(`[data-test-group="Key parameters"] ${SELECTORS.formField}`)
      .exists({ count: 0 }, '0 form fields under keyParams toggle');

    this.set('type', 'exported');
    await fillIn(SELECTORS.typeField, this.type);
    assert
      .dom(SELECTORS.toggleGroupDescription)
      .hasText(
        'This certificate type is exported. This means the private key will be returned in the response. Below, you will name the key and define its type and key bits.',
        `has correct description for type=${this.type}`
      );
    assert.strictEqual(this.model.type, this.type);
    assert
      .dom(`[data-test-group="Key parameters"] ${SELECTORS.formField}`)
      .exists({ count: 3 }, '3 form fields under keyParams toggle');
    assert.dom(SELECTORS.fieldByName('keyName')).exists(`Key name field shown when type=${this.type}`);
    assert.dom(SELECTORS.fieldByName('keyType')).exists(`Key type field shown when type=${this.type}`);
    assert.dom(SELECTORS.fieldByName('keyBits')).exists(`Key bits field shown when type=${this.type}`);

    this.set('type', 'internal');
    await fillIn(SELECTORS.typeField, this.type);
    assert
      .dom(SELECTORS.toggleGroupDescription)
      .hasText(
        'This certificate type is internal. This means that the private key will not be returned and cannot be retrieved later. Below, you will name the key and define its type and key bits.',
        `has correct description for type=${this.type}`
      );
    assert.strictEqual(this.model.type, this.type);
    assert
      .dom(`[data-test-group="Key parameters"] ${SELECTORS.formField}`)
      .exists({ count: 3 }, '3 form fields under keyParams toggle');
    assert.dom(SELECTORS.fieldByName('keyName')).exists(`Key name field shown when type=${this.type}`);
    assert.dom(SELECTORS.fieldByName('keyType')).exists(`Key type field shown when type=${this.type}`);
    assert.dom(SELECTORS.fieldByName('keyBits')).exists(`Key bits field shown when type=${this.type}`);

    this.set('type', 'existing');
    await fillIn(SELECTORS.typeField, this.type);
    assert
      .dom(SELECTORS.toggleGroupDescription)
      .hasText(
        'You chose to use an existing key. This means that we’ll use the key reference to create the CSR or root. Please provide the reference to the key.',
        `has correct description for type=${this.type}`
      );
    assert.strictEqual(this.model.type, this.type);
    assert
      .dom(`[data-test-group="Key parameters"] ${SELECTORS.formField}`)
      .exists({ count: 1 }, '1 form field under keyParams toggle');
    assert.dom(SELECTORS.fieldByName('keyRef')).exists(`Key reference field shown when type=${this.type}`);

    this.set('type', 'kms');
    await fillIn(SELECTORS.typeField, this.type);
    assert
      .dom(SELECTORS.toggleGroupDescription)
      .hasText(
        'This certificate type is kms, meaning managed keys will be used. Below, you will name the key and tell Vault where to find it in your KMS or HSM. Learn more about managed keys.',
        `has correct description for type=${this.type}`
      );
    assert.strictEqual(this.model.type, this.type);
    assert
      .dom(`[data-test-group="Key parameters"] ${SELECTORS.formField}`)
      .exists({ count: 3 }, '3 form fields under keyParams toggle');
    assert.dom(SELECTORS.fieldByName('keyName')).exists(`Key name field shown when type=${this.type}`);
    assert
      .dom(SELECTORS.fieldByName('managedKeyName'))
      .exists(`Managed key name field shown when type=${this.type}`);
    assert
      .dom(SELECTORS.fieldByName('managedKeyId'))
      .exists(`Managed key id field shown when type=${this.type}`);
  });

  test('it shows errors before submit if form is invalid', async function (assert) {
    const saveSpy = Sinon.spy();
    this.set('onSave', saveSpy);
    await render(
      hbs`<PkiGenerateRoot @model={{this.model}} @onSave={{this.onSave}} @onCancel={{this.onCancel}} />`,
      {
        owner: this.engine,
      }
    );

    await click(SELECTORS.saveButton);
    assert.dom(SELECTORS.formInvalidError).hasText('The form has errors. Please fix before continuing.');
    assert.ok(saveSpy.notCalled);
  });
});
