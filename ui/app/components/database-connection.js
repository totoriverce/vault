import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';

const LIST_ROOT_ROUTE = 'vault.cluster.secrets.backend.list-root';
const SHOW_ROUTE = 'vault.cluster.secrets.backend.show';

const getErrorMessage = errors => {
  let errorMessage = errors?.join('. ') || 'Something went wrong. Check the Vault logs for more information.';
  if (errorMessage.indexOf('failed to verify') >= 0) {
    errorMessage =
      'There was a verification error for this connection. Check the Vault logs for more information.';
  }
  return errorMessage;
};

export default class DatabaseConnectionEdit extends Component {
  @service store;
  @service router;
  @service flashMessages;
  @service wizard;

  @tracked
  showPasswordField = false; // used for edit mode

  @tracked
  showSaveModal = false; // used for create mode

  constructor() {
    super(...arguments);
    if (this.wizard.featureState === 'details' || this.wizard.featureState === 'connection') {
      this.wizard.transitionFeatureMachine(this.wizard.featureState, 'CONTINUE', 'database');
    }
  }

  rotateCredentials(backend, name) {
    let adapter = this.store.adapterFor('database/connection');
    return adapter.rotateRootCredentials(backend, name);
  }

  transitionToRoute() {
    return this.router.transitionTo(...arguments);
  }

  @action
  updateShowPassword(showForm) {
    this.showPasswordField = showForm;
    if (!showForm) {
      // unset password if hidden
      this.args.model.password = undefined;
    }
  }

  @action
  updatePassword(attr, evt) {
    const value = evt.target.value;
    this.args.model[attr] = value;
  }

  @action
  async handleCreateConnection(evt) {
    evt.preventDefault();
    let secret = this.args.model;
    let secretId = secret.name;
    secret.set('id', secretId);
    secret
      .save()
      .then(() => {
        this.showSaveModal = true;
      })
      .catch(e => {
        const errorMessage = getErrorMessage(e.errors);
        this.flashMessages.danger(errorMessage);
      });
  }

  @action
  continueWithoutRotate(evt) {
    evt.preventDefault();
    const { name } = this.args.model;
    this.transitionToRoute(SHOW_ROUTE, name);
  }

  @action
  continueWithRotate(evt) {
    evt.preventDefault();
    const { backend, name } = this.args.model;
    this.rotateCredentials(backend, name)
      .then(() => {
        this.flashMessages.success(`Successfully rotated root credentials for connection "${name}"`);
        this.transitionToRoute(SHOW_ROUTE, name);
      })
      .catch(e => {
        this.flashMessages.danger(`Error rotating root credentials: ${e.errors}`);
        this.transitionToRoute(SHOW_ROUTE, name);
      });
  }

  @action
  handleUpdateConnection(evt) {
    evt.preventDefault();
    let secret = this.args.model;
    let secretId = secret.name;
    secret
      .save()
      .then(() => {
        this.transitionToRoute(SHOW_ROUTE, secretId);
      })
      .catch(e => {
        const errorMessage = getErrorMessage(e.errors);
        this.flashMessages.danger(errorMessage);
      });
  }

  @action
  delete(evt) {
    evt.preventDefault();
    const secret = this.args.model;
    const backend = secret.backend;
    secret.destroyRecord().then(() => {
      this.transitionToRoute(LIST_ROOT_ROUTE, backend);
    });
  }

  @action
  reset() {
    const { name, backend } = this.args.model;
    let adapter = this.store.adapterFor('database/connection');
    adapter
      .resetConnection(backend, name)
      .then(() => {
        // TODO: Why isn't the confirmAction closing?
        this.flashMessages.success('Successfully reset connection');
      })
      .catch(e => {
        const errorMessage = getErrorMessage(e.errors);
        this.flashMessages.danger(errorMessage);
      });
  }

  @action
  rotate() {
    const { name, backend } = this.args.model;
    this.rotateCredentials(backend, name)
      .then(() => {
        // TODO: Why isn't the confirmAction closing?
        this.flashMessages.success('Successfully rotated credentials');
      })
      .catch(e => {
        const errorMessage = getErrorMessage(e.errors);
        this.flashMessages.danger(errorMessage);
      });
  }
}
