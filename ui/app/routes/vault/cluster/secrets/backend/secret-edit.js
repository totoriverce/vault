import AdapterError from '@ember-data/adapter/error';
import { set } from '@ember/object';
import { resolve } from 'rsvp';
import { inject as service } from '@ember/service';
import Route from '@ember/routing/route';
import utils from 'vault/lib/key-utils';
import UnloadModelRoute from 'vault/mixins/unload-model-route';
import { encodePath, normalizePath } from 'vault/utils/path-encoding-helpers';

export default Route.extend(UnloadModelRoute, {
  pathHelp: service('path-help'),
  secretParam() {
    let { secret } = this.paramsFor(this.routeName);
    return secret ? normalizePath(secret) : '';
  },
  enginePathParam() {
    let { backend } = this.paramsFor('vault.cluster.secrets.backend');
    return backend;
  },
  capabilities(secret, modelType) {
    const backend = this.enginePathParam();
    let backendModel = this.modelFor('vault.cluster.secrets.backend');
    let backendType = backendModel.engineType;
    let path;
    if (backendModel.isV2KV) {
      path = `${backend}/data/${secret}`;
    } else if (backendType === 'transit') {
      path = backend + '/keys/' + secret;
    } else if (backendType === 'ssh' || backendType === 'aws') {
      path = backend + '/roles/' + secret;
    } else if (modelType.startsWith('transform/')) {
      path = this.buildTransformPath(backend, secret, modelType);
    } else {
      path = backend + '/' + secret;
    }
    return this.store.findRecord('capabilities', path);
  },

  buildTransformPath(backend, secret, modelType) {
    let noun = modelType.split('/')[1];
    return `${backend}/${noun}/${secret}`;
  },

  modelTypeForTransform(secretName) {
    if (!secretName) return 'transform';
    if (secretName.startsWith('role/')) {
      return 'transform/role';
    }
    if (secretName.startsWith('template/')) {
      return 'transform/template';
    }
    if (secretName.startsWith('alphabet/')) {
      return 'transform/alphabet';
    }
    return 'transform'; // TODO: transform/transformation;
  },

  transformSecretName(secret, modelType) {
    const noun = modelType.split('/')[1];
    return secret.replace(`${noun}/`, '');
  },

  backendType() {
    return this.modelFor('vault.cluster.secrets.backend').get('engineType');
  },

  templateName: 'vault/cluster/secrets/backend/secretEditLayout',

  beforeModel() {
    let secret = this.secretParam();
    return this.buildModel(secret).then(() => {
      const parentKey = utils.parentKeyForKey(secret);
      const mode = this.routeName.split('.').pop();
      if (mode === 'edit' && utils.keyIsFolder(secret)) {
        if (parentKey) {
          return this.transitionTo('vault.cluster.secrets.backend.list', encodePath(parentKey));
        } else {
          return this.transitionTo('vault.cluster.secrets.backend.list-root');
        }
      }
    });
  },

  buildModel(secret) {
    const backend = this.enginePathParam();

    let modelType = this.modelType(backend, secret);
    if (['secret', 'secret-v2'].includes(modelType)) {
      return resolve();
    }
    return this.pathHelp.getNewModel(modelType, backend);
  },

  modelType(backend, secret) {
    let backendModel = this.modelFor('vault.cluster.secrets.backend', backend);
    let type = backendModel.get('engineType');
    let types = {
      database: secret && secret.startsWith('role/') ? 'database/role' : 'database/connection',
      transit: 'transit-key',
      ssh: 'role-ssh',
      transform: this.modelTypeForTransform(secret),
      aws: 'role-aws',
      pki: secret && secret.startsWith('cert/') ? 'pki-certificate' : 'role-pki',
      cubbyhole: 'secret',
      kv: backendModel.get('modelTypeForKV'),
      generic: backendModel.get('modelTypeForKV'),
    };
    return types[type];
  },

  getTargetVersion(currentVersion, paramsVersion) {
    if (currentVersion) {
      // we have the secret metadata, so we can read the currentVersion but give priority to any
      // version passed in via the url
      return parseInt(paramsVersion || currentVersion, 10);
    } else {
      // we've got a stub model because don't have read access on the metadata endpoint
      return paramsVersion ? parseInt(paramsVersion, 10) : null;
    }
  },

  async fetchV2Models(capabilities, secretModel, params) {
    let backend = this.enginePathParam();
    let backendModel = this.modelFor('vault.cluster.secrets.backend', backend);
    let targetVersion = this.getTargetVersion(secretModel.currentVersion, params.version);

    // if we have the metadata, a list of versions are part of the payload
    let version = secretModel.versions && secretModel.versions.findBy('version', targetVersion);
    // if it didn't fail the server read, and the version is not attached to the metadata,
    // this should 404
    if (!version && secretModel.failedServerRead !== true) {
      let error = new AdapterError();
      set(error, 'httpStatus', 404);
      throw error;
    }
    // manually set the related model
    secretModel.set('engine', backendModel);

    secretModel.set(
      'selectedVersion',
      await this.fetchV2VersionModel(capabilities, secretModel, version, targetVersion)
    );
    return secretModel;
  },

  async fetchV2VersionModel(capabilities, secretModel, version, targetVersion) {
    let secret = this.secretParam();
    let backend = this.enginePathParam();

    // v2 versions have a composite ID, we generated one here if we need to manually set it
    // after a failed fetch later;
    let versionId = targetVersion ? [backend, secret, targetVersion] : [backend, secret];

    let versionModel;
    try {
      if (secretModel.failedServerRead) {
        // we couldn't read metadata, so we want to directly fetch the version

        versionModel =
          this.store.peekRecord('secret-v2-version', JSON.stringify(versionId)) ||
          (await this.store.findRecord('secret-v2-version', JSON.stringify(versionId), {
            reload: true,
          }));
      } else {
        // we may have previously errored, so roll it back here
        version.rollbackAttributes();
        // if metadata read was successful, the version we have is only a partial model
        // trigger reload to fetch the whole version model
        versionModel = await version.reload();
      }
    } catch (error) {
      // cannot read the version data, but can write according to capabilities-self endpoint
      if (error.httpStatus === 403 && capabilities.get('canUpdate')) {
        // versionModel is then a partial model from the metadata (if we have read there), or
        // we need to create one on the client
        if (version) {
          version.set('failedServerRead', true);
          versionModel = version;
        } else {
          this.store.push({
            data: {
              type: 'secret-v2-version',
              id: JSON.stringify(versionId),
              attributes: {
                failedServerRead: true,
              },
            },
          });
          versionModel = this.store.peekRecord('secret-v2-version', JSON.stringify(versionId));
        }
      } else {
        throw error;
      }
    }
    return versionModel;
  },

  handleSecretModelError(capabilities, secretId, modelType, error) {
    // can't read the path and don't have update capability, so re-throw
    if (!capabilities.get('canUpdate') && modelType === 'secret') {
      throw error;
    }
    // don't have access to the metadata for v2 or the secret for v1,
    // so we make a stub model and mark it as `failedServerRead`
    this.store.push({
      data: {
        id: secretId,
        type: modelType,
        attributes: {
          failedServerRead: true,
        },
      },
    });
    let secretModel = this.store.peekRecord(modelType, secretId);
    return secretModel;
  },

  async model(params) {
    let secret = this.secretParam();
    let backend = this.enginePathParam();
    let modelType = this.modelType(backend, secret);
    if (!secret) {
      secret = '\u0020';
    }
    if (modelType === 'pki-certificate') {
      secret = secret.replace('cert/', '');
    }
    if (modelType.startsWith('transform/')) {
      secret = this.transformSecretName(secret, modelType);
    }
    if (modelType === 'database/role') {
      secret = secret.replace('role/', '');
    }
    let secretModel;

    let capabilities = this.capabilities(secret, modelType);
    try {
      secretModel = await this.store.queryRecord(modelType, { id: secret, backend });
    } catch (err) {
      // we've failed the read request, but if it's a kv-type backend, we want to
      // do additional checks of the capabilities
      if (err.httpStatus === 403 && (modelType === 'secret-v2' || modelType === 'secret')) {
        await capabilities;
        secretModel = this.handleSecretModelError(capabilities, secret, modelType, err);
      } else {
        throw err;
      }
    }
    await capabilities;
    if (modelType === 'secret-v2') {
      // after the the base model fetch, kv-v2 has a second associated
      // version model that contains the secret data
      secretModel = await this.fetchV2Models(capabilities, secretModel, params);
    }
    return {
      secret: secretModel,
      capabilities,
    };
  },

  setupController(controller, model) {
    this._super(...arguments);
    let secret = this.secretParam();
    let backend = this.enginePathParam();
    const preferAdvancedEdit =
      /* eslint-disable-next-line ember/no-controller-access-in-routes */
      this.controllerFor('vault.cluster.secrets.backend').get('preferAdvancedEdit') || false;
    const backendType = this.backendType();
    model.secret.setProperties({ backend });
    controller.setProperties({
      model: model.secret,
      capabilities: model.capabilities,
      baseKey: { id: secret },
      // mode will be 'show', 'edit', 'create'
      mode: this.routeName
        .split('.')
        .pop()
        .replace('-root', ''),
      backend,
      preferAdvancedEdit,
      backendType,
    });
  },

  resetController(controller) {
    if (controller.reset && typeof controller.reset === 'function') {
      controller.reset();
    }
  },

  actions: {
    error(error) {
      let secret = this.secretParam();
      let backend = this.enginePathParam();
      set(error, 'keyId', backend + '/' + secret);
      set(error, 'backend', backend);
      return true;
    },

    refreshModel() {
      this.refresh();
    },

    willTransition(transition) {
      /* eslint-disable-next-line ember/no-controller-access-in-routes */
      let { mode, model } = this.controller;
      let version = model.get('selectedVersion');
      let changed = model.changedAttributes();
      let changedKeys = Object.keys(changed);
      // until we have time to move `backend` on a v1 model to a relationship,
      // it's going to dirty the model state, so we need to look for it
      // and explicity ignore it here
      if (
        (mode !== 'show' && (changedKeys.length && changedKeys[0] !== 'backend')) ||
        (mode !== 'show' && version && Object.keys(version.changedAttributes()).length)
      ) {
        if (
          window.confirm(
            'You have unsaved changes. Navigating away will discard these changes. Are you sure you want to discard your changes?'
          )
        ) {
          version && version.rollbackAttributes();
          this.unloadModel();
          return true;
        } else {
          transition.abort();
          return false;
        }
      }
      return this._super(...arguments);
    },
  },
});
