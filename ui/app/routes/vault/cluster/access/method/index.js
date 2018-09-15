import Route from '@ember/routing/route';

export default Route.extend({
  beforeModel() {
    return this.transitionTo('vault.cluster.access.method.section', 'configuration');
  },
});
