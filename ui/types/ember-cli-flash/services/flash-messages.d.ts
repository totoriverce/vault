/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

declare module 'ember-cli-flash/services/flash-messages' {
  import Service from '@ember/service';
  import FlashObject from 'ember-cli-flash/flash/object';
  import { A } from '@ember/array';

  type Partial<T> = { [K in keyof T]?: T[K] };

  interface MessageOptions {
    type: string;
    priority: number;
    timeout: number;
    sticky: boolean;
    showProgress: boolean;
    extendedTimeout: number;
    destroyOnClick: boolean;
    onDestroy: () => void;
    [key: string]: unknown;
  }

  interface CustomMessageInfo extends Partial<MessageOptions> {
    message: string;
  }

  interface FlashFunction {
    (message: string, options?: Partial<MessageOptions>): FlashMessages;
  }

  class FlashMessages extends Service {
    queue: A<FlashObject>;
    success: FlashFunction;
    warning: FlashFunction;
    info: FlashFunction;
    error: FlashFunction;
    danger: FlashFunction;
    alert: FlashFunction;
    secondary: FlashFunction;
    add(messageInfo: CustomMessageInfo): FlashMessages;
    clearMessages(): FlashMessages;
    registerTypes(types: string[]): FlashMessages;
    getFlashObject(): FlashObject;
  }

  export default FlashMessages;
}
