import { hash } from 'rsvp';
import FetchConfigRoute from './fetch-config';

export default class KubernetesOverviewRoute extends FetchConfigRoute {
  async model() {
    const backend = this.secretMountPath.get();

    return hash({
      config: this.configModel,
      backend: this.modelFor('application'),
      roles: this.store.query('kubernetes/role', { backend }),
    });
  }
}
