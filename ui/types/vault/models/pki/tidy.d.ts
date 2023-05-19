import Model from '@ember-data/model';
import { FormField, FormFieldGroups } from 'vault/vault/app-types';

export default class PkiTidyModel extends Model {
  version: string;
  acmeAccountSafetyBuffer: string;
  tidyAcme: boolean;
  enabled: boolean;
  intervalDuration: string;
  issuerSafetyBuffer: string;
  pauseDuration: string;
  revocationQueueSafetyBuffer: string;
  safetyBuffer: string;
  tidyCertStore: boolean;
  tidyCrossClusterRevokedCerts: boolean;
  tidyExpiredIssuers: boolean;
  tidyMoveLegacyCaBundle: boolean;
  tidyRevocationQueue: boolean;
  tidyRevokedCertIssuerAssociations: boolean;
  tidyRevokedCerts: boolean;
  get useOpenAPI(): boolean;
  getHelpUrl(backend: string): string;
  _allByKey: {
    intervalDuration: FormField[];
  };
  get allGroups(): FormFieldGroups[];
  get sharedFields(): FormFieldGroups[];
  get formFieldGroups(): FormFieldGroups[];
}
