import Component from '@glimmer/component';
import { messageTypes } from 'core/helpers/message-types';

/**
 * @module Modal
 * Modal components are used to overlay content on top of the page. Has a darkened background,
 * a title, and in order to close it you must pass an onClose function.
 *
 * @example
 * ```js
 * <Modal @title={'myTitle'} @showCloseButton={true} @onClose={{this.closeModalAndExportData}}/>
 * ```
 * @callback onClose - onClose is the action taken when someone clicks the modal background or close button (if shown).
 * @param {boolean} isActive=false - whether or not modal displays
 * @param {string} [title] - This text shows up in the header section of the modal.
 * @param {boolean} [showCloseButton=false] - controls whether the close button in the top right corner shows.
 * @param {string} type=null - The header type. This comes from the message-types helper.
 */

export default class ModalComponent extends Component {
  get isActive() {
    return this.args.isActive || false;
  }

  get showCloseButton() {
    return this.args.showCloseButton || false;
  }

  get glyph() {
    if (!this.args.type) {
      return null;
    }
    return messageTypes([this.args.type]);
  }

  get modalClass() {
    if (!this.args.type) {
      return 'modal';
    }
    return 'modal ' + messageTypes([this.args.type]).class;
  }
}
