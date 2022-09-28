import { hash } from 'rsvp';
import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default Route.extend({
  store: service(),

  model() {
    return hash({
      cluster: this.modelFor('vault.cluster'),
      seal: this.store.findRecord('capabilities', 'sys/seal'),
    });
  },
});
