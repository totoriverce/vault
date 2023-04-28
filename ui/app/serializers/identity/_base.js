/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { assign } from '@ember/polyfills';
import ApplicationSerializer from '../application';

export default ApplicationSerializer.extend({
  normalizeItems(payload) {
    if (payload.data.keys && Array.isArray(payload.data.keys)) {
      return payload.data.keys.map((key) => {
        const model = payload.data.key_info[key];
        model.id = key;
        return model;
      });
    }
    assign(payload, payload.data);
    delete payload.data;
    return payload;
  },
});
