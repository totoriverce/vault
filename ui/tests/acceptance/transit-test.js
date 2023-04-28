/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { click, fillIn, find, currentURL, settled, visit, waitUntil, findAll } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { v4 as uuidv4 } from 'uuid';

import { encodeString } from 'vault/utils/b64';
import authPage from 'vault/tests/pages/auth';
import enablePage from 'vault/tests/pages/settings/mount-secret-backend';
import secretListPage from 'vault/tests/pages/secrets/backend/list';

const keyTypes = [
  {
    name: (ts) => `aes-${ts}`,
    type: 'aes128-gcm96',
    exportable: true,
    supportsEncryption: true,
  },
  {
    name: (ts) => `aes-convergent-${ts}`,
    type: 'aes128-gcm96',
    convergent: true,
    supportsEncryption: true,
  },
  {
    name: (ts) => `aes-${ts}`,
    type: 'aes256-gcm96',
    exportable: true,
    supportsEncryption: true,
  },
  {
    name: (ts) => `aes-convergent-${ts}`,
    type: 'aes256-gcm96',
    convergent: true,
    supportsEncryption: true,
  },
  {
    name: (ts) => `chacha-${ts}`,
    type: 'chacha20-poly1305',
    exportable: true,
    supportsEncryption: true,
  },
  {
    name: (ts) => `chacha-convergent-${ts}`,
    type: 'chacha20-poly1305',
    convergent: true,
    supportsEncryption: true,
    autoRotate: true,
  },
  {
    name: (ts) => `ecdsa-${ts}`,
    type: 'ecdsa-p256',
    exportable: true,
    supportsSigning: true,
  },
  {
    name: (ts) => `ecdsa-${ts}`,
    type: 'ecdsa-p384',
    exportable: true,
    supportsSigning: true,
  },
  {
    name: (ts) => `ecdsa-${ts}`,
    type: 'ecdsa-p521',
    exportable: true,
    supportsSigning: true,
  },
  {
    name: (ts) => `ed25519-${ts}`,
    type: 'ed25519',
    derived: true,
    supportsSigning: true,
  },
  {
    name: (ts) => `rsa-2048-${ts}`,
    type: `rsa-2048`,
    supportsSigning: true,
    supportsEncryption: true,
  },
  {
    name: (ts) => `rsa-3072-${ts}`,
    type: `rsa-3072`,
    supportsSigning: true,
    supportsEncryption: true,
  },
  {
    name: (ts) => `rsa-4096-${ts}`,
    type: `rsa-4096`,
    supportsSigning: true,
    supportsEncryption: true,
    autoRotate: true,
  },
];

const generateTransitKey = async function (key, now) {
  const name = key.name(now);
  await click('[data-test-secret-create]');

  await fillIn('[data-test-transit-key-name]', name);
  await fillIn('[data-test-transit-key-type]', key.type);
  if (key.exportable) {
    await click('[data-test-transit-key-exportable]');
  }
  if (key.derived) {
    await click('[data-test-transit-key-derived]');
  }
  if (key.convergent) {
    await click('[data-test-transit-key-convergent-encryption]');
  }
  if (key.autoRotate) {
    await click('[data-test-toggle-label="Auto-rotation period"]');
  }
  await click('[data-test-transit-key-create]');
  await settled(); // eslint-disable-line
  // link back to the list
  await click('[data-test-secret-root-link]');

  return name;
};

const testConvergentEncryption = async function (assert, keyName) {
  const tests = [
    // raw bytes for plaintext and context
    {
      plaintext: 'NaXud2QW7KjyK6Me9ggh+zmnCeBGdG93LQED49PtoOI=',
      context: 'nqR8LiVgNh/lwO2rArJJE9F9DMhh0lKo4JX9DAAkCDw=',
      encodePlaintext: false,
      encodeContext: false,
      assertAfterEncrypt: (key) => {
        assert.dom('.modal.is-active').exists(`${key}: Modal opens after encrypt`);
        assert.ok(
          /vault:/.test(find('[data-test-encrypted-value="ciphertext"]').innerText),
          `${key}: ciphertext shows a vault-prefixed ciphertext`
        );
      },
      assertBeforeDecrypt: (key) => {
        assert.dom('.modal.is-active').doesNotExist(`${key}: Modal not open before decrypt`);
        assert
          .dom('[data-test-transit-input="context"]')
          .hasValue(
            'nqR8LiVgNh/lwO2rArJJE9F9DMhh0lKo4JX9DAAkCDw=',
            `${key}: the ui shows the base64-encoded context`
          );
      },

      assertAfterDecrypt: (key) => {
        assert.dom('.modal.is-active').exists(`${key}: Modal opens after decrypt`);
        assert.strictEqual(
          find('[data-test-encrypted-value="plaintext"]').innerText,
          'NaXud2QW7KjyK6Me9ggh+zmnCeBGdG93LQED49PtoOI=',
          `${key}: the ui shows the base64-encoded plaintext`
        );
      },
    },
    // raw bytes for plaintext, string for context
    {
      plaintext: 'NaXud2QW7KjyK6Me9ggh+zmnCeBGdG93LQED49PtoOI=',
      context: encodeString('context'),
      encodePlaintext: false,
      encodeContext: false,
      assertAfterEncrypt: (key) => {
        assert.dom('.modal.is-active').exists(`${key}: Modal opens after encrypt`);
        assert.ok(
          /vault:/.test(find('[data-test-encrypted-value="ciphertext"]').innerText),
          `${key}: ciphertext shows a vault-prefixed ciphertext`
        );
      },
      assertBeforeDecrypt: (key) => {
        assert.dom('.modal.is-active').doesNotExist(`${key}: Modal not open before decrypt`);
        assert
          .dom('[data-test-transit-input="context"]')
          .hasValue(encodeString('context'), `${key}: the ui shows the input context`);
      },
      assertAfterDecrypt: (key) => {
        assert.dom('.modal.is-active').exists(`${key}: Modal opens after decrypt`);
        assert.strictEqual(
          find('[data-test-encrypted-value="plaintext"]').innerText,
          'NaXud2QW7KjyK6Me9ggh+zmnCeBGdG93LQED49PtoOI=',
          `${key}: the ui shows the base64-encoded plaintext`
        );
      },
    },
    // base64 input
    {
      plaintext: encodeString('This is the secret'),
      context: encodeString('context'),
      encodePlaintext: false,
      encodeContext: false,
      assertAfterEncrypt: (key) => {
        assert.dom('.modal.is-active').exists(`${key}: Modal opens after encrypt`);
        assert.ok(
          /vault:/.test(find('[data-test-encrypted-value="ciphertext"]').innerText),
          `${key}: ciphertext shows a vault-prefixed ciphertext`
        );
      },
      assertBeforeDecrypt: (key) => {
        assert.dom('.modal.is-active').doesNotExist(`${key}: Modal not open before decrypt`);
        assert
          .dom('[data-test-transit-input="context"]')
          .hasValue(encodeString('context'), `${key}: the ui shows the input context`);
      },
      assertAfterDecrypt: (key) => {
        assert.dom('.modal.is-active').exists(`${key}: Modal opens after decrypt`);
        assert.strictEqual(
          find('[data-test-encrypted-value="plaintext"]').innerText,
          encodeString('This is the secret'),
          `${key}: the ui decodes plaintext`
        );
      },
    },

    // string input
    {
      plaintext: 'There are many secrets 🤐',
      context: 'secret 2',
      encodePlaintext: true,
      encodeContext: true,
      assertAfterEncrypt: (key) => {
        assert.dom('.modal.is-active').exists(`${key}: Modal opens after encrypt`);
        assert.ok(
          /vault:/.test(find('[data-test-encrypted-value="ciphertext"]').innerText),
          `${key}: ciphertext shows a vault-prefixed ciphertext`
        );
      },
      assertBeforeDecrypt: (key) => {
        assert.dom('.modal.is-active').doesNotExist(`${key}: Modal not open before decrypt`);
        assert
          .dom('[data-test-transit-input="context"]')
          .hasValue(encodeString('secret 2'), `${key}: the ui shows the encoded context`);
      },
      assertAfterDecrypt: (key) => {
        assert.dom('.modal.is-active').exists(`${key}: Modal opens after decrypt`);
        assert.strictEqual(
          find('[data-test-encrypted-value="plaintext"]').innerText,
          encodeString('There are many secrets 🤐'),
          `${key}: the ui decodes plaintext`
        );
      },
    },
  ];

  for (const testCase of tests) {
    await click('[data-test-transit-action-link="encrypt"]');

    find('#plaintext-control .CodeMirror').CodeMirror.setValue(testCase.plaintext);
    await fillIn('[data-test-transit-input="context"]', testCase.context);

    if (!testCase.encodePlaintext) {
      // If value is already encoded, check the box
      await click('input[data-test-transit-input="encodedBase64"]');
    }
    if (testCase.encodeContext) {
      await click('[data-test-transit-b64-toggle="context"]');
    }
    assert.dom('.modal.is-active').doesNotExist(`${name}: is not open before encrypt`);
    await click('[data-test-button-encrypt]');

    if (testCase.assertAfterEncrypt) {
      await settled();
      testCase.assertAfterEncrypt(keyName);
    }
    // store ciphertext for decryption step
    const copiedCiphertext = find('[data-test-encrypted-value="ciphertext"]').innerText;
    await click('.modal.is-active [data-test-modal-background]');

    assert.dom('.modal.is-active').doesNotExist(`${name}: Modal closes after background clicked`);
    await click('[data-test-transit-action-link="decrypt"]');

    if (testCase.assertBeforeDecrypt) {
      await settled();
      testCase.assertBeforeDecrypt(keyName);
    }
    find('#ciphertext-control .CodeMirror').CodeMirror.setValue(copiedCiphertext);
    await click('[data-test-button-decrypt]');

    if (testCase.assertAfterDecrypt) {
      await settled();
      testCase.assertAfterDecrypt(keyName);
    }

    await click('.modal.is-active [data-test-modal-background]');

    assert.dom('.modal.is-active').doesNotExist(`${name}: Modal closes after background clicked`);
  }
};
module('Acceptance | transit', function (hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function () {
    const uid = uuidv4();
    await authPage.login();
    await settled();
    this.uid = uid;
    this.path = `transit-${uid}`;

    await enablePage.enable('transit', `transit-${uid}`);
    await settled();
  });

  test(`transit backend: list menu`, async function (assert) {
    await generateTransitKey(keyTypes[0], this.uid);
    await secretListPage.secrets.objectAt(0).menuToggle();
    await settled();
    assert.strictEqual(secretListPage.menuItems.length, 2, 'shows 2 items in the menu');
  });
  for (const key of keyTypes) {
    test(`transit backend: ${key.type}`, async function (assert) {
      assert.expect(key.convergent ? 43 : 7);
      const name = await generateTransitKey(key, this.uid);
      await visit(`vault/secrets/${this.path}/show/${name}`);

      const expectedRotateValue = key.autoRotate ? '30 days' : 'Key will not be automatically rotated';
      assert
        .dom('[data-test-row-value="Auto-rotation period"]')
        .hasText(expectedRotateValue, 'Has expected auto rotate value');

      await click('[data-test-transit-link="versions"]');
      // wait for capabilities

      assert.dom('[data-test-transit-key-version-row]').exists({ count: 1 }, `${name}: only one key version`);
      await waitUntil(() => find('[data-test-confirm-action-trigger]'));
      await click('[data-test-confirm-action-trigger]');

      await click('[data-test-confirm-button]');
      // wait for rotate call
      await waitUntil(() => findAll('[data-test-transit-key-version-row]').length >= 2);
      assert
        .dom('[data-test-transit-key-version-row]')
        .exists({ count: 2 }, `${name}: two key versions after rotate`);
      await click('[data-test-transit-key-actions-link]');

      assert.ok(
        currentURL().startsWith(`/vault/secrets/${this.path}/show/${name}?tab=actions`),
        `${name}: navigates to transit actions`
      );

      const keyAction = key.supportsEncryption ? 'encrypt' : 'sign';
      const actionTitle = find(`[data-test-transit-action-title=${keyAction}]`).innerText.toLowerCase();

      assert.true(
        actionTitle.includes(keyAction),
        `shows a card with title that links to the ${name} transit action`
      );

      await click(`[data-test-transit-card=${keyAction}]`);

      assert
        .dom('[data-test-transit-key-version-select]')
        .exists(`${name}: the rotated key allows you to select versions`);
      if (key.exportable) {
        assert
          .dom('[data-test-transit-action-link="export"]')
          .exists(`${name}: exportable key has a link to export action`);
      } else {
        assert
          .dom('[data-test-transit-action-link="export"]')
          .doesNotExist(`${name}: non-exportable key does not link to export action`);
      }
      if (key.convergent && key.supportsEncryption) {
        await testConvergentEncryption(assert, name);
        await settled();
      }
      await settled();
    });
  }
});
