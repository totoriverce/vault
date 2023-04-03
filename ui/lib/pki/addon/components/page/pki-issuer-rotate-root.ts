import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import { tracked } from '@glimmer/tracking';
import { action } from '@ember/object';
import { waitFor } from '@ember/test-waiters';
import { task } from 'ember-concurrency';
import errorMessage from 'vault/utils/error-message';
// TYPES
import Store from '@ember-data/store';
import Router from '@ember/routing/router';
import FlashMessageService from 'ember-cli-flash/services/flash-messages';
import SecretMountPath from 'vault/services/secret-mount-path';
import PkiIssuerModel from 'vault/models/pki/issuer';
import PkiActionModel from 'vault/vault/models/pki/action';
import { Breadcrumb } from 'vault/vault/app-types';
import { parsedParameters } from 'vault/utils/parse-pki-cert-oids';

interface Args {
  oldRoot: PkiIssuerModel;
  newRootModel: PkiActionModel;
  breadcrumbs: Breadcrumb;
  parsingErrors: string;
}

const RADIO_BUTTON_KEY = {
  oldSettings: 'use-old-settings',
  customizeNew: 'customize',
};

export default class PagePkiIssuerRotateRootComponent extends Component<Args> {
  @service declare readonly store: Store;
  @service declare readonly router: Router;
  @service declare readonly flashMessages: FlashMessageService;
  @service declare readonly secretMountPath: SecretMountPath;

  @tracked displayedForm = RADIO_BUTTON_KEY.oldSettings;
  @tracked showOldSettings = false;
  // form alerts below are only for "use old settings" option
  // validations/errors for "customize new root" are handled by <PkiGenerateRoot> component
  @tracked alertBanner = '';
  @tracked invalidFormAlert = '';
  @tracked modelValidations = null;

  get bannerType() {
    if (this.args.parsingErrors && !this.invalidFormAlert) {
      return {
        title: 'Not all of the certificate values could be parsed and transfered to new root',
        type: 'warning',
      };
    }
    return { type: 'danger' };
  }

  get generateOptions() {
    return [
      {
        key: RADIO_BUTTON_KEY.oldSettings,
        icon: 'certificate',
        label: 'Use old root settings',
        description: `Provide only a new common name and issuer name, using the old root’s settings. Selecting this option generates a root with Vault-internal key material.`,
      },
      {
        key: RADIO_BUTTON_KEY.customizeNew,
        icon: 'award',
        label: 'Customize new root certificate',
        description:
          'Generates a new self-signed CA certificate and private key. This generated root will sign its own CRL.',
      },
    ];
  }

  // for displaying old root details, and generated root details
  get displayFields() {
    const addKeyFields = ['privateKey', 'privateKeyType'];
    const defaultFields = [
      'certificate',
      'caChain',
      'issuerId',
      'issuerName',
      'issuingCa',
      'keyName',
      'keyId',
      'serialNumber',
      ...parsedParameters,
    ];
    return this.args.newRootModel.id ? [...defaultFields, ...addKeyFields] : defaultFields;
  }

  checkFormValidity() {
    if (this.args.newRootModel.validate) {
      const { isValid, state, invalidFormMessage } = this.args.newRootModel.validate();
      this.modelValidations = state;
      this.invalidFormAlert = invalidFormMessage;
      return isValid;
    }
    return true;
  }

  @task
  @waitFor
  *save(event: Event) {
    event.preventDefault();
    const continueSave = this.checkFormValidity();
    if (!continueSave) return;
    try {
      yield this.args.newRootModel.save({ adapterOptions: { actionType: 'rotate-root' } });
      this.flashMessages.success('Successfully generated root.');
    } catch (e) {
      this.alertBanner = errorMessage(e);
      this.invalidFormAlert = 'There was a problem generating root.';
    }
  }

  @action
  async fetchDataForDownload(format: string) {
    const endpoint = `/v1/${this.secretMountPath.currentPath}/issuer/${this.args.newRootModel.issuerId}/${format}`;
    const adapter = this.store.adapterFor('application');
    try {
      return adapter
        .rawRequest(endpoint, 'GET', { unauthenticated: true })
        .then(function (response: Response) {
          if (format === 'der') {
            return response.blob();
          }
          return response.text();
        });
    } catch (e) {
      return null;
    }
  }
}
