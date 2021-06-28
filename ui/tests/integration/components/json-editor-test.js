import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { create } from 'ember-cli-page-object';
import { render } from '@ember/test-helpers';
import hbs from 'htmlbars-inline-precompile';
import jsonEditor from '../../pages/components/json-editor';

const component = create(jsonEditor);

module('Integration | Component | json-editor', function(hooks) {
  setupRenderingTest(hooks);

  const setup = async function(context, title, value, options, showToolbar = true) {
    context.set('value', JSON.stringify(value));
    context.set('options', options);
    context.set('title', title);
    context.set('showToolbar', showToolbar);
    await render(hbs`{{json-editor title=title value=value options=options showToolbar=showToolbar}}`);
  };

  test('it renders', async function(assert) {
    let value = '';
    await setup(this, 'Test title', value, null);
    assert.equal(component.title, 'Test title', 'renders the provided title');
    assert.equal(component.hasToolbar, true, 'renders the toolbar');
    assert.equal(component.hasJSONEditor, true, 'renders the ivy code mirror component');
    assert.equal(component.canEdit, true, 'json editor is in read only mode');
  });

  test('it renders in read only mode', async function(assert) {
    let value = '';
    let options = {
      readOnly: true,
    };
    await setup(this, 'Test title', value, options);
    assert.equal(component.canEdit, false, 'editor is in read only mode');
  });

  test('it renders the editor without toolbar', async function(assert) {
    let value = '';
    await setup(this, 'Test title', value, null, false);
    assert.equal(component.hasToolbar, false, 'toolbar is not rendered');
  });
});
