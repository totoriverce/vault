import { singularize } from 'ember-inflector';
import { inject as service } from '@ember/service';
import Route from '@ember/routing/route';

export default Route.extend({
  pathHelp: service('path-help'),
  model() {
    const { item_id: itemName } = this.paramsFor(this.routeName);
    const { item_type: itemType } = this.paramsFor('vault.cluster.access.method.item');
    const { path: method } = this.paramsFor('vault.cluster.access.method');
    const methodModel = this.modelFor('vault.cluster.access.method');
    const { type } = methodModel;
    const modelType = `generated-${singularize(itemType)}-${type}`;
    return this.store.findRecord(modelType, itemName, {
      adapterOptions: { path: `${method}/${itemType}` },
    });
  },

  setupController(controller, model) {
    this._super(...arguments);
    const { item_id: itemId } = this.paramsFor(this.routeName);
    const { item_type: itemType } = this.paramsFor('vault.cluster.access.method.item');
    let { path } = this.paramsFor('vault.cluster.access.method');
    controller.set('itemType', singularize(itemType));
    controller.set('method', path);
    controller.set('props', model.toJSON());
    controller.set('id', itemId);
  },
});
