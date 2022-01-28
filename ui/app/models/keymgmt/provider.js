import Model, { attr } from '@ember-data/model';
import { tracked } from '@glimmer/tracking';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';
import { validator, buildValidations } from 'ember-cp-validations';

const CRED_PROPS = {
  azurekeyvault: ['client_id', 'client_secret', 'tenant_id'],
  awskms: ['access_key', 'secret_key', 'session_token', 'endpoint'],
  gcpckms: ['service_account_file'],
};
const OPTIONAL_CRED_PROPS = ['session_token', 'endpoint'];
// since we have dynamic credential attributes based on provider we need a dynamic presence validator
// add validators for all cred props and return true for value if not associated with selected provider
const credValidators = Object.keys(CRED_PROPS).reduce((obj, providerKey) => {
  CRED_PROPS[providerKey].forEach((prop) => {
    if (!OPTIONAL_CRED_PROPS.includes(prop)) {
      obj[`credentials.${prop}`] = validator('presence', {
        presence: true,
        value(model) {
          return model.credentialProps.includes(prop) ? model.credentials[prop] : true;
        },
      });
    }
  });
  return obj;
}, {});
const Validations = buildValidations({
  name: validator('presence', true),
  keyCollection: validator('presence', true),
  ...credValidators,
});
const ValidationsModel = Model.extend(Validations);

export default class KeymgmtProviderModel extends ValidationsModel {
  @attr('string') backend;
  @attr('string', {
    label: 'Provider name',
    subText: 'This is the name of the provider that will be displayed in Vault. This cannot be edited later.',
  })
  name;

  @attr('string', {
    label: 'Type',
    subText: 'Choose the provider type.',
    possibleValues: ['azurekeyvault', 'awskms', 'gcpckms'],
    defaultValue: 'azurekeyvault',
  })
  provider;

  @attr('string', {
    label: 'Key Vault instance name',
    subText: 'The name of a Key Vault instance must be supplied. This cannot be edited later.',
  })
  keyCollection;

  @attr('date') created;

  idPrefix = 'provider/';
  type = 'provider';

  @tracked keys = [];
  @tracked credentials = null; // never returned from API -- set only during create/edit

  get icon() {
    return {
      azurekeyvault: 'azure-color',
      awskms: 'aws-color',
      gcpckms: 'gcp-color',
    }[this.provider];
  }
  get typeName() {
    return {
      azurekeyvault: 'Azure Key Vault',
      awskms: 'AWS Key Management Service',
      gcpckms: 'Google Cloud Key Management Service',
    }[this.provider];
  }
  get showFields() {
    const attrs = expandAttributeMeta(this, ['name', 'created', 'keyCollection']);
    attrs.splice(1, 0, { hasBlock: true, label: 'Type', value: this.typeName, icon: this.icon });
    const l = this.keys.length;
    const value = l ? `${l} ${l > 1 ? 'keys' : 'key'}` : 'None';
    attrs.push({ hasBlock: true, isLink: l, label: 'Keys', value });
    return attrs;
  }
  get credentialProps() {
    return CRED_PROPS[this.provider];
  }
  get credentialFields() {
    const [creds, fields] = this.credentialProps.reduce(
      ([creds, fields], prop) => {
        creds[prop] = null;
        fields.push({ name: `credentials.${prop}`, type: 'string', options: { label: prop } });
        return [creds, fields];
      },
      [{}, []]
    );
    this.credentials = creds;
    return fields;
  }
  get createFields() {
    return expandAttributeMeta(this, ['provider', 'name', 'keyCollection']);
  }

  async fetchKeys(page) {
    try {
      this.keys = await this.store.lazyPaginatedQuery('keymgmt/key', {
        backend: 'keymgmt',
        provider: this.name,
        responsePath: 'data.keys',
        page,
      });
    } catch (error) {
      this.keys = [];
      if (error.httpStatus !== 404) {
        throw error;
      }
    }
  }
}
