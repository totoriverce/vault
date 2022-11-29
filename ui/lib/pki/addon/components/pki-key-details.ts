import { action } from '@ember/object';
import Component from '@glimmer/component';
import { inject as service } from '@ember/service';
import errorMessage from 'vault/utils/error-message'
interface Args {
  key: {
    rollbackAttributes: () => void;
    destroyRecord: () => void;
    backend: string;
    keyName: string;
    keyId: string;
  };
}

export default class PkiKeyDetails extends Component<Args> {
  @service declare router: { transitionTo: (route: string) => void; };
  @service declare flashMessages: { success: (successMessage: string) => void; danger: (errorMessage: string) => void; } 

  get breadcrumbs() {
    return [
      { label: 'secrets', route: 'secrets', linkExternal: true },
      { label: this.args.key.backend || 'pki', route: 'overview' },
      { label: 'keys', route: 'keys.index' },
      { label: this.args.key.keyId },
    ];
  }

  @action 
  async deleteKey() {
    try {
      await this.args.key.destroyRecord();
      this.flashMessages.success('Key deleted successfully');
      this.router.transitionTo('vault.cluster.secrets.backend.pki.keys.index');
    } catch (error) {
      this.args.key.rollbackAttributes();
      this.flashMessages.danger( errorMessage(error));
    }
  }
}
