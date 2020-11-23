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

type U2fUserFactorProfile struct {
	CredentialId string `json:"credentialId,omitempty"`
}

func NewU2fUserFactorProfile() *U2fUserFactorProfile {
	return &U2fUserFactorProfile{}
}

func (a *U2fUserFactorProfile) IsUserFactorInstance() bool {
	return true
}
