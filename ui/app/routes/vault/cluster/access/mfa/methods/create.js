import Route from '@ember/routing/route';

export default class MfaLoginEnforcementCreateRoute extends Route {
  setupController(controller) {
    super.setupController(...arguments);
    // if route was refreshed after type select recreate method model
    const { type } = controller;
    if (type) {
      controller.set('method', this.store.createRecord('mfa-method', { type }));
    }
  }
  resetController(controller, isExiting) {
    if (isExiting) {
      // reset type query param when user saves or cancels
      // this will not trigger when refreshing the page which preserves intended functionality
      controller.set('type', null);
    }
  }
}
