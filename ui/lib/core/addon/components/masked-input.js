/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { debug } from '@ember/debug';
import { action } from '@ember/object';
import { guidFor } from '@ember/object/internals';
import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import autosize from 'autosize';

/**
 * @module MaskedInput
 * `MaskedInput` components are textarea inputs where the input is hidden. They are used to enter sensitive information like passwords.
 *
 * @example
 * <MaskedInput
 *  @value={{get @model attr.name}}
 *  @allowCopy={{true}}
 *  @allowDownload={{true}}
 *  @onChange={{this.handleChange}}
 *  @onKeyUp={{this.handleKeyUp}}
 * />
 *
 * @param value {String} - The value to display in the input.
 * @param name {String} - The key correlated to the value. Used for the download file name.
 * @param [onChange=Callback] {Function|action} - Callback triggered on change, sends new value. Must set the value of @value
 * @param [allowCopy=false] {bool} - Whether or not the input should render with a copy button.
 * @param [displayOnly=false] {bool} - Whether or not to display the value as a display only `pre` element or as an input.
 *
 */
export default class MaskedInputComponent extends Component {
  textareaId = 'textarea-' + guidFor(this);
  @tracked showValue = false;

  constructor() {
    super(...arguments);
    if (!this.args.onChange && !this.args.displayOnly) {
      debug('onChange is required for editable Masked Input!');
    }
    this.updateSize();
  }

  updateSize() {
    autosize(document.getElementById(this.textareaId));
  }

  @action onChange(evt) {
    const value = evt.target.value;
    if (this.args.onChange) {
      this.args.onChange(value);
    }
  }

  @action handleKeyUp(name, value) {
    this.updateSize();
    if (this.onKeyUp) {
      this.onKeyUp(name, value);
    }
  }
  @action toggleMask() {
    this.showValue = !this.showValue;
  }
}
