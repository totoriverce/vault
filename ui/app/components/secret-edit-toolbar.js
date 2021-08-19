// ARG TODO turn into Glimmer
/**
 * @module SecretEditToolbar
 * SecretEditToolbar component is the toolbar component displaying the JSON toggle and the actions like delete in the show mode.
 *
 * @example
 * ```js
 * <SecretEditToolbar
 * @mode={{mode}}
 * @model={{this.model}}
 * @isV2={{isV2}}
 * @isWriteWithoutRead={{isWriteWithoutRead}}
 * @secretDataIsAdvanced={{secretDataIsAdvanced}}
 * @showAdvancedMode={{showAdvancedMode}}
 * @modelForData={{this.modelForData}}
 * @navToNearestAncestor={{this.navToNearestAncestor}}
 * @canUpdateSecretData={{canUpdateSecretData}}
 * @codemirrorString={{codemirrorString}}
 * @wrappedData={{wrappedData}}
 * @editActions={{hash
    toggleAdvanced=(action "toggleAdvanced")
    refresh=(action "refresh")
  }}
 * />
 * ```
 
 * @param {string} mode - show, create, edit. The view.
 * @param {object} model - the model passed from the parent secret-edit
 * @param {boolean} isV2 - KV type
 * @param {boolean} isWriteWithoutRead - boolean describing permissions
 * @param {boolean} secretDataIsAdvanced - used to determine if show JSON toggle
 * @param {object} modelForData - a modified version of the model with secret data
 * @param {string} navToNearestAncestor - route to nav to if press cancel
 * @param {boolean} canUpdateSecretData - permissions that show the create new version button or not.
 * @param {boolean} canEdit - permissions
 * @param {string} codemirrorString - used in JSON editor
 * @param {object} wrappedData - when copy the data it's the token of the secret returned.
 * @param {object} editActions - actions passed from parent to child
 */

import Component from '@ember/component';
import { not } from '@ember/object/computed';
import { inject as service } from '@ember/service';

export default Component.extend({
  store: service(),

  wrappedData: null,
  isWrapping: false,
  showWrapButton: not('wrappedData'),

  actions: {
    handleWrapClick() {
      this.set('isWrapping', true);
      if (this.isV2) {
        this.store
          .adapterFor('secret-v2-version')
          .queryRecord(this.modelForData.id, { wrapTTL: 1800 })
          .then(resp => {
            this.set('wrappedData', resp.wrap_info.token);
            this.flashMessages.success('Secret Successfully Wrapped!');
          })
          .catch(() => {
            this.flashMessages.danger('Could Not Wrap Secret');
          })
          .finally(() => {
            this.set('isWrapping', false);
          });
      } else {
        this.store
          .adapterFor('secret')
          .queryRecord(null, null, { backend: this.model.backend, id: this.modelForData.id, wrapTTL: 1800 })
          .then(resp => {
            this.set('wrappedData', resp.wrap_info.token);
            this.flashMessages.success('Secret Successfully Wrapped!');
          })
          .catch(() => {
            this.flashMessages.danger('Could Not Wrap Secret');
          })
          .finally(() => {
            this.set('isWrapping', false);
          });
      }
    },

    clearWrappedData() {
      console.log(this.wrappedData, 'meep');
      this.set('wrappedData', null);
    },

    handleCopySuccess() {
      this.flashMessages.success('Copied Wrapped Data!');
      this.send('clearWrappedData');
    },

    handleCopyError() {
      this.flashMessages.danger('Could Not Copy Wrapped Data');
      this.send('clearWrappedData');
    },
  },
});
