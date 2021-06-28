import { belongsTo, attr } from '@ember-data/model';
import { alias } from '@ember/object/computed';
import { computed } from '@ember/object';
import IdentityModel from './_base';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import identityCapabilities from 'vault/macros/identity-capabilities';

export default IdentityModel.extend({
  formFields: computed('type', function() {
    let fields = ['name', 'type', 'policies', 'metadata'];
    if (this.type === 'internal') {
      return fields.concat(['memberGroupIds', 'memberEntityIds']);
    }
    return fields;
  }),
  name: attr('string'),
  type: attr('string', {
    defaultValue: 'internal',
    possibleValues: ['internal', 'external'],
  }),
  creationTime: attr('string', {
    readOnly: true,
  }),
  lastUpdateTime: attr('string', {
    readOnly: true,
  }),
  numMemberEntities: attr('number', {
    readOnly: true,
  }),
  numParentGroups: attr('number', {
    readOnly: true,
  }),
  metadata: attr('object', {
    editType: 'kv',
  }),
  policies: attr({
    label: 'Policies',
    editType: 'searchSelect',
    fallbackComponent: 'string-list',
    models: ['policy/acl', 'policy/rgp'],
  }),
  memberGroupIds: attr({
    label: 'Member Group IDs',
    editType: 'searchSelect',
    fallbackComponent: 'string-list',
    models: ['identity/group'],
  }),
  parentGroupIds: attr({
    label: 'Parent Group IDs',
    editType: 'searchSelect',
    fallbackComponent: 'string-list',
    models: ['identity/group'],
  }),
  memberEntityIds: attr({
    label: 'Member Entity IDs',
    editType: 'searchSelect',
    fallbackComponent: 'string-list',
    models: ['identity/entity'],
  }),
  hasMembers: computed(
    'memberEntityIds',
    'memberEntityIds.[]',
    'memberGroupIds',
    'memberGroupIds.[]',
    function() {
      let { memberEntityIds, memberGroupIds } = this;
      let numEntities = (memberEntityIds && memberEntityIds.length) || 0;
      let numGroups = (memberGroupIds && memberGroupIds.length) || 0;
      return numEntities + numGroups > 0;
    }
  ),

  alias: belongsTo('identity/group-alias', { async: false, readOnly: true }),
  updatePath: identityCapabilities(),
  canDelete: alias('updatePath.canDelete'),
  canEdit: alias('updatePath.canUpdate'),

  aliasPath: lazyCapabilities(apiPath`identity/group-alias`),
  canAddAlias: computed('aliasPath.canCreate', 'type', 'alias', function() {
    let type = this.type;
    let alias = this.alias;
    // internal groups can't have aliases, and external groups can only have one
    if (type === 'internal' || alias) {
      return false;
    }
    return this.aliasPath.canCreate;
  }),
});
