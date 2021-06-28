import { inject as service } from '@ember/service';
import Component from '@glimmer/component';

/**
 * @module ClusterInfo
 *
 * @example
 * ```js
 * <ClusterInfo @cluster={{cluster}} @onLinkClick={{action}} />
 * ```
 *
 * @param {object} cluster - details of the current cluster, passed from the parent.
 * @param {Function} onLinkClick - parent action which determines the behavior on link click
 */
export default class ClusterInfoComponent extends Component {
  @service auth;
  @service store;
  @service version;

  get activeCluster() {
    return this.store.peekRecord('cluster', this.auth.activeCluster);
  }

  transitionToRoute() {
    this.router.transitionTo(...arguments);
  }
}
