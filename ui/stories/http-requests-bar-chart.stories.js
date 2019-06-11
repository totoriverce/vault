/* eslint-disable import/extensions */
import hbs from 'htmlbars-inline-precompile';
import { storiesOf } from '@storybook/ember';
import { withKnobs, object } from '@storybook/addon-knobs';
import notes from './http-requests-bar-chart.md';

const COUNTERS = [
  { start_time: '2019-04-01T00:00:00.000Z', total: 5500 },
  { start_time: '2019-05-01T00:00:00.000Z', total: 4500 },
  { start_time: '2019-06-01T00:00:00.000Z', total: 5000 },
];

storiesOf('HttpRequests/BarChart/', module)
  .addParameters({ options: { showPanel: true } })
  .addDecorator(
    withKnobs()
  )
  .add(`HttpRequestsBarChart`, () => ({
    template: hbs`
        <h5 class="title is-5">Http Requests Bar Chart</h5>
        <HttpRequestsBarChart @counters={{counters}}/>
    `,
    context: {
      counters: object('counters', COUNTERS)
    },
  }),
  {notes}
);
