import Ember from 'ember';
const { get } = Ember;

// this mixin relies on `unload-model-route` also being used
export default Ember.Mixin.create({
  actions: {
    willTransition(transition) {
      const model = this.controller.get('model');
      if (!model) {
        return true;
      }
      if (get(model, 'hasDirtyAttributes')) {
        if (
          window.confirm(
            'You have unsaved changes. Navigating away will discard these changes. Are you sure you want to discard your changes?'
          )
        ) {
          this.unloadModel();
          return true;
        } else {
          transition.abort();
          return false;
        }
      }
      return true;
    },
  },
});
