import Ember from 'ember';
import flat from 'flat';
import deepmerge from 'deepmerge';
import keyUtils from 'vault/lib/key-utils';
import { task, timeout } from 'ember-concurrency';

const { ancestorKeysForKey } = keyUtils;
const { unflatten } = flat;
const { Component, computed, inject } = Ember;
const DOT_REPLACEMENT = '☃';
const ANIMATION_DURATION = 250;

export default Component.extend({
  tagName: '',
  namespaceService: inject.service('namespace'),
  auth: inject.service(),
  namespace: null,

  init() {
    this._super(...arguments);
    this.get('namespaceService.findNamespacesForUser').perform();
  },

  didReceiveAttrs() {
    this._super(...arguments);

    let ns = this.get('namespace');
    let oldNS = this.get('oldNamespace');
    if (!oldNS || ns !== oldNS) {
      this.get('setForAnimation').perform();
    }
    this.set('oldNamespace', ns);
  },

  setForAnimation: task(function*() {
    let leaves = this.get('menuLeaves');
    let lastLeaves = this.get('lastMenuLeaves');
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
  }).restartable(),

  isAnimating: computed.alias('setForAnimation.isRunning'),

  namespacePath: computed.alias('namespaceService.path'),

  // this is an array of namespace paths that the current user
  // has access to
  accessibleNamespaces: computed.alias('namespaceService.accessibleNamespaces'),

  namespaceTree: computed('accessibleNamespaces', function() {
    let nsList = this.maybeAddRoot(this.get('accessibleNamespaces'));

    if (!nsList) {
      return [];
    }
    // first sort the list by length, then alphanumeric
    nsList = nsList.slice(0).sort((a, b) => b.length - a.length || b.localeCompare(a));
    // then reduce to an array
    // and we remove all of the items that have a string
    // that starts with the same prefix from the list
    // so if we have "foo/bar/baz", both "foo" and "foo/bar"
    // won't be included in the list
    let nsTree = nsList.reduce((accumulator, ns) => {
      let prefixInList = accumulator.some(nsPath => nsPath.startsWith(ns));
      if (!prefixInList) {
        accumulator.push(ns);
      }
      return accumulator;
    }, []);

    // after the reduction we're left with an array that contains
    // strings that represent the longest branches
    // we'll replace the dots in the paths, then expand the path
    // to a nested object that we can then query with Ember.get
    return deepmerge.all(
      nsTree.map(ns => {
        ns = ns.replace(/\.+/g, DOT_REPLACEMENT);
        return unflatten({ [ns]: null }, { delimiter: '/' });
      })
    );
  }),

  maybeAddRoot(leaves) {
    let userRoot = this.get('auth.authData.userRootNamespace');
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
  menuLeaves: computed('namespacePath', 'namespaceTree', function() {
    let ns = this.get('namespacePath');
    let leaves = ancestorKeysForKey(ns) || [];
    leaves.push(ns);
    leaves = this.maybeAddRoot(leaves);

    leaves = leaves.map(this.pathToLeaf);
    return leaves;
  }),

  // the nodes at the root of the namespace tree
  // these will get rendered as the bottom layer
  rootLeaves: computed('namespaceTree', function() {
    let tree = this.get('namespaceTree');
    let leaves = Object.keys(tree);
    return leaves;
  }),

  currentLeaf: computed.alias('lastMenuLeaves.lastObject'),
  canAccessMultipleNamespaces: computed.gt('accessibleNamespaces.length', 1),
  isUserRootNamespace: computed('auth.authData.userRootNamespace', 'namespacePath', function() {
    return this.get('auth.authData.userRootNamespace') === this.get('namespacePath');
  }),

  namespaceDisplay: computed('namespacePath', 'accessibleNamespaces', 'accessibleNamespaces.[]', function() {
    let namespace = this.get('namespacePath');
    if (namespace === '') {
      return '';
    }
    let parts = namespace.split('/');
    return parts[parts.length - 1];
  }),
});
