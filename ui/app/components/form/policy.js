import Component from '@glimmer/component';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';
import { task } from 'ember-concurrency';
import trimRight from 'vault/utils/trim-right';
import { tracked } from '@glimmer/tracking';

/**
 * @module Form::Policy
 * Form::Policy components are the forms used to display the create and edit forms for all types of policies. This is only the form, not the outlying layout, and expects that the form model is passed from the parent.
 *
 * @example
 *  <Form::Policy
 *    @model={{this.model}}
 *    @onSave={{transition-to "vault.cluster.policy.show" this.model.policyType this.model.name}}
 *    @onCancel={{transition-to "vault.cluster.policies.index"}}
 *  />
 * ```
 * @callback onCancel - callback triggered when cancel button is clicked
 * @callback onSave - callback triggered when save button is clicked. Passes saved model
 * @param {object} model - ember data model from createRecord
 */

export default class FormPolicyComponent extends Component {
  @service flashMessages;
  @service wizard;

  @tracked errorBanner = '';

  @task
  *save(event) {
    event.preventDefault();
    try {
      yield this.args.model.save();
      // parent is in charge of flash messages, closing modals, transitions, etc
      this.args.onSave(this.args.model);
    } catch (error) {
      const message = error.errors ? error.errors.join('. ') : error.message;
      this.errorBanner = message;
    }
  }

  @action
  setModelName({ target }) {
    this.args.model.name = target.value.toLowerCase();
  }

  @action
  setPolicyFromFile(index, fileInfo) {
    const { value, fileName } = fileInfo;
    this.args.model.policy = value;
    if (!this.args.model.name) {
      const trimmedFileName = trimRight(fileName, ['.json', '.txt', '.hcl', '.policy']);
      this.args.model.name = trimmedFileName.toLowerCase();
    }
    this.showFileUpload = false;
  }

  @action
  cancel() {
    this.cleanup();
    this.args.onCancel();
  }

  cleanup() {
    const method = this.args.model.isNew ? 'unloadRecord' : 'rollbackAttributes';
    this.args.model[method]();
  }
}
