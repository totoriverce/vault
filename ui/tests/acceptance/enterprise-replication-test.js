import { clickTrigger } from 'ember-power-select/test-support/helpers';
import { click, fillIn, findAll, currentURL, find, visit, settled } from '@ember/test-helpers';
import { module, test } from 'qunit';
import { setupApplicationTest } from 'ember-qunit';
import authPage from 'vault/tests/pages/auth';
import { pollCluster } from 'vault/tests/helpers/poll-cluster';
import { create } from 'ember-cli-page-object';
import flashMessage from 'vault/tests/pages/components/flash-message';
import ss from 'vault/tests/pages/components/search-select';

const searchSelect = create(ss);
const flash = create(flashMessage);

const disableReplication = async (type, assert) => {
  // disable performance replication
  await visit(`/vault/replication/${type}`);
  await settled();
  if (findAll('[data-test-replication-link="manage"]').length) {
    await click('[data-test-replication-link="manage"]');
    await settled();
    await click('[data-test-disable-replication] button');

    const typeDisplay = type === 'dr' ? 'Disaster Recovery' : 'Performance';
    await fillIn('[data-test-confirmation-modal-input="disable"]', typeDisplay);
    await click('[data-test-confirm-button]');
    await settled();
    if (assert) {
      assert.equal(currentURL(), `/vault/replication`, 'redirects to the replication page');
      assert.equal(
        flash.latestMessage,
        'This cluster is having replication disabled. Vault will be unavailable for a brief period and will resume service shortly.',
        'renders info flash when disabled'
      );
    }
    await settled();
  }
};

module('Acceptance | Enterprise | replication', function(hooks) {
  setupApplicationTest(hooks);

  hooks.beforeEach(async function() {
    await authPage.login();
    await settled();
    await disableReplication('dr');
    await settled();
    await disableReplication('performance');
    await settled();
  });

  hooks.afterEach(async function() {
    await disableReplication('dr');
    await settled();
    await disableReplication('performance');
    await settled();
  });

  test('replication', async function(assert) {
    const secondaryName = 'firstSecondary';
    const mode = 'deny';
    let mountPath;

    // confirm unable to visit dr secondary details page when both replications are disabled
    await visit('/vault/replication-dr-promote/details');
    await settled();
    assert.dom('[data-test-component="empty-state"]').exists();
    assert
      .dom('[data-test-empty-state-title]')
      .includesText('Disaster Recovery secondary not set up', 'shows the correct title of the empty state');

    assert
      .dom('[data-test-empty-state-message]')
      .hasText(
        'This cluster has not been enabled as a Disaster Recovery Secondary. You can do so by enabling replication and adding a secondary from the Disaster Recovery Primary.',
        'renders default message specific to when no replication is enabled'
      );

    await visit('/vault/replication');
    await settled();
    assert.equal(currentURL(), '/vault/replication');

    // enable perf replication
    await click('[data-test-replication-type-select="performance"]');
    await settled();
    await fillIn('[data-test-replication-cluster-mode-select]', 'primary');

    await click('[data-test-replication-enable]');
    await settled();
    await pollCluster(this.owner);
    await settled();

    // confirm that the details dashboard shows
    assert.dom('[data-test-replication-dashboard]').exists();

    // add a secondary with a mount filter config
    await click('[data-test-replication-link="secondaries"]');
    await settled();
    await click('[data-test-secondary-add]');
    await settled();
    await fillIn('[data-test-replication-secondary-id]', secondaryName);

    await click('#deny');
    await clickTrigger();
    mountPath = searchSelect.options.objectAt(0).text;
    await searchSelect.options.objectAt(0).click();
    await click('[data-test-secondary-add]');
    await settled();

    await pollCluster(this.owner);
    // click into the added secondary's mount filter config
    await click('[data-test-replication-link="secondaries"]');
    await settled();
    await click('[data-test-popup-menu-trigger]');

    await click('[data-test-replication-mount-filter-link]');
    await settled();
    assert.equal(currentURL(), `/vault/replication/performance/secondaries/config/show/${secondaryName}`);
    assert.ok(
      find('[data-test-mount-config-mode]')
        .textContent.trim()
        .toLowerCase()
        .includes(mode),
      'show page renders the correct mode'
    );
    assert
      .dom('[data-test-mount-config-paths]')
      .includesText(mountPath, 'show page renders the correct mount path');

    // delete config by choosing "no filter" in the edit screen
    await click('[data-test-replication-link="edit-mount-config"]');
    await settled();
    await click('#no-filtering');
    await settled();

    await click('[data-test-config-save]');
    await settled();
    assert.equal(
      flash.latestMessage,
      `The performance mount filter config for the secondary ${secondaryName} was successfully deleted.`,
      'renders success flash upon deletion'
    );
    assert.equal(
      currentURL(),
      `/vault/replication/performance/secondaries`,
      'redirects to the secondaries page'
    );
    // nav back to details page and confirm secondary is in the known secondaries table
    await click('[data-test-replication-link="details"]');
    await settled();
    assert
      .dom(`[data-test-secondaries=row-for-${secondaryName}]`)
      .exists('shows a table row the recently added secondary');

    // nav to DR
    await visit('/vault/replication/dr');
    await settled();
    await fillIn('[data-test-replication-cluster-mode-select]', 'secondary');
    assert
      .dom('[data-test-replication-enable]')
      .isDisabled('dr secondary enable is disabled when other replication modes are on');

    // disable performance replication
    await disableReplication('performance', assert);
    await settled();
    await pollCluster(this.owner);

    // enable dr replication
    await visit('vault/replication/dr');
    await settled();
    await fillIn('[data-test-replication-cluster-mode-select]', 'primary');
    await click('button[type="submit"]');
    await settled();

    await pollCluster(this.owner);
    // empty state inside of know secondaries table
    assert.dom('[data-test-empty-state-title]').exists();
    assert
      .dom('[data-test-empty-state-title]')
      .includesText(
        'No known dr secondary clusters associated with this cluster',
        'shows the correct title of the empty state'
      );

    assert.ok(
      find('[data-test-replication-title]').textContent.includes('Disaster Recovery'),
      'it displays the replication type correctly'
    );
    assert.ok(
      find('[data-test-replication-mode-display]').textContent.includes('primary'),
      'it displays the cluster mode correctly'
    );

    // add dr secondary
    await click('[data-test-replication-link="secondaries"]');
    await settled();
    await click('[data-test-secondary-add]');
    await settled();
    await fillIn('[data-test-replication-secondary-id]', secondaryName);

    await click('[data-test-secondary-add]');
    await settled();
    await pollCluster(this.owner);
    await click('[data-test-replication-link="secondaries"]');
    await settled();
    assert
      .dom('[data-test-secondary-name]')
      .includesText(secondaryName, 'it displays the secondary in the list of known secondaries');
  });

  test('disabling dr primary when perf replication is enabled', async function(assert) {
    await visit('vault/replication/performance');
    await settled();
    // enable perf replication
    await fillIn('[data-test-replication-cluster-mode-select]', 'primary');
    await click('[data-test-replication-enable]');
    await settled();
    await pollCluster(this.owner);

    // enable dr replication
    await visit('/vault/replication/dr');
    await settled();
    await fillIn('[data-test-replication-cluster-mode-select]', 'primary');
    await settled();
    await click('[data-test-replication-enable]');
    await settled();
    await pollCluster(this.owner);
    await visit('/vault/replication/dr/manage');
    await settled();
    await click('[data-test-demote-replication] [data-test-replication-action-trigger]');
    await settled();
    assert.ok(findAll('[data-test-demote-warning]').length, 'displays the demotion warning');
  });

  test('navigating to dr secondary details page when dr secondary is not enabled', async function(assert) {
    // enable dr replication

    await visit('/vault/replication/dr');
    await settled();
    await fillIn('[data-test-replication-cluster-mode-select]', 'primary');
    await click('[data-test-replication-enable]');
    await settled();
    await pollCluster(this.owner);
    await visit('/vault/replication-dr-promote/details');
    await settled();

    assert.dom('[data-test-component="empty-state"]').exists();
    assert
      .dom('[data-test-empty-state-message]')
      .hasText(
        'This Disaster Recovery secondary has not been enabled. You can do so from the Disaster Recovery Primary.',
        'renders message when replication is enabled'
      );
  });

  test('add secondary and navigate through token generation modal', async function(assert) {
    const secondaryNameFirst = 'firstSecondary';
    const secondaryNameSecond = 'secondSecondary';
    await visit('/vault/replication');
    await settled();
    // enable perf replication
    await click('[data-test-replication-type-select="performance"]');
    await settled();
    await fillIn('[data-test-replication-cluster-mode-select]', 'primary');
    await click('[data-test-replication-enable]');
    await settled();
    await pollCluster(this.owner);
    await settled();

    // add a secondary with default TTL
    await click('[data-test-replication-link="secondaries"]');
    await settled();
    await click('[data-test-secondary-add]');
    await settled();
    await fillIn('[data-test-replication-secondary-id]', secondaryNameFirst);
    await click('[data-test-secondary-add]');
    await settled();
    await pollCluster(this.owner);
    await settled();
    let modalDefaultTtl = document.querySelector('[data-test-row-value="TTL"]').innerText;
    // checks on secondary token modal
    assert.dom('#modal-wormhole').exists();
    assert.equal(modalDefaultTtl, '1800s', 'shows the correct TTL of 1800s');
    // click off the modal to make sure you don't just have to click on the copy-close button to copy the token
    await click('[data-test-modal-background]');
    await settled();
    // add another secondary not using the default ttl
    await click('[data-test-secondary-add]');
    await settled();
    await fillIn('[data-test-replication-secondary-id]', secondaryNameSecond);
    await click('[data-test-toggle-input]');
    await settled();
    await fillIn('[data-test-ttl-value]', 3);
    await click('[data-test-secondary-add]');
    await settled();
    await pollCluster(this.owner);
    await settled();
    let modalTtl = document.querySelector('[data-test-row-value="TTL"]').innerText;
    assert.equal(modalTtl, '180s', 'shows the correct TTL of 180s');
    await click('[data-test-modal-background]');
    await settled();
    // confirm you were redirected to the secondaries page
    assert.equal(
      currentURL(),
      `/vault/replication/performance/secondaries`,
      'redirects to the secondaries page'
    );
    assert
      .dom('[data-test-secondary-name]')
      .includesText(secondaryNameFirst, 'it displays the secondary in the list of secondaries');
  });

  test('render performance and dr primary and navigate to details page', async function(assert) {
    // enable perf primary replication
    await visit('/vault/replication');
    await click('[data-test-replication-type-select="performance"]');
    await settled();
    await fillIn('[data-test-replication-cluster-mode-select]', 'primary');
    await click('[data-test-replication-enable]');
    await settled();
    await pollCluster(this.owner);
    await settled();

    await visit('/vault/replication');
    await settled();
    assert
      .dom(`[data-test-replication-summary-card]`)
      .doesNotExist(`does not render replication summary card when both modes are not enabled as primary`);

    // enable DR primary replication
    const enableButton = document.querySelector('.is-primary');

    await click(enableButton);
    await settled();
    await click('[data-test-replication-enable="true"]');
    await settled();
    await pollCluster(this.owner);
    await settled();

    // navigate using breadcrumbs back to replication.index
    await click('[data-test-replication-breadcrumb]');
    await settled();
    assert
      .dom('[data-test-replication-summary-card]')
      .exists({ count: 2 }, 'renders two replication-summary-card components');

    // navigate to details page using the "Details" link
    await click('[data-test-manage-link="Disaster Recovery"]');
    await settled();
    assert
      .dom('[data-test-selectable-card-container="primary"]')
      .exists('shows the correct card on the details dashboard');
    assert.equal(currentURL(), '/vault/replication/dr');
  });

  test('render performance secondary and navigate to the details page', async function(assert) {
    // enable perf replication
    await visit('/vault/replication');
    await settled();
    await click('[data-test-replication-type-select="performance"]');
    await settled();
    await fillIn('[data-test-replication-cluster-mode-select]', 'primary');
    await click('[data-test-replication-enable]');
    await settled();
    await pollCluster(this.owner);
    await settled();

    // demote perf primary to a secondary
    await click('[data-test-replication-link="manage"]');
    await settled();
    // open demote modal
    await click('[data-test-demote-replication] [data-test-replication-action-trigger]');
    await settled();
    // enter confirmation text
    await fillIn('[data-test-confirmation-modal-input="demote"]', 'Performance');
    // Click confirm button
    await click('[data-test-confirm-button="demote"]');
    await settled();
    await click('[data-test-replication-link="details"]');
    await settled();
    assert.dom('[data-test-replication-dashboard]').exists();
    assert.dom('[data-test-selectable-card-container="secondary"]').exists();
    assert.ok(
      find('[data-test-replication-mode-display]').textContent.includes('secondary'),
      'it displays the cluster mode correctly'
    );
  });
});
