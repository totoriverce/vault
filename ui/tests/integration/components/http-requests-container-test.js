import { module, test } from 'qunit';
import { setupRenderingTest } from 'ember-qunit';
import { render } from '@ember/test-helpers';
import { create } from 'ember-cli-page-object';
import httpRequestsContainer from 'vault/tests/pages/components/http-requests-container';
import hbs from 'htmlbars-inline-precompile';

const COUNTERS = [
  { start_time: '2018-04-01T00:00:00Z', total: 5500 },
  { start_time: '2019-05-01T00:00:00Z', total: 4500 },
  { start_time: '2019-06-01T00:00:00Z', total: 5000 },
];

module('Integration | Component | http-requests-container', function(hooks) {
  setupRenderingTest(hooks);

  hooks.beforeEach(function() {
    this.set('counters', COUNTERS);
  });

  test('it renders', async function(assert) {
    await render(hbs`<HttpRequestsContainer @counters={{counters}}/>`);

    assert.dom('.http-requests-container').exists();
    assert.dom('.http-requests-dropdown').exists();
    assert.dom('.http-requests-bar-chart-container').exists();
    assert.dom('.http-requests-table').exists();
  });

  test('it does not render a bar chart for less than one month of data', async function(assert) {
    const one_month_counter = [
      {
        start_time: '2019-05-01T00:00:00Z',
        total: 50,
      },
    ];
    this.set('one_month_counter', one_month_counter);

    await render(hbs`<HttpRequestsContainer @counters={{one_month_counter}}/>`);

    assert.dom('.http-requests-table').exists();
    assert.dom('.http-requests-bar-chart-container').doesNotExist();
  });

  test('it filters the data according to the dropdown', async function(assert) {
    await render(hbs`<HttpRequestsContainer @counters={{counters}}/>`);
    const component = create(httpRequestsContainer);
    await component.select('2018');

    assert.equal(
      this.element.querySelectorAll('.shadow-bars> .bar').length,
      1,
      'filters the bar chart to the selected year'
    );

    assert.equal(
      this.element.querySelectorAll('.start-time').length,
      1,
      'filters the table to the selected year'
    );
  });
});
