import Model, { attr } from '@ember-data/model';
import { withFormFields } from 'vault/decorators/model-form-fields';
import lazyCapabilities, { apiPath } from 'vault/macros/lazy-capabilities';
import { service } from '@ember/service';

const issuerUrls = ['issuingCertificates', 'crlDistributionPoints', 'ocspServers'];
const inputFields = [
  'issuerName',
  'leafNotAfterBehavior',
  'usage',
  'manualChain',
  'revocationSignatureAlgorithm',
  ...issuerUrls,
];
const displayFields = [
  {
    default: [
      'certificate',
      'caChain',
      'commonName',
      'issuerName',
      'issuerId',
      'serialNumber',
      'keyId',
      'altNames',
      'uriSans',
      'ipSans',
      'otherSans',
      'notValidBefore',
      'notValidAfter',
    ],
  },
  { 'Issuer URLs': issuerUrls },
];
@withFormFields(inputFields, displayFields)
export default class PkiIssuerModel extends Model {
  @service secretMountPath;
  // TODO use openAPI after removing route extension (see pki/roles route for example)
  get useOpenAPI() {
    return false;
  }
  get backend() {
    return this.secretMountPath.currentPath;
  }
  get issuerRef() {
    return this.issuerName || this.issuerId;
  }

  // READ ONLY
  @attr isDefault;
  @attr('string', { label: 'Issuer ID' }) issuerId;
  @attr('string', { label: 'Default key ID' }) keyId;
  @attr({ label: 'CA Chain', masked: true }) caChain;
  @attr({ masked: true }) certificate;

  // parsed from certificate contents in serializer (see parse-pki-cert.js)
  @attr commonName;
  @attr('number', { formatDate: true }) notValidAfter;
  @attr('number', { formatDate: true }) notValidBefore;
  @attr serialNumber;
  @attr({ label: 'Subject Alternative Names (SANs)' }) altNames;
  @attr({ label: 'IP SANs' }) ipSans;
  @attr({ label: 'URI SANs' }) uriSans;
  @attr({ label: 'Other SANs' }) otherSans;

  // UPDATING
  @attr('string') issuerName;

  @attr({
    label: 'Leaf notAfter behavior',
    subText:
      'What happens when a leaf certificate is issued, but its NotAfter field (and therefore its expiry date) exceeds that of this issuer.',
    docLink: '/vault/api-docs/secret/pki#update-issuer',
    editType: 'yield',
    valueOptions: ['err', 'truncate', 'permit'],
  })
  leafNotAfterBehavior;

  @attr('string', {
    subText:
      "An advanced field useful when automatic chain building isn't desired. The first element must be the present issuer's reference.",
  })
  manualChain;

  @attr({
    subText: 'Allowed usages for this issuer. It can always be read.',
    editType: 'yield',
    valueOptions: [
      { label: 'Issuing certificates', value: 'issuing-certificates' },
      { label: 'Signing CRLs', value: 'crl-signing' },
      { label: 'Signing OCSPs', value: 'ocsp-signing' },
    ],
  })
  usage;

  @attr({
    subText:
      'The signature algorithm to use when building CRLs. The default value (empty string) is for Go to select the signature algorithm automatically, which may not always work.',
    noDefault: true,
    possibleValues: [
      'sha256withrsa',
      'ecdsawithsha384',
      'sha256withrsapss',
      'ed25519',
      'sha384withrsapss',
      'sha512withrsapss',
      'pureed25519',
      'sha384withrsa',
      'sha512withrsa',
      'ecdsawithsha256',
      'ecdsawithsha512',
    ],
  })
  revocationSignatureAlgorithm;

  @attr('string', {
    subText:
      'The URL values for the Issuing Certificate field. These are different URLs for the same resource, and should be added individually, not in a comma-separated list.',
    editType: 'stringArray',
  })
  issuingCertificates;

  @attr('string', {
    label: 'CRL distribution points',
    subText: 'Specifies the URL values for the CRL Distribution Points field.',
    editType: 'stringArray',
  })
  crlDistributionPoints;

  @attr('string', {
    label: 'OCSP servers',
    subText: 'Specifies the URL values for the OCSP Servers field.',
    editType: 'stringArray',
  })
  ocspServers;

  // IMPORTING
  @attr('string') pemBundle;
  // readonly attrs returned after importing
  @attr importedIssuers;
  @attr importedKeys;
  @attr mapping;

  @lazyCapabilities(apiPath`${'backend'}/issuer/${'issuerId'}`) issuerPath;
  @lazyCapabilities(apiPath`${'backend'}/root/rotate/exported`) rotateExported;
  @lazyCapabilities(apiPath`${'backend'}/root/rotate/internal`) rotateInternal;
  @lazyCapabilities(apiPath`${'backend'}/root/rotate/existing`) rotateExisting;
  @lazyCapabilities(apiPath`${'backend'}/intermediate/cross-sign`) crossSignPath;
  @lazyCapabilities(apiPath`${'backend'}/issuer/${'issuerId'}/sign-intermediate`) signIntermediate;
  get canRotateIssuer() {
    return (
      this.rotateExported.get('canUpdate') !== false ||
      this.rotateExisting.get('canUpdate') !== false ||
      this.rotateInternal.get('canUpdate') !== false
    );
  }
  get canCrossSign() {
    return this.crossSignPath.get('canUpdate') !== false;
  }
  get canSignIntermediate() {
    return this.signIntermediate.get('canUpdate') !== false;
  }
  get canConfigure() {
    return this.issuerPath.get('canUpdate') !== false;
  }
}
