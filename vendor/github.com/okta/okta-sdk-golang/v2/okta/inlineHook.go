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
	"github.com/okta/okta-sdk-golang/v2/okta/query"
	"time"
)

type InlineHookResource resource

type InlineHook struct {
	Links       interface{}        `json:"_links,omitempty"`
	Channel     *InlineHookChannel `json:"channel,omitempty"`
	Created     *time.Time         `json:"created,omitempty"`
	Id          string             `json:"id,omitempty"`
	LastUpdated *time.Time         `json:"lastUpdated,omitempty"`
	Name        string             `json:"name,omitempty"`
	Status      string             `json:"status,omitempty"`
	Type        string             `json:"type,omitempty"`
	Version     string             `json:"version,omitempty"`
}

func (m *InlineHookResource) CreateInlineHook(ctx context.Context, body InlineHook) (*InlineHook, *Response, error) {
	url := fmt.Sprintf("/api/v1/inlineHooks")

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	var inlineHook *InlineHook

	resp, err := m.client.requestExecutor.Do(ctx, req, &inlineHook)
	if err != nil {
		return nil, resp, err
	}

	return inlineHook, resp, nil
}

// Gets an inline hook by ID
func (m *InlineHookResource) GetInlineHook(ctx context.Context, inlineHookId string) (*InlineHook, *Response, error) {
	url := fmt.Sprintf("/api/v1/inlineHooks/%v", inlineHookId)

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var inlineHook *InlineHook

	resp, err := m.client.requestExecutor.Do(ctx, req, &inlineHook)
	if err != nil {
		return nil, resp, err
	}

	return inlineHook, resp, nil
}

// Updates an inline hook by ID
func (m *InlineHookResource) UpdateInlineHook(ctx context.Context, inlineHookId string, body InlineHook) (*InlineHook, *Response, error) {
	url := fmt.Sprintf("/api/v1/inlineHooks/%v", inlineHookId)

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("PUT", url, body)
	if err != nil {
		return nil, nil, err
	}

	var inlineHook *InlineHook

	resp, err := m.client.requestExecutor.Do(ctx, req, &inlineHook)
	if err != nil {
		return nil, resp, err
	}

	return inlineHook, resp, nil
}

// Deletes the Inline Hook matching the provided id. Once deleted, the Inline Hook is unrecoverable. As a safety precaution, only Inline Hooks with a status of INACTIVE are eligible for deletion.
func (m *InlineHookResource) DeleteInlineHook(ctx context.Context, inlineHookId string) (*Response, error) {
	url := fmt.Sprintf("/api/v1/inlineHooks/%v", inlineHookId)

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

func (m *InlineHookResource) ListInlineHooks(ctx context.Context, qp *query.Params) ([]*InlineHook, *Response, error) {
	url := fmt.Sprintf("/api/v1/inlineHooks")
	if qp != nil {
		url = url + qp.String()
	}

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var inlineHook []*InlineHook

	resp, err := m.client.requestExecutor.Do(ctx, req, &inlineHook)
	if err != nil {
		return nil, resp, err
	}

	return inlineHook, resp, nil
}

// Executes the Inline Hook matching the provided inlineHookId using the request body as the input. This will send the provided data through the Channel and return a response if it matches the correct data contract. This execution endpoint should only be used for testing purposes.
func (m *InlineHookResource) ExecuteInlineHook(ctx context.Context, inlineHookId string, body InlineHookPayload) (*InlineHookResponse, *Response, error) {
	url := fmt.Sprintf("/api/v1/inlineHooks/%v/execute", inlineHookId)

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	var inlineHookResponse *InlineHookResponse

	resp, err := m.client.requestExecutor.Do(ctx, req, &inlineHookResponse)
	if err != nil {
		return nil, resp, err
	}

	return inlineHookResponse, resp, nil
}

// Activates the Inline Hook matching the provided id
func (m *InlineHookResource) ActivateInlineHook(ctx context.Context, inlineHookId string) (*InlineHook, *Response, error) {
	url := fmt.Sprintf("/api/v1/inlineHooks/%v/lifecycle/activate", inlineHookId)

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("POST", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var inlineHook *InlineHook

	resp, err := m.client.requestExecutor.Do(ctx, req, &inlineHook)
	if err != nil {
		return nil, resp, err
	}

	return inlineHook, resp, nil
}

// Deactivates the Inline Hook matching the provided id
func (m *InlineHookResource) DeactivateInlineHook(ctx context.Context, inlineHookId string) (*InlineHook, *Response, error) {
	url := fmt.Sprintf("/api/v1/inlineHooks/%v/lifecycle/deactivate", inlineHookId)

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("POST", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var inlineHook *InlineHook

	resp, err := m.client.requestExecutor.Do(ctx, req, &inlineHook)
	if err != nil {
		return nil, resp, err
	}

	return inlineHook, resp, nil
}
