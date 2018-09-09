import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render, settled } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';

import { create } from 'ember-cli-page-object';
import kvObjectEditor from '../../pages/components/kv-object-editor';

import sinon from 'sinon';
const component = create(kvObjectEditor);

module('Integration | Component | kv object editor', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    component.setContext(this);
  });

  hooks.afterEach(function() {
    component.removeContext();
  });

  test('it renders with no initial value', async function(assert) {
    let spy = sinon.spy();
    this.set('onChange', spy);
    await render(hbs`{{kv-object-editor onChange=onChange}}`);
    assert.equal(component.rows().count, 1, 'renders a single row');
    component.addRow();
    return settled().then(() => {
      assert.equal(component.rows().count, 1, 'will only render row with a blank key');
    });
  });

  test('it calls onChange when the val changes', async function(assert) {
    let spy = sinon.spy();
    this.set('onChange', spy);
    await render(hbs`{{kv-object-editor onChange=onChange}}`);
    component
      .rows(0)
      .kvKey('foo')
      .kvVal('bar');
    settled().then(() => {
      assert.equal(spy.callCount, 2, 'calls onChange each time change is triggered');
      assert.deepEqual(
        spy.lastCall.args[0],
        { foo: 'bar' },
        'calls onChange with the JSON respresentation of the data'
      );
    });
    component.addRow();
    return settled().then(() => {
      assert.equal(component.rows().count, 2, 'adds a row when there is no blank one');
    });
  });

  test('it renders passed data', async function(assert) {
    let metadata = { foo: 'bar', baz: 'bop' };
    this.set('value', metadata);
    await render(hbs`{{kv-object-editor value=value}}`);
    assert.equal(
      component.rows().count,
      Object.keys(metadata).length + 1,
      'renders both rows of the metadata, plus an empty one'
    );
  });

  test('it deletes a row', async function(assert) {
    let spy = sinon.spy();
    this.set('onChange', spy);
    await render(hbs`{{kv-object-editor onChange=onChange}}`);
    component
      .rows(0)
      .kvKey('foo')
      .kvVal('bar');
    component.addRow();
    settled().then(() => {
      assert.equal(component.rows().count, 2);
      assert.equal(spy.callCount, 2, 'calls onChange for editing');
      component.rows(0).deleteRow();
    });

    return settled().then(() => {
      assert.equal(component.rows().count, 1, 'only the blank row left');
      assert.equal(spy.callCount, 3, 'calls onChange deleting row');
      assert.deepEqual(spy.lastCall.args[0], {}, 'last call to onChange is an empty object');
    });
  });

  test('it shows a warning if there are duplicate keys', async function(assert) {
    let metadata = { foo: 'bar', baz: 'bop' };
    this.set('value', metadata);
    await render(hbs`{{kv-object-editor value=value}}`);
    component.rows(0).kvKey('foo');

    return settled().then(() => {
      assert.ok(component.showsDuplicateError, 'duplicate keys are allowed but an error message is shown');
    });
  });
});
