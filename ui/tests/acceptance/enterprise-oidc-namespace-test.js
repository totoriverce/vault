import { visit, currentURL } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { create } from 'ember-cli-page-object';
import { setupMirage } from 'ember-cli-mirage/test-support';
import parseURL from 'core/utils/parse-url';
import consoleClass from 'vault/tests/pages/components/console/ui-panel';
import authPage from 'vault/tests/pages/auth';

const shell = create(consoleClass);

const createNS = async (name) => {
  await shell.runCommands(`write sys/namespaces/${name} -force`);
};
const SELECTORS = {
  authTab: (path) => `[data-test-auth-method="${path}"] a`,
};

module('Acceptance | Enterprise | oidc auth namespace test', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(async function () {
    this.namespace = 'test-ns';
    this.rootOidc = 'root-oidc';
    this.nsOidc = 'ns-oidc';

    this.server.post(`/auth/:path/config`, () => {});

    this.enableOidc = (path, role = '') => {
      return shell.runCommands([
        `write sys/auth/${path} type=oidc`,
        `write auth/${path}/config default_role="${role}" oidc_discovery_url="https://example.com"`,
        // show method as tab
        `write sys/auth/${path}/tune listing_visibility="unauth"`,
      ]);
    };

    this.disableOidc = (path) => shell.runCommands([`delete /sys/auth/${path}`]);
  });

  test('oidc: request is made to auth_url when a namespace is inputted', async function (assert) {
    assert.expect(5);

    this.server.post(`/auth/${this.rootOidc}/oidc/auth_url`, (schema, req) => {
      const { redirect_uri } = JSON.parse(req.requestBody);
      const { pathname, search } = parseURL(redirect_uri);
      assert.strictEqual(
        pathname + search,
        `/ui/vault/auth/${this.rootOidc}/oidc/callback`,
        'request made to auth_url when the login page is visited'
      );
    });
    this.server.post(`/auth/${this.nsOidc}/oidc/auth_url`, (schema, req) => {
      const { redirect_uri } = JSON.parse(req.requestBody);
      const { pathname, search } = parseURL(redirect_uri);
      assert.strictEqual(
        pathname + search,
        `/ui/vault/auth/${this.nsOidc}/oidc/callback?namespace=${this.namespace}`,
        'request made to correct auth_url when namespace is filled in'
      );
    });

    await authPage.login();
    // enable oidc in root namespace, without default role
    await this.enableOidc(this.rootOidc);
    // create child namespace to enable oidc
    await createNS(this.namespace);
    // enable oidc in child namespace with default role
    await authPage.loginNs(this.namespace);
    await this.enableOidc(this.nsOidc, `${this.nsOidc}-role`);
    await authPage.logout();

    await visit('/vault/auth');
    assert.dom(SELECTORS.authTab(this.rootOidc)).exists('renders oidc method tab for root');
    await authPage.namespaceInput(this.namespace);
    assert.strictEqual(
      currentURL(),
      `/vault/auth?namespace=${this.namespace}&with=${this.nsOidc}%2F`,
      'url updates with namespace value'
    );
    assert.dom(SELECTORS.authTab(this.nsOidc)).exists('renders oidc method tab for child namespace');

    // disable methods to cleanup test state for re-running
    await authPage.login();
    await this.disableOidc(this.rootOidc);
    await this.disableOidc(this.nsOidc);
    await shell.runCommands([`delete /sys/auth/${this.namespace}`]);
    await authPage.logout();
  });
});
