import Route from '@ember/routing/route';

export default class OidcScopeRoute extends Route {
  model({ name }) {
    return this.store.findRecord('oidc/scope', name);
  }
}
