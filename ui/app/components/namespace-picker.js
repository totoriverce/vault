import Ember from 'ember';
import flat from 'flat';
import deepmerge from 'deepmerge';
import keyUtils from 'vault/lib/key-utils';
import { task, timeout } from 'ember-concurrency';

const { ancestorKeysForKey } = keyUtils;
const { unflatten } = flat;
const { Component, computed, inject } = Ember;
const DOT_REPLACEMENT = '☃';
const ANIMATION_DURATION = 300;

export default Component.extend({
  namespaceService: inject.service('namespace'),
  auth: inject.service(),

  init() {
    this._super(...arguments);
    this.get('namespaceService.findNamespacesForUser').perform();
  },

  didRender() {
    this._super(...arguments);
    this.get('setForAnimation').perform();
  },

  setForAnimation: task(function*() {
    let leaves = this.get('menuLeaves');
    let lastLeaves = this.get('lastMenuLeaves');
    if (!lastLeaves) {
      yield timeout(0);
      this.set('lastMenuLeaves', leaves);
      return;
    }
    let isAdding = leaves.length > lastLeaves.length;
    let changedLeaf = (isAdding ? leaves : lastLeaves).get('lastObject');
    this.set('isAdding', isAdding);
    this.set('changedLeaf', changedLeaf);
    // if we're adding we want to render immediately an animate it in
    // if we're not adding, we need time to move the item out before
    // a rerender removes it
    yield timeout(isAdding ? 0 : ANIMATION_DURATION);
    this.set('lastMenuLeaves', leaves);
  }),

  namespacePath: computed.alias('namespaceService.path'),

  // this is an array of namespace paths that the current user
  // has access to
  accessibleNamespaces: computed.alias('namespaceService.accessibleNamespaces'),

  namespaceTree: computed('accessibleNamespaces', function() {
    let nsList = this.get('accessibleNamespaces');
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
  menuLeaves: computed('namespacePath', 'namespaceTree', function() {
    let ns = this.get('namespacePath');
    let leaves = ancestorKeysForKey(ns) || [];
    leaves.push(ns);
    return leaves.map(this.pathToLeaf);
  }),

  // the nodes at the root of the namespace tree
  // these will get rendered as the bottom layer
  rootLeaves: computed('namespaceTree', function() {
    let leaves = Object.keys(this.get('namespaceTree'));
    return leaves.map(this.pathToLeaf);
  }),

  currentLeaf: computed.alias('menuLeaves.lastObject'),
  canAccessMultipleNamespaces: computed.gt('accessibleNamespaces.length', 1),
  isUserRootNamespace: computed('auth.authData.userRootNamespace', 'namespacePath', function() {
    return this.get('auth.authData.userRootNamespace') === this.get('namespacePath');
  }),

  namespaceDisplay: computed('namespacePath', 'accessibleNamespaces', 'accessibleNamespaces.[]', function() {
    let namespace = this.get('namespacePath');
    if (namespace === '') {
      return '';
    } else {
      let parts = namespace.split('/');
      return parts[parts.length - 1];
    }
  }),
});
