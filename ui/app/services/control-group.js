import Service, { inject as service } from '@ember/service';
import RSVP from 'rsvp';
import ControlGroupError from 'vault/lib/control-group-error';
import getStorage from 'vault/lib/token-storage';
import parseURL from 'core/utils/parse-url';

const CONTROL_GROUP_PREFIX = 'vault:cg-';
const TOKEN_SEPARATOR = '☃';

// list of endpoints that return wrapped responses
// without `wrap-ttl`
const WRAPPED_RESPONSE_PATHS = [
  'sys/wrapping/rewrap',
  'sys/wrapping/wrap',
  'sys/replication/performance/primary/secondary-token',
  'sys/replication/dr/primary/secondary-token',
];

const storageKey = (accessor, path) => {
  return `${CONTROL_GROUP_PREFIX}${accessor}${TOKEN_SEPARATOR}${path}`;
};

export { storageKey, CONTROL_GROUP_PREFIX, TOKEN_SEPARATOR };
export default Service.extend({
  version: service(),
  router: service(),

  storage() {
    return getStorage();
  },

  keyFromAccessor(accessor) {
    let keys = this.storage().keys() || [];
    let returnKey = keys
      .filter(k => k.startsWith(CONTROL_GROUP_PREFIX))
      .find(key => key.replace(CONTROL_GROUP_PREFIX, '').startsWith(accessor));
    return returnKey ? returnKey : null;
  },

  storeControlGroupToken(info) {
    let key = storageKey(info.accessor, info.creation_path);
    this.storage().setItem(key, info);
  },

  deleteControlGroupToken(accessor) {
    this.unmarkTokenForUnwrap();
    let key = this.keyFromAccessor(accessor);
    this.storage().removeItem(key);
  },

  deleteTokens() {
    let keys = this.storage().keys() || [];
    keys.filter(k => k.startsWith(CONTROL_GROUP_PREFIX)).forEach(key => this.storage().removeItem(key));
  },

  wrapInfoForAccessor(accessor) {
    let key = this.keyFromAccessor(accessor);
    return key ? this.storage().getItem(key) : null;
  },

  tokenToUnwrap: null,
  markTokenForUnwrap(accessor) {
    this.set('tokenToUnwrap', this.wrapInfoForAccessor(accessor));
  },

  unmarkTokenForUnwrap() {
    this.set('tokenToUnwrap', null);
  },

  tokenForUrl(url) {
    if (this.version.isOSS) {
      return null;
    }
    let pathForUrl = parseURL(url).pathname;
    pathForUrl = pathForUrl.replace('/v1/', '');
    let tokenInfo = this.tokenToUnwrap;
    if (tokenInfo && tokenInfo.creation_path === pathForUrl) {
      let { token, accessor, creation_time } = tokenInfo;
      return { token, accessor, creationTime: creation_time };
    }
    return null;
  },

  checkForControlGroup(callbackArgs, response, wasWrapTTLRequested) {
    let creationPath = response && response?.wrap_info?.creation_path;
    if (
      this.version.isOSS ||
      wasWrapTTLRequested ||
      !response ||
      (creationPath && WRAPPED_RESPONSE_PATHS.includes(creationPath)) ||
      !response.wrap_info
    ) {
      return RSVP.resolve(...callbackArgs);
    }
    let error = new ControlGroupError(response.wrap_info);
    return RSVP.reject(error);
  },

  paramsFromTransition(transitionTo, params, queryParams) {
    let returnedParams = params.slice();
    let qps = queryParams;
    transitionTo.paramNames.map(name => {
      let param = transitionTo.params[name];
      if (param.length) {
        // push on to the front of the array since were're started at the end
        returnedParams.unshift(param);
      }
    });
    qps = { ...queryParams, ...transitionTo.queryParams };
    // if there's a parent transition, recurse to get its route params
    if (transitionTo.parent) {
      [returnedParams, qps] = this.paramsFromTransition(transitionTo.parent, returnedParams, qps);
    }
    return [returnedParams, qps];
  },

  urlFromTransition(transitionObj) {
    let transition = transitionObj.to;
    let [params, queryParams] = this.paramsFromTransition(transition, [], {});
    let url = this.router.urlFor(transition.name, ...params, {
      queryParams,
    });
    return url.replace('/ui', '');
  },

  handleError(error, transition) {
    let { accessor, token, creation_path, creation_time, ttl } = error;
    let url = this.urlFromTransition(transition);
    let data = { accessor, token, creation_path, creation_time, ttl };
    data.uiParams = { url };
    this.storeControlGroupToken(data);
    return this.router.transitionTo('vault.cluster.access.control-group-accessor', accessor);
  },

  logFromError(error) {
    let { accessor, token, creation_path, creation_time, ttl } = error;
    let data = { accessor, token, creation_path, creation_time, ttl };
    this.storeControlGroupToken(data);

    let href = this.router.urlFor('vault.cluster.access.control-group-accessor', accessor);
    let lines = [
      `A Control Group was encountered at ${error.creation_path}.`,
      `The Control Group Token is ${error.token}.`,
      `The Accessor is ${error.accessor}.`,
      `Visit <a href='${href}'>${href}</a> for more details.`,
    ];
    return {
      type: 'error-with-html',
      content: lines.join('\n'),
    };
  },
});
