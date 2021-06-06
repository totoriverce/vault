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
	"context"
	"fmt"
	"time"
)

type GroupRuleResource resource

type GroupRule struct {
	Actions     *GroupRuleAction     `json:"actions,omitempty"`
	Conditions  *GroupRuleConditions `json:"conditions,omitempty"`
	Created     *time.Time           `json:"created,omitempty"`
	Id          string               `json:"id,omitempty"`
	LastUpdated *time.Time           `json:"lastUpdated,omitempty"`
	Name        string               `json:"name,omitempty"`
	Status      string               `json:"status,omitempty"`
	Type        string               `json:"type,omitempty"`
}

// Updates a group rule. Only &#x60;INACTIVE&#x60; rules can be updated.
func (m *GroupRuleResource) UpdateGroupRule(ctx context.Context, ruleId string, body GroupRule) (*GroupRule, *Response, error) {
	url := fmt.Sprintf("/api/v1/groups/rules/%v", ruleId)

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("PUT", url, body)
	if err != nil {
		return nil, nil, err
	}

	var groupRule *GroupRule

	resp, err := m.client.requestExecutor.Do(ctx, req, &groupRule)
	if err != nil {
		return nil, resp, err
	}

	return groupRule, resp, nil
}

// Removes a specific group rule by id from your organization
func (m *GroupRuleResource) DeleteGroupRule(ctx context.Context, ruleId string) (*Response, error) {
	url := fmt.Sprintf("/api/v1/groups/rules/%v", ruleId)

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := m.client.requestExecutor.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
