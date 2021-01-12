import { inject as service } from '@ember/service';
import { alias, gt } from '@ember/object/computed';
import Component from '@ember/component';
import { computed } from '@ember/object';
import keyUtils from 'vault/lib/key-utils';
import pathToTree from 'vault/lib/path-to-tree';
import { task, timeout } from 'ember-concurrency';

const { ancestorKeysForKey } = keyUtils;
const DOT_REPLACEMENT = '☃';
const ANIMATION_DURATION = 250;

export default Component.extend({
  tagName: '',
  namespaceService: service('namespace'),
  auth: service(),
  store: service(),
  namespace: null,
  listCapability: null,
  canList: false,

  init() {
    this._super(...arguments);
    this.namespaceService?.findNamespacesForUser.perform();
  },

  didReceiveAttrs() {
    this._super(...arguments);

    let ns = this.namespace;
    let oldNS = this.oldNamespace;
    if (!oldNS || ns !== oldNS) {
      this.setForAnimation.perform();
      this.fetchListCapability.perform();
    }
    this.set('oldNamespace', ns);
  },

  fetchListCapability: task(function*() {
    try {
      let capability = yield this.store.findRecord('capabilities', 'sys/namespaces/');
      this.set('listCapability', capability);
      this.set('canList', true);
    } catch (e) {
      // If error out on findRecord call it's because you don't have permissions
      // and therefor don't have permission to manage namespaces
      this.set('canList', false);
    }
  }),
  setForAnimation: task(function*() {
    let leaves = this.menuLeaves;
    let lastLeaves = this.lastMenuLeaves;
    if (!lastLeaves) {
      this.set('lastMenuLeaves', leaves);
      yield timeout(0);
      return;
    }
    let isAdding = leaves.length > lastLeaves.length;
    let changedLeaf = (isAdding ? leaves : lastLeaves).get('lastObject');
    this.set('isAdding', isAdding);
    this.set('changedLeaf', changedLeaf);

    // if we're adding we want to render immediately an animate it in
    // if we're not adding, we need time to move the item out before
    // a rerender removes it
    if (isAdding) {
      this.set('lastMenuLeaves', leaves);
      yield timeout(0);
      return;
    }
    yield timeout(ANIMATION_DURATION);
    this.set('lastMenuLeaves', leaves);
  }).drop(),

  isAnimating: alias('setForAnimation.isRunning'),

  namespacePath: alias('namespaceService.path'),

  // this is an array of namespace paths that the current user
  // has access to
  accessibleNamespaces: alias('namespaceService.accessibleNamespaces'),
  inRootNamespace: alias('namespaceService.inRootNamespace'),

  namespaceTree: computed('accessibleNamespaces', function() {
    let nsList = this.accessibleNamespaces;

    if (!nsList) {
      return [];
    }
    return pathToTree(nsList);
  }),

  maybeAddRoot(leaves) {
    let userRoot = this.auth.authData.userRootNamespace;
    if (userRoot === '') {
      leaves.unshift('');
    }

    return leaves.uniq();
  },

  pathToLeaf(path) {
    // dots are allowed in namespace paths
    // so we need to preserve them, and replace slashes with dots
    // in order to use Ember's get function on the namespace tree
    // to pull out the correct level
    return (
      path
        // trim trailing slash
        .replace(/\/$/, '')
        // replace dots with snowman
        .replace(/\.+/g, DOT_REPLACEMENT)
        // replace slash with dots
        .replace(/\/+/g, '.')
    );
  },

  // an array that keeps track of what additional panels to render
  // on the menu stack
  // if you're in  'foo/bar/baz',
  // this array will be: ['foo', 'foo.bar', 'foo.bar.baz']
  // the template then iterates over this, and does  Ember.get(namespaceTree, leaf)
  // to render the nodes of each leaf

  // gets set as  'lastMenuLeaves' in the ember concurrency task above
  menuLeaves: computed('namespacePath', 'namespaceTree', 'pathToLeaf', function() {
    let ns = this.namespacePath;
    ns = ns.replace(/^\//, '');
    let leaves = ancestorKeysForKey(ns);
    leaves.push(ns);
    leaves = this.maybeAddRoot(leaves);

    leaves = leaves.map(this.pathToLeaf);
    return leaves;
  }),

  // the nodes at the root of the namespace tree
  // these will get rendered as the bottom layer
  rootLeaves: computed('namespaceTree', function() {
    let tree = this.namespaceTree;
    let leaves = Object.keys(tree);
    return leaves;
  }),

  currentLeaf: alias('lastMenuLeaves.lastObject'),
  canAccessMultipleNamespaces: gt('accessibleNamespaces.length', 1),
  isUserRootNamespace: computed('auth.authData.userRootNamespace', 'namespacePath', function() {
    return this.auth.authData.userRootNamespace === this.namespacePath;
  }),

  namespaceDisplay: computed('namespacePath', 'accessibleNamespaces', 'accessibleNamespaces.[]', function() {
    let namespace = this.namespacePath;
    if (namespace === '') {
      return '';
    }
    let parts = namespace.split('/');
    return parts[parts.length - 1];
  }),

  actions: {
    refreshNamespaceList() {
      this.namespaceService.findNamespacesForUser.perform();
    },
  },
});
