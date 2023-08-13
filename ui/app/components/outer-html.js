/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

// THIS COMPONENT IS ONLY FOR EXTENDING
// You should use this component if you want to use outerHTML semantics
// in your components - this is the default for upcoming Glimmer components
import Component from '@ember/component';

export default Component.extend({
  tagName: '',
});

// yep! that's it, it's more of a way to keep track of what components
// use tagless semantics to make the upgrade to glimmer components easier
