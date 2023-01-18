import { action } from '@ember/object';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import localStorage from 'vault/lib/local-storage';
import Component from '@glimmer/component';

/**
 * @module SecretListHeader
 * SecretListHeader component is breadcrumb, title with icon and menu with tabs component.
 *
 * @example
 * ```js
 * <SecretListHeader
   @model={{this.model}}
   @backendCrumb={{hash
    label=this.model.id
    text=this.model.id
    path="vault.cluster.secrets.backend.list-root"
    model=this.model.id
   }}
  />
 * ```
 * @param {object} model - Model used to pull information about icon and title and backend type for navigation.
 * @param {string} [baseKey] - Provided for navigation on the breadcrumbs.
 * @param {object} [backendCrumb] - Includes label, text, path and model ID.
 * @param {boolean} [isEngine=false] - Changes link type if the component is being used inside an Ember engine.
 */

export default class SecretListHeader extends Component {
  @service router;
  @tracked hideBetaModal;

  // constructor() {
  //   super(...arguments);
  //   console.log('engineType', this.args.model.engineType);
  //   console.log('isEngine', this.args.isEngine);
  // }

  get isKV() {
    return ['kv', 'generic'].includes(this.args.model.engineType);
  }

  get isPki() {
    return this.args.model.engineType === 'pki';
  }

  get shouldHideBetaModal() {
    return localStorage.getItem('hideBetaModal');
  }

  @action
  transitionToNewPki() {
    this.router.transitionTo('vault.cluster.secrets.backend.pki.overview', this.args.model.engineType);
  }

  @action
  toggleHideBetaModal() {
    this.hideBetaModal = !this.hideBetaModal;

    this.hideBetaModal
      ? localStorage.setItem('hideBetaModal', true)
      : localStorage.removeItem('hideBetaModal');
  }
}
