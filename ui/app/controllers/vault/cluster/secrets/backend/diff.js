/* eslint-disable no-undef */
import Controller from '@ember/controller';
import BackendCrumbMixin from 'vault/mixins/backend-crumb';
import { inject as service } from '@ember/service';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

export default class DiffController extends Controller.extend(BackendCrumbMixin) {
  @service store;
  @tracked
  leftSideVersionDataSelected = null;
  @tracked
  rightSideVersionDataSelected = null;
  @tracked
  leftSideVersionSelected = null;
  @tracked
  rightSideVersionSelected = null;
  @tracked
  statesMatch = false;

  adapter = this.store.adapterFor('secret-v2-version');

  get leftSideDataInit() {
    // return secretData from hitting the get secret endpoint
    let string = `["${this.model.engineId}", "${this.model.id}", "${this.model.currentVersion}"]`;
    return this.adapter
      .querySecretDataByVersion(string)
      .then(response => response.data) // using ember promise helpers to await in the hbs file
      .catch(() => null);
  }

  get getRightSideVersionInit() {
    // initial value of right side version is one less than the current version
    return this.model.currentVersion === 1 ? 0 : this.model.currentVersion - 1;
  }

  get rightSideDataInit() {
    // return secretData from hitting the get secret endpoint
    let string = `["${this.model.engineId}", "${this.model.id}", "${this.getRightSideVersionInit}"]`;
    return this.adapter
      .querySecretDataByVersion(string)
      .then(response => response.data) // using ember promise helpers to await in the hbs file\
      .catch(() => null);
  }

  async createVisualDiff() {
    let diffpatcher = jsondiffpatch.create({});
    let leftSideVersionData = this.leftSideVersionDataSelected || (await this.leftSideDataInit);
    let rightSideVersionData = this.rightSideVersionDataSelected || (await this.rightSideDataInit);
    let delta = diffpatcher.diff(leftSideVersionData, rightSideVersionData);

    if (delta === undefined) {
      this.statesMatch = true;
      return JSON.stringify(leftSideVersionData, undefined, 2); // value, replacer (all properties included), space (white space and indentation, line break, etc.)
    } else {
      this.statesMatch = false;
    }

    return jsondiffpatch.formatters.html.format(delta, this.leftSideVersionData);
  }

  get visualDiff() {
    return this.createVisualDiff();
  }

  // ARG TODO I believe I can remove this but double check
  @action
  refreshModel() {
    this.send('refreshModel');
  }
  @action
  async selectLeftSideVersion(selectedVersion, actions) {
    let string = `["${this.model.engineId}", "${this.model.id}", "${selectedVersion}"]`;
    let secretData = await this.adapter.querySecretDataByVersion(string);
    this.leftSideVersionDataSelected = secretData.data;
    this.leftSideVersionSelected = selectedVersion;
    // close dropdown menu.
    actions.close();
  }
  @action
  async selectRightSideVersion(selectedVersion, actions) {
    let string = `["${this.model.engineId}", "${this.model.id}", "${selectedVersion}"]`;
    let secretData = await this.adapter.querySecretDataByVersion(string);
    this.rightSideVersionDataSelected = secretData.data;
    this.rightSideVersionSelected = selectedVersion;
    // close dropdown menu.
    actions.close();
  }
}
