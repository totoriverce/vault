/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@glimmer/component';
import errorMessage from 'vault/utils/error-message';
import { action } from '@ember/object';
import { inject as service } from '@ember/service';
import { task } from 'ember-concurrency';
import { waitFor } from '@ember/test-waiters';
import { tracked } from '@glimmer/tracking';

import RouterService from '@ember/routing/router-service';
import type PkiTidyModel from 'vault/models/pki/tidy';
import type { FormField, TtlEvent } from 'vault/app-types';

interface Args {
  tidy: PkiTidyModel;
  tidyType: string;
  onSave: CallableFunction;
  onCancel: CallableFunction;
}

interface PkiTidyTtls {
  intervalDuration: string;
}
interface PkiTidyBooleans {
  enabled: boolean;
}

export default class PkiTidyForm extends Component<Args> {
  @service declare readonly router: RouterService;

  @tracked errorBanner = '';
  @tracked invalidFormAlert = '';

  @task
  @waitFor
  *save(event: Event) {
    event.preventDefault();
    try {
      yield this.args.tidy.save({ adapterOptions: { tidyType: this.args.tidyType } });
      this.args.onSave();
    } catch (e) {
      this.errorBanner = errorMessage(e);
      this.invalidFormAlert = 'There was an error submitting this form.';
    }
  }

  @action
  handleTtl(attr: FormField, e: TtlEvent) {
    const { enabled, goSafeTimeString } = e;
    const ttlAttr = attr.name;
    this.args.tidy[ttlAttr as keyof PkiTidyTtls] = goSafeTimeString;
    this.args.tidy[attr.options.mapToBoolean as keyof PkiTidyBooleans] = enabled;
  }
}
