import { encodePath } from 'vault/utils/path-encoding-helpers';
import PkiCertificateBaseAdapter from './base';

export default class PkiCertificateGenerateAdapter extends PkiCertificateBaseAdapter {
  urlForCreateRecord(modelName, snapshot) {
    const { name, backend } = snapshot.record;
    if (!name || !backend) {
      throw new Error('URL for create record is missing required attributes');
    }
    return `${this.buildURL()}/${encodePath(backend)}/issue/${encodePath(name)}`;
  }
}
