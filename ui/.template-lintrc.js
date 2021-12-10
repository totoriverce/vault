'use strict';

module.exports = {
  extends: 'octane',
  rules: {
    'no-bare-strings': 'off',
    'no-action': 'off',
    'no-duplicate-landmark-elements': 'warn',
    'no-implicit-this': {
      allow: ['supported-auth-backends'],
    },
    'require-input-label': 'off',
    'no-down-event-binding': 'warn',
  },
  ignore: ['lib/story-md'],
};
