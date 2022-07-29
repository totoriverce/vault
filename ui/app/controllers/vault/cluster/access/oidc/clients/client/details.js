import Controller from '@ember/controller';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';

export default class OidcClientDetailsController extends Controller {
  @service router;
  @service flashMessages;

  queryParams = ['tab'];
  tab = 'details';

  @action
  async delete() {
    try {
      await this.model.destroyRecord();
      this.flashMessages.success('Application deleted successfully');
      this.router.transitionTo('vault.cluster.access.oidc.clients');
    } catch (error) {
      const message = error.errors ? error.errors.join('. ') : error.message;
      this.flashMessages.danger(message);
    }
  }
}
