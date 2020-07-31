/**
 * @module BoxRadio
 * BoxRadio components are used to...
 *
 * @example
 * ```js
 * <BoxRadio @requiredParam={requiredParam} @optionalParam={optionalParam} @param1={{param1}}/>
 * ```
 * @param {object} requiredParam - requiredParam is...
 * @param {string} [optionalParam] - optionalParam is...
 * @param {string} [param1=defaultValue] - param1 is...
 */

import Component from '@ember/component';
import layout from '../templates/components/box-radio';

export default Component.extend({
  layout,
});
