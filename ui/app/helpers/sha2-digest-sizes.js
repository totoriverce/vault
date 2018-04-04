import Ember from 'ember';

const SHA2_DIGEST_SIZES = ['sha2-224', 'sha2-256', 'sha2-384', 'sha2-512'];

export function sha2DigestSizes() {
  return SHA2_DIGEST_SIZES;
}

export default Ember.Helper.helper(sha2DigestSizes);
