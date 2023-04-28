/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import { Factory } from 'ember-cli-mirage';

export default Factory.extend({
  use_signature: true,
  idp_url: 'https://foobar.pingidentity.com/pingid',
  admin_url: 'https://foobar.pingidentity.com/pingid',
  authenticator_url: 'https://authenticator.pingone.com/pingid/ppm',
  org_alias: 'foobarbaz',
  type: 'pingid',
  username_template: '',
  namespace_id: 'root',
});
