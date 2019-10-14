import Route from '@ember/routing/route';
import ClusterRoute from 'vault/mixins/cluster-route';

export default Route.extend(ClusterRoute, {
  model() {
    return this.store.findAll('server');
  },

  actions: {
    doRefresh() {
      this.refresh();
    },
  },
});
