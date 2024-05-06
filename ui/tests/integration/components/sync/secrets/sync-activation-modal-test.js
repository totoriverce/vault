import { module, test } from 'qunit';
import { setupRenderingTest } from 'vault/tests/helpers';
import { setupEngine } from 'ember-engines/test-support';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { render, click } from '@ember/test-helpers';
import { hbs } from 'ember-cli-htmlbars';
import sinon from 'sinon';
import { overrideResponse } from 'vault/tests/helpers/stubs';
import { PAGE } from 'vault/tests/helpers/sync/sync-selectors';

const SELECTORS = PAGE.overview.activationModal;

module('Integration | Component | Secrets::SyncActivationModal', function (hooks) {
  setupRenderingTest(hooks);
  setupEngine(hooks, 'sync');
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.onClose = sinon.stub();
    this.onError = sinon.stub();

    this.renderComponent = async () => {
      await render(
        hbs`
      <Secrets::SyncActivationModal @onClose={{this.onClose}} @onError={{this.onError}} />
    `,
        { owner: this.engine }
      );
    };
  });

  test('it renders with correct text', async function (assert) {
    await this.renderComponent();

    assert
      .dom(SELECTORS.container)
      .hasTextContaining(
        "By enabling the Secrets Sync feature you may incur additional costs. Please review our documentation to learn more. I've read the above linked document"
      );
  });

  test('it calls onClose', async function (assert) {
    await this.renderComponent();

    await click(SELECTORS.cancel);

    assert.true(this.onClose.called);
  });

  test('it disables submit until user has confirmed docs', async function (assert) {
    await this.renderComponent();

    assert.dom(SELECTORS.checkbox).isNotChecked('checkbox is initially unchecked');
    assert.dom(SELECTORS.confirm).isDisabled('submit is disabled');
    await click(SELECTORS.checkbox);

    assert.dom(SELECTORS.checkbox).isChecked();
    assert.dom(SELECTORS.confirm).isNotDisabled('submit is enabled once user has confirmed');
  });

  module('on submit', function (hooks) {
    hooks.beforeEach(function () {
      const router = this.owner.lookup('service:router');
      this.transitionStub = sinon.stub(router, 'transitionTo');
    });

    module('success', function (hooks) {
      hooks.beforeEach(function () {
        this.server.post('/sys/activation-flags/secrets-sync/activate', () => {
          return {};
        });
      });

      test('it calls the activate endpoint', async function (assert) {
        assert.expect(1);

        this.server.post('/sys/activation-flags/secrets-sync/activate', () => {
          assert.true(true, 'POST to secrets-sync/activate is called');
          return {};
        });

        await this.renderComponent();

        await click(SELECTORS.checkbox);
        await click(SELECTORS.confirm);
      });

      test('it transitions back to sync overview', async function (assert) {
        await this.renderComponent();

        await click(SELECTORS.checkbox);
        await click(SELECTORS.confirm);

        assert.true(this.transitionStub.calledWith('vault.cluster.sync.secrets.overview'));
      });
    });

    module('on error', function (hooks) {
      hooks.beforeEach(function () {
        this.server.post('/sys/activation-flags/secrets-sync/activate', () => {
          return overrideResponse(403);
        });

        const flashMessages = this.owner.lookup('service:flash-messages');
        this.flashDangerSpy = sinon.spy(flashMessages, 'danger');
      });

      test('it handles errors', async function (assert) {
        await this.renderComponent();

        await click(SELECTORS.checkbox);
        await click(SELECTORS.confirm);

        assert.true(this.onError.called, 'calls the onError arg');
        assert.true(this.flashDangerSpy.called, 'triggers an error flash message');
      });
    });
  });
});
