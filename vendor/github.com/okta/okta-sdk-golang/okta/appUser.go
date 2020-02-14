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

import (
	"fmt"
	"github.com/okta/okta-sdk-golang/okta/query"
	"time"
)

type AppUserResource resource

type AppUser struct {
	Embedded        interface{}         `json:"_embedded,omitempty"`
	Links           interface{}         `json:"_links,omitempty"`
	Created         *time.Time          `json:"created,omitempty"`
	Credentials     *AppUserCredentials `json:"credentials,omitempty"`
	ExternalId      string              `json:"externalId,omitempty"`
	Id              string              `json:"id,omitempty"`
	LastSync        *time.Time          `json:"lastSync,omitempty"`
	LastUpdated     *time.Time          `json:"lastUpdated,omitempty"`
	PasswordChanged *time.Time          `json:"passwordChanged,omitempty"`
	Profile         interface{}         `json:"profile,omitempty"`
	Scope           string              `json:"scope,omitempty"`
	Status          string              `json:"status,omitempty"`
	StatusChanged   *time.Time          `json:"statusChanged,omitempty"`
	SyncState       string              `json:"syncState,omitempty"`
}

func NewAppUser() *AppUser {
	return &AppUser{}
}

func (a *AppUser) IsApplicationInstance() bool {
	return true
}

func (m *AppUserResource) UpdateApplicationUser(appId string, userId string, body AppUser) (*AppUser, *Response, error) {
	url := fmt.Sprintf("/api/v1/apps/%v/users/%v", appId, userId)
	req, err := m.client.requestExecutor.NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	var appUser *AppUser
	resp, err := m.client.requestExecutor.Do(req, &appUser)
	if err != nil {
		return nil, resp, err
	}
	return appUser, resp, nil
}
func (m *AppUserResource) DeleteApplicationUser(appId string, userId string, qp *query.Params) (*Response, error) {
	url := fmt.Sprintf("/api/v1/apps/%v/users/%v", appId, userId)
	if qp != nil {
		url = url + qp.String()
	}
	req, err := m.client.requestExecutor.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := m.client.requestExecutor.Do(req, nil)
	if err != nil {
		return resp, err
	}
	return resp, nil
}
