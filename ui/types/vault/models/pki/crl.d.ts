import Model from '@ember-data/model';
import { FormField } from 'vault/app-types';

export default class PkiCrlModel extends Model {
  get useOpenAPI(): boolean;
  getHelpUrl(backendPath: string): string;
  autoRebuild: boolean;
  autoRebuildGracePeriod: string;
  enableDelta: boolean;
  expiry: string;
  deltaRebuildInterval: string;
  disable: boolean;
  ocspExpiry: string;
  ocspDisable: boolean;
  crlPath: string;
  formFields: FormField[];
  get canSet(): boolean;
}
