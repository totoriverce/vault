import { currentRouteName } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import page from 'vault/tests/pages/settings/auth/enable';
import listPage from 'vault/tests/pages/access/methods';

module('Acceptance | settings/auth/enable', function(hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(function() {
    return authLogin();
  });

  test('it mounts and redirects', function(assert) {
    // always force the new mount to the top of the list
    const path = `approle-${new Date().getTime()}`;
    const type = 'approle';
    page.visit();
    assert.equal(currentRouteName(), 'vault.cluster.settings.auth.enable');
    page.form.mount(type, path);
    assert.equal(
      page.flash.latestMessage,
      `Successfully mounted ${type} auth method at ${path}.`,
      'success flash shows'
    );
    assert.equal(
      currentRouteName(),
      'vault.cluster.access.methods',
      'redirects to the auth backend list page'
    );
    assert.ok(listPage.backendLinks().findById(path), 'mount is present in the list');
  });
});
