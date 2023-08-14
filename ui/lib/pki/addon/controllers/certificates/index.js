/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Controller from '@ember/controller';
import { getOwner } from '@ember/application';
import { action } from '@ember/object';

export default class PkiCertificatesIndexController extends Controller {
  get mountPoint() {
    return getOwner(this).mountPoint;
  }

  @action setFilter(val) {
    this.filter = val;
  }
  @action setFilterFocus(bool) {
    this.filterFocused = bool;
  }
}
