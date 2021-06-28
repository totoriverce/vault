import Component from '@glimmer/component';
import { action } from '@ember/object';

/**
 * @module ConfigureAwsSecretComponent
 *
 * @example
 * ```js
 * <ConfigureAwsSecret 
    @model={{model}} 
    @tab={{tab}} 
    @accessKey={{accessKey}}
    @secretKey={{secretKey}}
    @region={{region}}
    @iamEndpoint={{iamEndpoint}}
    @stsEndpoint={{stsEndpoint}}
    @saveAWSRoot={{action "save" "saveAWSRoot"}}
    @saveAWSLease={{action "save" "saveAWSLease"}} />
 * ```
 *
 * @param {object} model - aws secret engine model
 * @param {string} tab - current tab selection
 * @param {string} accessKey - AWS access key
 * @param {string} secretKey - AWS secret key
 * @param {string} region - AWS region
 * @param {string} iamEndpoint - IAM endpoint
 * @param {string} stsEndpoint - Sts endpoint
 * @param {Function} saveAWSRoot - parent action which saves AWS root credentials
 * @param {Function} saveAWSLease - parent action which updates AWS lease information
 * 
 */
export default class ConfigureAwsSecretComponent extends Component {
  @action
  saveRootCreds(data) {
    this.args.saveAWSRoot(data);
  }

  @action
  saveLease(data) {
    this.args.saveAWSLease(data);
  }
}
