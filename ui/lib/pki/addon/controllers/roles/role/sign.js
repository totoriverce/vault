/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Controller from '@ember/controller';
import { action } from '@ember/object';
import { tracked } from '@glimmer/tracking';

export default class PkiRolesSignController extends Controller {
  @tracked hasSubmitted = false;

  @action
  toggleTitle() {
    this.hasSubmitted = !this.hasSubmitted;
  }
}
