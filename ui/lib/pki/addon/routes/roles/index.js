import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';

export default class RolesIndexRoute extends Route {
  @service store;
  @service secretMountPath;
  @service pathHelp;

  model() {
    // the pathHelp service is needed for adding openAPI to the model
    this.pathHelp.getNewModel('pki/pki-role-engine', 'pki');

    return this.store
      .query('pki/pki-role-engine', { backend: this.secretMountPath.currentPath })
      .catch((err) => {
        if (err.httpStatus === 404) {
          return [];
        } else {
          throw err;
        }
      });
  }
}
