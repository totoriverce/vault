import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import { currentURL, settled, click, visit, fillIn } from '@ember/test-helpers';
import { create } from 'ember-cli-page-object';
import { selectChoose, clickTrigger } from 'ember-power-select/test-support/helpers';

import mountSecrets from 'vault/tests/pages/settings/mount-secret-backend';
import connectionPage from 'vault/tests/pages/secrets/backend/database/connection';
import rolePage from 'vault/tests/pages/secrets/backend/database/role';
import apiStub from 'vault/tests/helpers/noop-all-api-requests';
import authPage from 'vault/tests/pages/auth';
import logout from 'vault/tests/pages/logout';
import consoleClass from 'vault/tests/pages/components/console/ui-panel';
import searchSelect from 'vault/tests/pages/components/search-select';

const searchSelectComponent = create(searchSelect);

const consoleComponent = create(consoleClass);

const MODEL = {
  engineType: 'database',
  id: 'database-name',
};

const mount = async () => {
  let path = `database-${Date.now()}`;
  await mountSecrets.enable('database', path);
  await settled();
  return path;
};

const newConnection = async (backend, plugin = 'mongodb-database-plugin') => {
  const name = `connection-${Date.now()}`;
  await connectionPage.visitCreate({ backend });
  await connectionPage.dbPlugin(plugin);
  await connectionPage.name(name);
  await connectionPage.url(`mongodb://127.0.0.1:4321/${name}`);
  await connectionPage.toggleVerify();
  await connectionPage.save();
  await connectionPage.enable();
  return name;
};

const connectionTests = [
  {
    name: 'mongodb-connection',
    plugin: 'mongodb-database-plugin',
    url: `mongodb://127.0.0.1:4321/test`,
    requiredFields: async (assert, name) => {
      assert.dom('[data-test-input="username"]').exists(`Username field exists for ${name}`);
      assert.dom('[data-test-input="password"]').exists(`Password field exists for ${name}`);
      assert.dom('[data-test-input="write_concern"]').exists(`Write concern field exists for ${name}`);
      assert.dom('[data-test-toggle-group="TLS options"]').exists('TLS options toggle exists');
      assert
        .dom('[data-test-input="root_rotation_statements"]')
        .exists(`Root rotation statements exists for ${name}`);
    },
  },
  {
    name: 'mssql-connection',
    plugin: 'mssql-database-plugin',
    url: `mssql://127.0.0.1:4321/test`,
    requiredFields: async (assert, name) => {
      assert.dom('[data-test-input="username"]').exists(`Username field exists for ${name}`);
      assert.dom('[data-test-input="password"]').exists(`Password field exists for ${name}`);
      assert
        .dom('[data-test-input="max_open_connections"]')
        .exists(`Max open connections exists for ${name}`);
      assert
        .dom('[data-test-input="max_idle_connections"]')
        .exists(`Max idle connections exists for ${name}`);
      assert
        .dom('[data-test-input="max_connection_lifetime"]')
        .exists(`Max connection lifetime exists for ${name}`);
      assert
        .dom('[data-test-input="root_rotation_statements"]')
        .exists(`Root rotation statements exists for ${name}`);
    },
  },
  {
    name: 'mysql-connection',
    plugin: 'mysql-database-plugin',
    url: `{{username}}:{{password}}@tcp(127.0.0.1:3306)/test`,
    requiredFields: async (assert, name) => {
      assert.dom('[data-test-input="username"]').exists(`Username field exists for ${name}`);
      assert.dom('[data-test-input="password"]').exists(`Password field exists for ${name}`);
      assert
        .dom('[data-test-input="max_open_connections"]')
        .exists(`Max open connections exists for ${name}`);
      assert
        .dom('[data-test-input="max_idle_connections"]')
        .exists(`Max idle connections exists for ${name}`);
      assert
        .dom('[data-test-input="max_connection_lifetime"]')
        .exists(`Max connection lifetime exists for ${name}`);
      assert.dom('[data-test-toggle-group="TLS options"]').exists('TLS options toggle exists');
      assert
        .dom('[data-test-input="root_rotation_statements"]')
        .exists(`Root rotation statements exists for ${name}`);
    },
  },
  {
    name: 'mysql-aurora-connection',
    plugin: 'mysql-aurora-database-plugin',
    url: `{{username}}:{{password}}@tcp(127.0.0.1:3306)/test`,
    requiredFields: async (assert, name) => {
      assert.dom('[data-test-input="username"]').exists(`Username field exists for ${name}`);
      assert.dom('[data-test-input="password"]').exists(`Password field exists for ${name}`);
      assert
        .dom('[data-test-input="max_open_connections"]')
        .exists(`Max open connections exists for ${name}`);
      assert
        .dom('[data-test-input="max_idle_connections"]')
        .exists(`Max idle connections exists for ${name}`);
      assert
        .dom('[data-test-input="max_connection_lifetime"]')
        .exists(`Max connection lifetime exists for ${name}`);
      assert.dom('[data-test-toggle-group="TLS options"]').exists('TLS options toggle exists');
      assert
        .dom('[data-test-input="root_rotation_statements"]')
        .exists(`Root rotation statements exists for ${name}`);
    },
  },
  {
    name: 'mysql-rds-connection',
    plugin: 'mysql-rds-database-plugin',
    url: `{{username}}:{{password}}@tcp(127.0.0.1:3306)/test`,
    requiredFields: async (assert, name) => {
      assert.dom('[data-test-input="username"]').exists(`Username field exists for ${name}`);
      assert.dom('[data-test-input="password"]').exists(`Password field exists for ${name}`);
      assert
        .dom('[data-test-input="max_open_connections"]')
        .exists(`Max open connections exists for ${name}`);
      assert
        .dom('[data-test-input="max_idle_connections"]')
        .exists(`Max idle connections exists for ${name}`);
      assert
        .dom('[data-test-input="max_connection_lifetime"]')
        .exists(`Max connection lifetime exists for ${name}`);
      assert.dom('[data-test-toggle-group="TLS options"]').exists('TLS options toggle exists');
      assert
        .dom('[data-test-input="root_rotation_statements"]')
        .exists(`Root rotation statements exists for ${name}`);
    },
  },
  {
    name: 'mysql-legacy-connection',
    plugin: 'mysql-legacy-database-plugin',
    url: `{{username}}:{{password}}@tcp(127.0.0.1:3306)/test`,
    requiredFields: async (assert, name) => {
      assert.dom('[data-test-input="username"]').exists(`Username field exists for ${name}`);
      assert.dom('[data-test-input="password"]').exists(`Password field exists for ${name}`);
      assert
        .dom('[data-test-input="max_open_connections"]')
        .exists(`Max open connections exists for ${name}`);
      assert
        .dom('[data-test-input="max_idle_connections"]')
        .exists(`Max idle connections exists for ${name}`);
      assert
        .dom('[data-test-input="max_connection_lifetime"]')
        .exists(`Max connection lifetime exists for ${name}`);
      assert.dom('[data-test-toggle-group="TLS options"]').exists('TLS options toggle exists');
      assert
        .dom('[data-test-input="root_rotation_statements"]')
        .exists(`Root rotation statements exists for ${name}`);
    },
  },
];

module('Acceptance | secrets/database/*', function(hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function() {
    this.server = apiStub({ usePassthrough: true });
    return authPage.login();
  });
  hooks.afterEach(function() {
    this.server.shutdown();
  });

  test('can enable the database secrets engine', async function(assert) {
    let backend = `database-${Date.now()}`;
    await mountSecrets.enable('database', backend);
    await settled();
    assert.equal(
      currentURL(),
      `/vault/secrets/${backend}/list`,
      'Mounts and redirects to connection list page'
    );
    assert.dom('[data-test-component="empty-state"]').exists('Empty state exists');
    assert
      .dom('.is-active[data-test-secret-list-tab="Connections"]')
      .exists('Has Connections tab which is active');
    await click('[data-test-tab="overview"]');
    assert.equal(currentURL(), `/vault/secrets/${backend}/overview`, 'Tab links to overview page');
    assert.dom('[data-test-component="empty-state"]').exists('Empty state also exists on overview page');
    assert.dom('[data-test-secret-list-tab="Roles"]').exists('Has Roles tab');
  });

  test('Connection create and edit form for each plugin', async function(assert) {
    const backend = await mount();
    for (let testCase of connectionTests) {
      await connectionPage.visitCreate({ backend });
      assert.equal(currentURL(), `/vault/secrets/${backend}/create`, 'Correct creation URL');
      assert
        .dom('[data-test-empty-state-title')
        .hasText('No plugin selected', 'No plugin is selected by default and empty state shows');
      await connectionPage.dbPlugin(testCase.plugin);
      assert.dom('[data-test-empty-state]').doesNotExist('Empty state goes away after plugin selected');
      await connectionPage.name(testCase.name);
      await connectionPage.url(testCase.url);
      testCase.requiredFields(assert, testCase.name);
      await connectionPage.toggleVerify();
      await connectionPage.save();
      await connectionPage.enable();
      assert
        .dom('[data-test-modal-title]')
        .hasText('Rotate your root credentials?', 'Modal appears asking to rotate root credentials');
      await connectionPage.enable();
      assert.ok(
        currentURL().startsWith(`/vault/secrets/${backend}/show/${testCase.name}`),
        `Saves connection and takes you to show page for ${testCase.name}`
      );
      assert
        .dom(`[data-test-row-value="Password"]`)
        .doesNotExist(`Does not show Password value on show page for ${testCase.name}`);
      await connectionPage.edit();
      assert.ok(
        currentURL().startsWith(`/vault/secrets/${backend}/edit/${testCase.name}`),
        `Edit connection button and takes you to edit page for ${testCase.name}`
      );
      assert.dom(`[data-test-input="name"]`).hasAttribute('readonly');
      assert.dom(`[data-test-input="plugin_name"]`).hasAttribute('readonly');
      assert.dom('[data-test-input="password"]').doesNotExist('Password is not displayed on edit form');
      assert.dom('[data-test-toggle-input="show-password"]').exists('Update password toggle exists');
      await connectionPage.toggleVerify();
      await connectionPage.save();
      // click "Add Role"
      await connectionPage.addRole();
      await settled();
      assert.equal(
        searchSelectComponent.selectedOptions[0].text,
        testCase.name,
        'Database connection is pre-selected on the form'
      );
      await click('[data-test-secret-breadcrumb]');
    }
  });

  test('Can create and delete a connection', async function(assert) {
    const backend = await mount();
    const connectionDetails = {
      plugin: 'mongodb-database-plugin',
      id: 'horses-db',
      fields: [
        { label: 'Connection Name', name: 'name', value: 'horses-db' },
        { label: 'Connection url', name: 'connection_url', value: 'mongodb://127.0.0.1:235/horses' },
        { label: 'Username', name: 'username', value: 'user', hideOnShow: true },
        { label: 'Password', name: 'password', password: 'so-secure', hideOnShow: true },
        { label: 'Write concern', name: 'write_concern' },
      ],
    };
    assert.equal(
      currentURL(),
      `/vault/secrets/${backend}/list`,
      'Mounts and redirects to connection list page'
    );
    await connectionPage.createLink();
    assert.equal(currentURL(), `/vault/secrets/${backend}/create`, 'Create link goes to create page');
    assert
      .dom('[data-test-empty-state-title')
      .hasText('No plugin selected', 'No plugin is selected by default and empty state shows');
    await connectionPage.dbPlugin(connectionDetails.plugin);
    assert.dom('[data-test-empty-state]').doesNotExist('Empty state goes away after plugin selected');
    connectionDetails.fields.forEach(async ({ name, value }) => {
      assert
        .dom(`[data-test-input="${name}"]`)
        .exists(`Field ${name} exists for ${connectionDetails.plugin}`);
      if (value) {
        await fillIn(`[data-test-input="${name}"]`, value);
      }
    });
    // uncheck verify for the save step to work
    await connectionPage.toggleVerify();
    await connectionPage.save();
    assert
      .dom('[data-test-modal-title]')
      .hasText('Rotate your root credentials?', 'Modal appears asking to ');
    await connectionPage.enable();
    assert.equal(
      currentURL(),
      `/vault/secrets/${backend}/show/${connectionDetails.id}`,
      'Saves connection and takes you to show page'
    );
    connectionDetails.fields.forEach(({ label, name, value, hideOnShow }) => {
      if (hideOnShow) {
        assert
          .dom(`[data-test-row-value="${label}"]`)
          .doesNotExist(`Does not show ${name} value on show page for ${connectionDetails.plugin}`);
      } else if (!value) {
        assert.dom(`[data-test-row-value="${label}"]`).hasText('Default');
      } else {
        assert.dom(`[data-test-row-value="${label}"]`).hasText(value);
      }
    });
    await connectionPage.delete();
    assert.equal(currentURL(), `/vault/secrets/${backend}/list`, 'Redirects to connection list page');
    assert
      .dom('[data-test-empty-state-title')
      .hasText('No connections in this backend', 'No connections listed because it was deleted');
  });

  test('buttons show up for managing connection', async function(assert) {
    const backend = await mount();
    const connection = await newConnection(backend);
    await connectionPage.visitShow({ backend, id: connection });
    assert
      .dom('[data-test-database-connection-delete]')
      .hasText('Delete connection', 'Delete connection button exists with correct text');
    assert
      .dom('[data-test-database-connection-reset]')
      .hasText('Reset connection', 'Reset button exists with correct text');
    assert.dom('[data-test-secret-create]').hasText('Add role', 'Add role button exists with correct text');
    assert.dom('[data-test-edit-link]').hasText('Edit configuration', 'Edit button exists with correct text');
    const CONNECTION_VIEW_ONLY = `
path "${backend}/*" {
  capabilities = ["deny"]
}
path "${backend}/config" {
  capabilities = ["list"]
}
path "${backend}/config/*" {
  capabilities = ["read"]
}
    `;
    await consoleComponent.runCommands([
      `write sys/mounts/${backend} type=database`,
      `write sys/policies/acl/test-policy policy=${btoa(CONNECTION_VIEW_ONLY)}`,
      'write -field=client_token auth/token/create policies=test-policy ttl=1h',
    ]);
    let token = consoleComponent.lastTextOutput;
    await logout.visit();
    await authPage.login(token);
    await connectionPage.visitShow({ backend, id: connection });
    assert.equal(currentURL(), `/vault/secrets/${backend}/show/${connection}`, 'Allows reading connection');
    assert
      .dom('[data-test-database-connection-delete]')
      .doesNotExist('Delete button does not show due to permissions');
    assert
      .dom('[data-test-database-connection-reset]')
      .doesNotExist('Reset button does not show due to permissions');
    assert.dom('[data-test-secret-create]').doesNotExist('Add role button does not show due to permissions');
    assert.dom('[data-test-edit-link]').doesNotExist('Edit button does not show due to permissions');
    await visit(`/vault/secrets/${backend}/overview`);
    assert.dom('[data-test-selectable-card="Connections"]').exists('Connections card exists on overview');
    assert
      .dom('[data-test-selectable-card="Roles"]')
      .doesNotExist('Roles card does not exist on overview w/ policy');
    assert.dom('.title-number').hasText('1', 'Lists the correct number of connections');
  });

  test('Role create form', async function(assert) {
    const backend = await mount();
    // Connection needed for role fields
    await newConnection(backend);
    await rolePage.visitCreate({ backend });
    await rolePage.name('bar');
    assert
      .dom('[data-test-component="empty-state"]')
      .exists({ count: 2 }, 'Two empty states exist before selections made');
    await clickTrigger('#database');
    assert.equal(searchSelectComponent.options.length, 1, 'list shows existing connections so far');
    await selectChoose('#database', '.ember-power-select-option', 0);
    assert
      .dom('[data-test-component="empty-state"]')
      .exists({ count: 2 }, 'Two empty states exist before selections made');
    await rolePage.roleType('static');
    assert.dom('[data-test-component="empty-state"]').doesNotExist('Empty states go away');
    assert.dom('[data-test-input="username"]').exists('Username field appears for static role');
    assert
      .dom('[data-test-toggle-input="Rotation period"]')
      .exists('Rotation period field appears for static role');
    await rolePage.roleType('dynamic');
    assert
      .dom('[data-test-toggle-input="Generated credentials’s Time-to-Live (TTL)"]')
      .exists('TTL field exists for dynamic');
    assert
      .dom('[data-test-toggle-input="Generated credentials’s maximum Time-to-Live (Max TTL)"]')
      .exists('Max TTL field exists for dynamic');
    // Real connection (actual running db) required to save role, so we aren't testing that flow yet
  });

  test('root and limited access', async function(assert) {
    this.set('model', MODEL);
    let backend = 'database';
    const NO_ROLES_POLICY = `
      path "database/roles/*" {
        capabilities = ["delete"]
      }
      path "database/static-roles/*" {
        capabilities = ["delete"]
      }
      path "database/config/*" {
        capabilities = ["list", "create", "read", "update"]
      }
      path "database/creds/*" {
        capabilities = ["list", "create", "read", "update"]
      }
    `;
    await consoleComponent.runCommands([
      `write sys/mounts/${backend} type=database`,
      `write sys/policies/acl/test-policy policy=${btoa(NO_ROLES_POLICY)}`,
      'write -field=client_token auth/token/create policies=test-policy ttl=1h',
    ]);
    let token = consoleComponent.lastTextOutput;

    // test root user flow
    await settled();

    // await click('[data-test-secret-backend-row="database"]');
    // skipping the click because occasionally is shows up on the second page and cannot be found
    await visit(`/vault/secrets/database/overview`);
    await settled();
    assert.dom('[data-test-component="empty-state"]').exists('renders empty state');
    assert.dom('[data-test-secret-list-tab="Connections"]').exists('renders connections tab');
    assert.dom('[data-test-secret-list-tab="Roles"]').exists('renders connections tab');

    await click('[data-test-secret-create="connections"]');
    assert.equal(currentURL(), '/vault/secrets/database/create');

    // Login with restricted policy
    await logout.visit();
    await authPage.login(token);
    await settled();
    // skipping the click because occasionally is shows up on the second page and cannot be found
    await visit(`/vault/secrets/database/overview`);
    assert.dom('[data-test-tab="overview"]').exists('renders overview tab');
    assert.dom('[data-test-secret-list-tab="Connections"]').exists('renders connections tab');
    assert
      .dom('[data-test-secret-list-tab="Roles]')
      .doesNotExist(`does not show the roles tab because it does not have permissions`);
    assert
      .dom('[data-test-selectable-card="Connections"]')
      .exists({ count: 1 }, 'renders only the connection card');

    await click('[data-test-action-text="Configure new"]');
    assert.equal(currentURL(), '/vault/secrets/database/create?itemType=connection');
  });
});
