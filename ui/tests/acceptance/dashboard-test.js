/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { module, test } from 'qunit';
import { visit, currentURL, click, fillIn } from '@ember/test-helpers';
import { setupApplicationTest } from 'vault/tests/helpers';
import { setupMirage } from 'ember-cli-mirage/test-support';
import { selectChoose, selectSearch } from 'ember-power-select/test-support/helpers';

import authPage from 'vault/tests/pages/auth';
import SECRETS_ENGINE_SELECTORS from 'vault/tests/helpers/components/dashboard/secrets-engines-card';
import VAULT_CONFIGURATION_SELECTORS from 'vault/tests/helpers/components/dashboard/vault-configuration-details-card';
import mountSecrets from 'vault/tests/pages/settings/mount-secret-backend';
import { runCommands } from 'vault/tests/helpers/pki/pki-run-commands';

module('Acceptance | landing page dashboard', function (hooks) {
  setupApplicationTest(hooks);
  setupMirage(hooks);

  hooks.beforeEach(function () {
    this.data = {
      api_addr: 'http://127.0.0.1:8200',
      cache_size: 0,
      cluster_addr: 'https://127.0.0.1:8201',
      cluster_cipher_suites: '',
      cluster_name: '',
      default_lease_ttl: 0,
      default_max_request_duration: 0,
      detect_deadlocks: '',
      disable_cache: false,
      disable_clustering: false,
      disable_indexing: false,
      disable_mlock: true,
      disable_performance_standby: false,
      disable_printable_check: false,
      disable_sealwrap: false,
      disable_sentinel_trace: false,
      enable_response_header_hostname: false,
      enable_response_header_raft_node_id: false,
      enable_ui: true,
      experiments: null,
      introspection_endpoint: false,
      listeners: [
        {
          config: {
            address: '0.0.0.0:8200',
            cluster_address: '0.0.0.0:8201',
            tls_disable: true,
          },
          type: 'tcp',
        },
      ],
      log_format: '',
      log_level: 'debug',
      log_requests_level: '',
      max_lease_ttl: '48h',
      pid_file: '',
      plugin_directory: '',
      plugin_file_permissions: 0,
      plugin_file_uid: 0,
      raw_storage_endpoint: true,
      seals: [
        {
          disabled: false,
          type: 'shamir',
        },
      ],
      storage: {
        cluster_addr: 'https://127.0.0.1:8201',
        disable_clustering: false,
        raft: {
          max_entry_size: '',
        },
        redirect_addr: 'http://127.0.0.1:8200',
        type: 'raft',
      },
      telemetry: {
        add_lease_metrics_namespace_labels: false,
        circonus_api_app: '',
        circonus_api_token: '',
        circonus_api_url: '',
        circonus_broker_id: '',
        circonus_broker_select_tag: '',
        circonus_check_display_name: '',
        circonus_check_force_metric_activation: '',
        circonus_check_id: '',
        circonus_check_instance_id: '',
        circonus_check_search_tag: '',
        circonus_check_tags: '',
        circonus_submission_interval: '',
        circonus_submission_url: '',
        disable_hostname: true,
        dogstatsd_addr: '',
        dogstatsd_tags: null,
        lease_metrics_epsilon: 3600000000000,
        maximum_gauge_cardinality: 500,
        metrics_prefix: '',
        num_lease_metrics_buckets: 168,
        prometheus_retention_time: 86400000000000,
        stackdriver_debug_logs: false,
        stackdriver_location: '',
        stackdriver_namespace: '',
        stackdriver_project_id: '',
        statsd_address: '',
        statsite_address: '',
        usage_gauge_period: 5000000000,
      },
    };
    return authPage.login();
  });

  // TODO LANDING PAGE: create a test that will navigate to dashboard if user opts into new dashboard ui
  test('navigate to dashboard on login', async function (assert) {
    assert.strictEqual(currentURL(), '/vault/dashboard');
  });

  test('display the version number for the title', async function (assert) {
    await visit('/vault/dashboard');
    assert.dom('[data-test-dashboard-version-header]').hasText('Vault v1.9.0');
  });

  module('secrets engines card', function () {
    test('shows a secrets engine card', async function (assert) {
      await mountSecrets.enable('pki', 'a-pki');
      await visit('/vault/dashboard');
      assert.dom(SECRETS_ENGINE_SELECTORS.cardTitle).hasText('Secrets engines');
      assert.dom(SECRETS_ENGINE_SELECTORS.getSecretEngineAccessor('a-pki')).exists();
    });

    test('it adds disabled css styling to unsupported secret engines', async function (assert) {
      await mountSecrets.enable('nomad', 'nomad');
      await visit('/vault/dashboard');
      assert.dom('[data-test-secrets-engines-row="nomad"] [data-test-view]').doesNotExist();
    });
  });

  module('learn more card', function () {
    test('shows the learn more card', async function (assert) {
      await visit('/vault/dashboard');
      assert.dom('[data-test-learn-more-title]').hasText('Learn more');
      assert
        .dom('[data-test-learn-more-subtext]')
        .hasText(
          'Explore the features of Vault and learn advance practices with the following tutorials and documentation.'
        );
      assert.dom('[data-test-learn-more-links] a').exists({ count: 4 });
    });
  });

  module('configuration details card', function () {
    test('shows the configuration details card', async function (assert) {
      this.server.get('sys/config/state/sanitized', () => ({
        data: this.data,
        wrap_info: null,
        warnings: null,
        auth: null,
      }));
      await authPage.login();
      await visit('/vault/dashboard');
      assert.dom(VAULT_CONFIGURATION_SELECTORS.cardTitle).hasText('Configuration details');
      assert.dom(VAULT_CONFIGURATION_SELECTORS.apiAddr).hasText('http://127.0.0.1:8200');
      assert.dom(VAULT_CONFIGURATION_SELECTORS.defaultLeaseTtl).hasText('0');
      assert.dom(VAULT_CONFIGURATION_SELECTORS.maxLeaseTtl).hasText('2 days');
      assert.dom(VAULT_CONFIGURATION_SELECTORS.tlsDisable).hasText('Enabled');
      assert.dom(VAULT_CONFIGURATION_SELECTORS.logFormat).hasText('None');
      assert.dom(VAULT_CONFIGURATION_SELECTORS.logLevel).hasText('debug');
      assert.dom(VAULT_CONFIGURATION_SELECTORS.storageType).hasText('raft');
    });
    test('shows the tls disabled if it is disabled', async function (assert) {
      this.server.get('sys/config/state/sanitized', () => {
        this.data.listeners[0].config.tls_disable = false;
        return {
          data: this.data,
          wrap_info: null,
          warnings: null,
          auth: null,
        };
      });
      await authPage.login();
      await visit('/vault/dashboard');
      assert.dom(VAULT_CONFIGURATION_SELECTORS.tlsDisable).hasText('Disabled');
    });
    test('shows the tls disabled if there is no tlsDisabled returned from server', async function (assert) {
      this.server.get('sys/config/state/sanitized', () => {
        this.data.listeners = [];

        return {
          data: this.data,
          wrap_info: null,
          warnings: null,
          auth: null,
        };
      });
      await authPage.login();
      await visit('/vault/dashboard');
      assert.dom(VAULT_CONFIGURATION_SELECTORS.tlsDisable).hasText('Disabled');
    });
  });
  module('quick actions card', function () {
    test('shows the quick actions card empty state when no engine is selected', async function (assert) {
      await authPage.login();
      assert.dom('[data-test-component="empty-state"]').exists({ count: 1 });
      assert.dom('[data-test-empty-state-title]').hasText('No mount selected');
      await selectChoose('.search-select', 'pki');
      assert.dom('[data-test-component="empty-state"]').doesNotExist();
    });
    test('shows the quick actions card for pki', async function (assert) {
      await authPage.login();
      await mountSecrets.enable('pki', 'b-pki');
      await runCommands([`write b-pki/root/generate/internal common_name="Hashicorp Test"`]);
      await runCommands([
        `write b-pki/roles/some-role \
      issuer_ref="default" \
      allowed_domains="example.com" \
      allow_subdomains=true \
      max_ttl="720h"`,
      ]);
      await visit('/vault/dashboard');
      await selectChoose('.search-select', 'b-pki');
      await fillIn('[data-test-select="action-select"]', 'Issue certificate');
      await selectChoose('.search-select', 'some-role');
      await click('[data-test-button="Issue leaf certificate"]');
      assert.strictEqual(currentURL(), `/vault/secrets/b-pki/pki/roles/some-role/generate`);
    });
    test('shows the quick actions card for db', async function (assert) {
      await authPage.login();
      await mountSecrets.enable('database', 'database');
      await runCommands([
        `vault write  database/roles/api-prod db_name=apiprod creation_statements="CREATE ROLE "{{name}}" WITH LOGIN PASSWORD '{{password}}' VALID UNTIL '{{expiration}}'; GRANT SELECT ON ALL TABLES IN SCHEMA public TO "{{name}}";" default_ttl=1h max_ttl=24h`,
      ]);
      await visit('/vault/dashboard');
      await selectChoose('.search-select', 'database');
      await fillIn('[data-test-select="action-select"]', 'Generate credentials for database');
      await selectChoose('.search-select', 'api-prod');
      await click('[data-test-button="Generate credentials"]');
      assert.strictEqual(currentURL(), `/vault/secrets/database/credentials/api-prod`);
    });
    test('shows the quick actions card for kv', async function (assert) {
      await authPage.login();
      await mountSecrets.enable('kv', 'kv');
      await visit('/vault/dashboard');
      // TODO: write more kv tests when kv work is merged!
      await selectChoose('.search-select', 'kv');
      await fillIn('[data-test-select="action-select"]', 'Find KV secrets');
      assert.dom('[data-test-button="Read secrets"]').exists({ count: 1 });
    });
    test('hides engines that are not pki, db, or kv for quick actions card', async function (assert) {
      await authPage.login();
      await mountSecrets.enable('consul', 'consul');
      await visit('/vault/dashboard');
      await selectSearch('[data-test-secrets-engines-select]', 'consul');
      assert.dom('.ember-power-select-option--no-matches-message').exists({ count: 1 });
    });
  });
});
