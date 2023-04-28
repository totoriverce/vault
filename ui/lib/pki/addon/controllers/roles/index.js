/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Controller from '@ember/controller';
import { getOwner } from '@ember/application';

export default class PkiRolesIndexController extends Controller {
  get mountPoint() {
    return getOwner(this).mountPoint;
  }
}
