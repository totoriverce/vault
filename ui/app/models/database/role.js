import Model, { attr } from '@ember-data/model';
import { computed } from '@ember/object';
import { alias } from '@ember/object/computed';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import { expandAttributeMeta } from 'vault/utils/field-to-attrs';

export default Model.extend({
  idPrefix: 'role/',
  backend: attr('string', { readOnly: true }),
  name: attr('string', {
    label: 'Role name',
  }),
  database: attr('array', {
    label: '',
    editType: 'searchSelect',
    fallbackComponent: 'string-list',
    models: ['database/connection'],
    selectLimit: 1,
    onlyAllowExisting: true,
    subLabel: 'Database name',
    subText: 'The database for which credentials will be generated.',
  }),
  type: attr('string', {
    label: 'Type of role',
    noDefault: true,
    possibleValues: ['static', 'dynamic'],
  }),
  ttl: attr({
    editType: 'ttl',
    defaultValue: '1h',
    label: 'Generated credentials’s Time-to-Live (TTL)',
    subText: 'Vault will use the engine default of 1 hour',
    defaultShown: 'Engine default',
  }),
  max_ttl: attr({
    editType: 'ttl',
    defaultValue: '24h',
    label: 'Generated credentials’s maximum Time-to-Live (Max TTL)',
    subText: 'Vault will use the engine default of 24 hours',
    defaultShown: 'Engine default',
  }),
  username: attr('string', {
    subText: 'The database username that this Vault role corresponds to.',
  }),
  rotation_period: attr({
    editType: 'ttl',
    defaultValue: '24h',
    subText:
      'Specifies the amount of time Vault should wait before rotating the password. The minimum is 5 seconds. Default is 24 hours.',
  }),
  creation_statements: attr('array', {
    editType: 'stringArray',
  }),
  revocation_statements: attr('array', {
    editType: 'stringArray',
    defaultShown: 'Default',
  }),
  rotation_statements: attr('array', {
    editType: 'stringArray',
    defaultShown: 'Default',
  }),
  rollback_statements: attr('array', {
    editType: 'stringArray',
    defaultShown: 'Default',
  }),
  renew_statements: attr('array', {
    editType: 'stringArray',
    defaultShown: 'Default',
  }),
  creation_statement: attr('string', {
    editType: 'json',
    allowReset: true,
    theme: 'hashi short',
    defaultShown: 'Default',
  }),
  revocation_statement: attr('string', {
    editType: 'json',
    allowReset: true,
    theme: 'hashi short',
    defaultShown: 'Default',
  }),

  /* FIELD ATTRIBUTES */
  get fieldAttrs() {
    // Main fields on edit/create form
    let fields = ['name', 'database', 'type'];
    return expandAttributeMeta(this, fields);
  },

  get showFields() {
    let fields = ['name', 'database', 'type'];
    if (this.type === 'dynamic') {
      fields = fields.concat(['ttl', 'max_ttl', 'creation_statements', 'revocation_statements']);
    } else {
      fields = fields.concat(['username', 'rotation_period']);
    }
    return expandAttributeMeta(this, fields);
  },

  roleSettingAttrs: computed(function() {
    // logic for which get displayed is on DatabaseRoleSettingForm
    let allRoleSettingFields = [
      'ttl',
      'max_ttl',
      'username',
      'rotation_period',
      'creation_statements',
      'creation_statement', // only for MongoDB (styling difference)
      'revocation_statements',
      'revocation_statement', // only for MongoDB (styling difference)
      'rotation_statements',
      'rollback_statements',
      'renew_statements',
    ];
    return expandAttributeMeta(this, allRoleSettingFields);
  }),

  /* CAPABILITIES */
  // only used for secretPath
  path: attr('string', { readOnly: true }),

  secretPath: lazyCapabilities(apiPath`${'backend'}/${'path'}/${'id'}`, 'backend', 'path', 'id'),
  canEditRole: alias('secretPath.canUpdate'),
  canDelete: alias('secretPath.canDelete'),
  dynamicPath: lazyCapabilities(apiPath`${'backend'}/roles/+`, 'backend'),
  canCreateDynamic: alias('dynamicPath.canCreate'),
  staticPath: lazyCapabilities(apiPath`${'backend'}/static-roles/+`, 'backend'),
  canCreateStatic: alias('staticPath.canCreate'),
  credentialPath: lazyCapabilities(apiPath`${'backend'}/creds/${'id'}`, 'backend', 'id'),
  canGenerateCredentials: alias('credentialPath.canRead'),
  databasePath: lazyCapabilities(apiPath`${'backend'}/config/${'database[0]'}`, 'backend', 'database'),
  canUpdateDb: alias('databasePath.canUpdate'),
});
