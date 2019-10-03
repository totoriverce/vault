import Component from '@ember/component';
import { set, get, computed } from '@ember/object';
import { inject as service } from '@ember/service';
import { readOnly } from '@ember/object/computed';
import { task } from 'ember-concurrency';

export default Component.extend({
  namespace: service(),
  store: service(),
  config: null,
  possiblePaths: null,
  currentNamespace: readOnly('namespace.path'),
  namespaces: readOnly('namespaces.accessibleNamespaces'),
  autoCompleteOptions: null,
  setAutoCompleteOptions: task(function*() {
    // fetch auth and secret methods from sys/internal/ui/mounts
    // for any namespaces that are already autocompleted
  }).keepLatest(),

  // singleton mounts are not eligible for per-mount-filtering
  singletonMountTypes: computed(function() {
    return ['cubbyhole', 'system', 'token', 'identity', 'ns_system', 'ns_identity'];
  }),

  actions: {
    onSelectChange() {},

    addOrRemovePath(path, e) {
      let config = get(this, 'config') || [];
      let paths = get(config, 'paths').slice();

      if (e.target.checked) {
        paths.addObject(path);
      } else {
        paths.removeObject(path);
      }

      set(config, 'paths', paths);
    },
  },
});
