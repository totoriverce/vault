import { inject as service } from '@ember/service';
import { or } from '@ember/object/computed';
import { isBlank } from '@ember/utils';
import { task, waitForEvent } from 'ember-concurrency';
import Component from '@ember/component';
import { set } from '@ember/object';
import FocusOnInsertMixin from 'vault/mixins/focus-on-insert';
import keys from 'vault/lib/keycodes';

const LIST_ROOT_ROUTE = 'vault.cluster.secrets.backend.list-root';
const SHOW_ROUTE = 'vault.cluster.secrets.backend.show';

export default Component.extend(FocusOnInsertMixin, {
  router: service(),
  wizard: service(),

  mode: null,
  emptyData: '{\n}',
  onDataChange() {},
  onRefresh() {},
  model: null,
  requestInFlight: or('model.isLoading', 'model.isReloading', 'model.isSaving'),

  didReceiveAttrs() {
    this._super(...arguments);
    if (
      (this.wizard.featureState === 'details' && this.mode === 'create') ||
      (this.wizard.featureState === 'role' && this.mode === 'show')
    ) {
      this.wizard.transitionFeatureMachine(this.wizard.featureState, 'CONTINUE', this.backendType);
    }
    if (this.wizard.featureState === 'displayRole') {
      this.wizard.transitionFeatureMachine(this.wizard.featureState, 'NOOP', this.backendType);
    }
  },

  willDestroyElement() {
    this._super(...arguments);
    if (this.model && this.model.isError) {
      this.model.rollbackAttributes();
    }
  },

  waitForKeyUp: task(function*() {
    while (true) {
      let event = yield waitForEvent(document.body, 'keyup');
      this.onEscape(event);
    }
  })
    .on('didInsertElement')
    .cancelOn('willDestroyElement'),

  transitionToRoute() {
    this.router.transitionTo(...arguments);
  },

  onEscape(e) {
    if (e.keyCode !== keys.ESC || this.mode !== 'show') {
      return;
    }
    this.transitionToRoute(LIST_ROOT_ROUTE);
  },

  hasDataChanges() {
    this.onDataChange(this.model.hasDirtyAttributes);
  },

  persist(method, successCallback) {
    const model = this.model;
    return model[method]().then(() => {
      if (!model.isError) {
        if (this.wizard.featureState === 'role') {
          this.wizard.transitionFeatureMachine('role', 'CONTINUE', this.backendType);
        }
        successCallback(model);
      }
    });
  },

  actions: {
    createOrUpdate(type, event) {
      event.preventDefault();

      const modelId = this.model.id;
      // prevent from submitting if there's no key
      // maybe do something fancier later
      if (type === 'create' && isBlank(modelId)) {
        return;
      }

      this.persist('save', () => {
        this.hasDataChanges();
        this.transitionToRoute(SHOW_ROUTE, modelId);
      });
    },

    setValue(key, event) {
      set(this.model, key, event.target.checked);
    },

    refresh() {
      this.onRefresh();
    },

    delete() {
      this.persist('destroyRecord', () => {
        this.hasDataChanges();
        this.transitionToRoute(LIST_ROOT_ROUTE);
      });
    },

    codemirrorUpdated(attr, val, codemirror) {
      codemirror.performLint();
      const hasErrors = codemirror.state.lint.marked.length > 0;

      if (!hasErrors) {
        set(this.model, attr, JSON.parse(val));
      }
    },
  },
});
