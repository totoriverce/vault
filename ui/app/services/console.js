/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

// Low level service that allows users to input paths to make requests to vault
// this service provides the UI synecdote to the cli commands read, write, delete, and list
import { filterBy } from '@ember/object/computed';

import Service from '@ember/service';

import { getOwner } from '@ember/application';
import { computed } from '@ember/object';
import { shiftCommandIndex } from 'vault/lib/console-helpers';
import { encodePath } from 'vault/utils/path-encoding-helpers';

export function sanitizePath(path) {
  //remove whitespace + remove trailing and leading slashes
  return path.trim().replace(/^\/+|\/+$/g, '');
}
export function ensureTrailingSlash(path) {
  return path.replace(/(\w+[^/]$)/g, '$1/');
}

const VERBS = {
  read: 'GET',
  list: 'GET',
  write: 'POST',
  delete: 'DELETE',
};

export default Service.extend({
  isOpen: false,

  adapter() {
    return getOwner(this).lookup('adapter:console');
  },
  commandHistory: filterBy('log', 'type', 'command'),
  log: computed(function () {
    return [];
  }),
  commandIndex: null,

  shiftCommandIndex(keyCode, setCommandFn = () => {}) {
    const [newIndex, newCommand] = shiftCommandIndex(keyCode, this.commandHistory, this.commandIndex);
    if (newCommand !== undefined && newIndex !== undefined) {
      this.set('commandIndex', newIndex);
      setCommandFn(newCommand);
    }
  },

  clearLog(clearAll = false) {
    const log = this.log;
    let history;
    if (!clearAll) {
      history = this.commandHistory.slice();
      history.setEach('hidden', true);
    }
    log.clear();
    if (history) {
      log.addObjects(history);
    }
  },

  logAndOutput(command, logContent) {
    const log = this.log;
    if (command) {
      log.pushObject({ type: 'command', content: command });
      this.set('commandIndex', null);
    }
    if (logContent) {
      log.pushObject(logContent);
    }
  },

  ajax(operation, path, options = {}) {
    const verb = VERBS[operation];
    const adapter = this.adapter();
    const url = adapter.buildURL(encodePath(path));
    const { data, wrapTTL } = options;
    return adapter.ajax(url, verb, {
      data,
      wrapTTL,
    });
  },

  read(path, data, wrapTTL) {
    return this.ajax('read', sanitizePath(path), { wrapTTL });
  },

  write(path, data, wrapTTL) {
    return this.ajax('write', sanitizePath(path), { data, wrapTTL });
  },

  delete(path) {
    return this.ajax('delete', sanitizePath(path));
  },

  list(path, data, wrapTTL) {
    const listPath = ensureTrailingSlash(sanitizePath(path));
    return this.ajax('list', listPath, {
      data: {
        list: true,
      },
      wrapTTL,
    });
  },
});
