/* eslint-env node */
'use strict';

module.exports = function(environment) {
  var ENV = {
    modulePrefix: 'vault',
    environment: environment,
    rootURL: '/ui/',
    serviceWorkerScope: '/v1/sys/storage/raft/snapshot',
    locationType: 'auto',
    EmberENV: {
      FEATURES: {
        // Here you can enable experimental features on an ember canary build
        // e.g. 'with-controller': true
      },
      EXTEND_PROTOTYPES: {
        // Prevent Ember Data from overriding Date.parse.
        Date: false,
      },
    },

    APP: {
      // endpoints that the UI polls
      POLLING_URLS: ['sys/health', 'sys/replication/status', 'sys/seal-status'],
      // endpoints that UI uses to determine the cluster state
      // calls to these endpoints will always go to the root namespace
      // these also need to be updated in the open-api-explorer engine
      NAMESPACE_ROOT_URLS: [
        'sys/health',
        'sys/seal-status',
        'sys/license/features',
        'sys/internal/counters/config',
      ],
      // number of records to show on a single page by default - this is used by the client-side pagination
      DEFAULT_PAGE_SIZE: 100,
    },
    flashMessageDefaults: {
      timeout: 7000,
      sticky: false,
    },
  };
  if (environment === 'development') {
    // ENV.APP.LOG_RESOLVER = true;
    // ENV.APP.LOG_ACTIVE_GENERATION = true;
    ENV.APP.LOG_TRANSITIONS = true;
    // ENV.APP.LOG_TRANSITIONS_INTERNAL = true;
    // ENV.APP.LOG_VIEW_LOOKUPS = true;
    // ENV['ember-cli-mirage'] = {
    //   enabled: true,
    // };
  }

  if (environment === 'test') {
    // Testem prefers this...
    ENV.locationType = 'none';
    // keep test console output quieter
    ENV.APP.LOG_ACTIVE_GENERATION = false;
    ENV.APP.LOG_VIEW_LOOKUPS = false;
    ENV.APP.rootElement = '#ember-testing';
    ENV.APP.autoboot = false;
    ENV.flashMessageDefaults.timeout = 50;
    ENV['ember-cli-mirage'] = {
      enabled: false,
    };
  }
  if (environment !== 'production') {
    ENV.APP.DEFAULT_PAGE_SIZE = 15;
    ENV.contentSecurityPolicyHeader = 'Content-Security-Policy';
    ENV.contentSecurityPolicyMeta = true;
    ENV.contentSecurityPolicy = {
      'connect-src': ["'self'"],
      'img-src': ["'self'", 'data:'],
      'font-src': ["'self'"],
      'form-action': ["'none'"],
      'script-src': ["'self'"],
      'style-src': ["'unsafe-inline'", "'self'"],
    };
  }

  ENV.welcomeMessage = process.env.UI_AUTH_WELCOME;

  return ENV;
};
