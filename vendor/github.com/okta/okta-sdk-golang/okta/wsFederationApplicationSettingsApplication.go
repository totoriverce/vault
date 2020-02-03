/*
* Copyright 2018 - Present Okta, Inc.
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*      http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

// AUTO-GENERATED!  DO NOT EDIT FILE DIRECTLY

package okta

import ()

type WsFederationApplicationSettingsApplication struct {
	AttributeStatements  string `json:"attributeStatements,omitempty"`
	AudienceRestriction  string `json:"audienceRestriction,omitempty"`
	AuthnContextClassRef string `json:"authnContextClassRef,omitempty"`
	GroupFilter          string `json:"groupFilter,omitempty"`
	GroupName            string `json:"groupName,omitempty"`
	GroupValueFormat     string `json:"groupValueFormat,omitempty"`
	NameIDFormat         string `json:"nameIDFormat,omitempty"`
	Realm                string `json:"realm,omitempty"`
	SiteURL              string `json:"siteURL,omitempty"`
	UsernameAttribute    string `json:"usernameAttribute,omitempty"`
	WReplyOverride       *bool  `json:"wReplyOverride,omitempty"`
	WReplyURL            string `json:"wReplyURL,omitempty"`
}

func NewWsFederationApplicationSettingsApplication() *WsFederationApplicationSettingsApplication {
	return &WsFederationApplicationSettingsApplication{}
}

func (a *WsFederationApplicationSettingsApplication) IsApplicationInstance() bool {
	return true
}
