import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { render } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';

module('Integration | Component | pki-issuer-cross-sign', function (hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function () {
    this.model = this.store.createRecord('pki/issuer', {
      issuerId: 'dcc69709-0056-b008-2ad1-a939cfae0c2a',
      issuerName: 'my-parent-issuer-name',
    });
  });

  // TODO finish
  test('it renders', async function (assert) {
    await render(hbs`<PkiIssuerCrossSign @parentIssuer={{this.model}} />`);
    assert.dom(this.element).exists();
  });
});
