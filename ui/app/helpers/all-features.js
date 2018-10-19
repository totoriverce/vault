import { helper as buildHelper } from '@ember/component/helper';

const ALL_FEATURES = [
  'HSM',
  'Performance Replication',
  'DR Replication',
  'MFA',
  'Sentinel',
  'AWS KMS Autounseal',
  'GCP CKMS Autounseal',
  'Seal Wrapping',
  'Control Groups',
  'Azure Key Vault Seal',
  'Performance Standby',
  'Namespaces',
];

export function allFeatures() {
  return ALL_FEATURES;
}

export default buildHelper(allFeatures);
