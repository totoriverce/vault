import Component from '@glimmer/component';
import { tracked } from '@glimmer/tracking';
import { inject as service } from '@ember/service';

/**
 * @module MfaLoginEnforcementHeader
 * MfaLoginEnforcementHeader components are used to display information when creating and editing login enforcements
 *
 * @example
 * ```js
 * <MfaLoginEnforcementHeader @heading="New enforcement" />
 * <MfaLoginEnforcementHeader @radioCardGroupValue={{this.enforcementPreference}} @onRadioCardSelect={{fn (mut this.enforcementPreference)}} @onEnforcementSelect={{fn (mut this.enforcement)}} />
 * ```
 * @callback onRadioCardSelect
 * @callback onEnforcementSelect
 * @param {boolean} [isInline] - toggle component display when used inline with mfa method form -- overrides heading and shows radio cards and enforcement select
 * @param {string} [heading] - page heading to display outside of inline mode
 * @param {string} [radioCardGroupValue] - selected value of the radio card group in inline mode -- new, existing or skip are the accepted values
 * @param {onRadioCardSelect} [onRadioCardSelect] - change event triggered on radio card select
 * @param {onEnforcementSelect} [onEnforcementSelect] - change event triggered on enforcement select when radioCardGroupValue is set to existing
 */

export default class MfaLoginEnforcementHeaderComponent extends Component {
  @service store;

  constructor() {
    super(...arguments);
    if (!this.args.heading) {
      this.fetchEnforcements();
    }
  }

  @tracked enforcements = [];

  async fetchEnforcements() {
    try {
      this.enforcements = (await this.store.query('mfa-login-enforcement', {})).toArray();
    } catch (error) {
      this.enforcements = [];
    }
  }
}
