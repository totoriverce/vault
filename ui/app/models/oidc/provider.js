import Model, { attr } from '@ember-data/model';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';
import { withModelValidations } from 'vault/decorators/model-validations';

const validations = {
  name: [{ type: 'presence', message: 'Name is required.' }],
};

@withModelValidations(validations)
export default class OidcProviderModel extends Model {
  @attr('string', { editDisabled: true }) name;
  @attr('string', {
    label: 'Issuer URL',
    subText: 'The issuer URL parameter is necessary for the validation of ID tokens by applications.',
    docLink: '/api-docs/secret/identity/oidc-provider#create-or-update-a-provider',
  })
  issuer;
  @attr('array', {
    label: 'Supported scopes',
    subText: 'Scopes define information about a user and the OIDC service. Optional.',
    editType: 'searchSelect',
    models: ['oidc/scope'],
  })
  scopesSupported;
  @attr('array', { label: 'Allowed applications' }) allowedClientIds; // no editType because does not use form-field component

  // TODO refactor when field-to-attrs is refactored as decorator
  _attributeMeta = null; // cache initial result of expandAttributeMeta in getter and return
  get formFields() {
    if (!this._attributeMeta) {
      this._attributeMeta = expandAttributeMeta(this, ['name', 'issuer', 'scopesSupported']);
    }
    return this._attributeMeta;
  }
  @lazyCapabilities(apiPath`identity/oidc/provider/${'name'}`, 'name') providerPath;
  @lazyCapabilities(apiPath`identity/oidc/provider`) providersPath;
  get canCreate() {
    return this.providerPath.get('canCreate');
  }
  get canRead() {
    return this.providerPath.get('canRead');
  }
  get canEdit() {
    return this.providerPath.get('canUpdate');
  }
  get canDelete() {
    return this.providerPath.get('canDelete');
  }
  get canList() {
    return this.providersPath.get('canList');
  }

  @lazyCapabilities(apiPath`identity/oidc/client`) clientsPath;
  get canListClients() {
    return this.clientsPath.get('canList');
  }
  @lazyCapabilities(apiPath`identity/oidc/scope`) scopesPath;
  get canListScopes() {
    return this.scopesPath.get('canList');
  }
}
