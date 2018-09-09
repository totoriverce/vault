import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { create } from 'ember-cli-page-object';
import hbs from 'htmlbars-inline-precompile';
import maskedInput from 'vault/tests/pages/components/masked-input';

const component = create(maskedInput);

module('Integration | Component | masked input', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    component.setContext(this);
  });

  hooks.afterEach(function() {
    component.removeContext();
  });

  const hasClass = (classString = '', classToFind) => {
    return classString.split(' ').contains(classToFind);
  };

  test('it renders', async function(assert) {
    await render(hbs`{{masked-input}}`);

    assert.ok(hasClass(component.wrapperClass, 'masked'));
  });

  test('it renders a textarea', async function(assert) {
    await render(hbs`{{masked-input}}`);

    assert.ok(component.textareaIsPresent);
  });

  test('it does not render a textarea when displayOnly is true', async function(assert) {
    await render(hbs`{{masked-input displayOnly=true}}`);

    assert.notOk(component.textareaIsPresent);
  });

  test('it unmasks text on focus', async function(assert) {
    this.set('value', 'value');
    await render(hbs`{{masked-input value=value}}`);

    assert.ok(hasClass(component.wrapperClass, 'masked'));

    component.focus();
    assert.notOk(hasClass(component.wrapperClass, 'masked'));
  });

  test('it remasks text on blur', async function(assert) {
    this.set('value', 'value');
    await render(hbs`{{masked-input value=value}}`);

    assert.ok(hasClass(component.wrapperClass, 'masked'));

    component.focus();
    component.blur();

    assert.ok(hasClass(component.wrapperClass, 'masked'));
  });

  test('it unmasks text when button is clicked', async function(assert) {
    this.set('value', 'value');
    await render(hbs`{{masked-input value=value}}`);

    assert.ok(hasClass(component.wrapperClass, 'masked'));

    component.toggleMasked();

    assert.notOk(hasClass(component.wrapperClass, 'masked'));
  });

  test('it remasks text when button is clicked', async function(assert) {
    this.set('value', 'value');
    await render(hbs`{{masked-input value=value}}`);

    component.toggleMasked();
    component.toggleMasked();

    assert.ok(hasClass(component.wrapperClass, 'masked'));
  });
});
