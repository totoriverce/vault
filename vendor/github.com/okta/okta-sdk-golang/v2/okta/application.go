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

type App interface {
	IsApplicationInstance() bool
}

type ApplicationResource resource

type Application struct {
	Embedded      interface{}               `json:"_embedded,omitempty"`
	Links         interface{}               `json:"_links,omitempty"`
	Accessibility *ApplicationAccessibility `json:"accessibility,omitempty"`
	Created       *time.Time                `json:"created,omitempty"`
	Credentials   *ApplicationCredentials   `json:"credentials,omitempty"`
	Features      []string                  `json:"features,omitempty"`
	Id            string                    `json:"id,omitempty"`
	Label         string                    `json:"label,omitempty"`
	LastUpdated   *time.Time                `json:"lastUpdated,omitempty"`
	Licensing     *ApplicationLicensing     `json:"licensing,omitempty"`
	Name          string                    `json:"name,omitempty"`
	Profile       interface{}               `json:"profile,omitempty"`
	Settings      *ApplicationSettings      `json:"settings,omitempty"`
	SignOnMode    string                    `json:"signOnMode,omitempty"`
	Status        string                    `json:"status,omitempty"`
	Visibility    *ApplicationVisibility    `json:"visibility,omitempty"`
}

func NewApplication() *Application {
	return &Application{}
}

func (a *Application) IsApplicationInstance() bool {
	return true
}

// Fetches an application from your Okta organization by &#x60;id&#x60;.
func (m *ApplicationResource) GetApplication(ctx context.Context, appId string, appInstance App, qp *query.Params) (App, *Response, error) {
	url := fmt.Sprintf("/api/v1/apps/%v", appId)
	if qp != nil {
		url = url + qp.String()
	}

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	application := appInstance

	resp, err := m.client.requestExecutor.Do(ctx, req, &application)
	if err != nil {
		return nil, resp, err
	}

	return application, resp, nil
}

// Updates an application in your organization.
func (m *ApplicationResource) UpdateApplication(ctx context.Context, appId string, body App) (App, *Response, error) {
	url := fmt.Sprintf("/api/v1/apps/%v", appId)

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("PUT", url, body)
	if err != nil {
		return nil, nil, err
	}

	application := body

	resp, err := m.client.requestExecutor.Do(ctx, req, &application)
	if err != nil {
		return nil, resp, err
	}

	return application, resp, nil
}

// Removes an inactive application.
func (m *ApplicationResource) DeleteApplication(ctx context.Context, appId string) (*Response, error) {
	url := fmt.Sprintf("/api/v1/apps/%v", appId)

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

// Enumerates apps added to your organization with pagination. A subset of apps can be returned that match a supported filter expression or query.
func (m *ApplicationResource) ListApplications(ctx context.Context, qp *query.Params) ([]App, *Response, error) {
	url := fmt.Sprintf("/api/v1/apps")
	if qp != nil {
		url = url + qp.String()
	}

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var application []Application

	resp, err := m.client.requestExecutor.Do(ctx, req, &application)
	if err != nil {
		return nil, resp, err
	}

	apps := make([]App, len(application))
	for i := range application {
		apps[i] = &application[i]
	}
	return apps, resp, nil

}

// Adds a new application to your Okta organization.
func (m *ApplicationResource) CreateApplication(ctx context.Context, body App, qp *query.Params) (App, *Response, error) {
	url := fmt.Sprintf("/api/v1/apps")
	if qp != nil {
		url = url + qp.String()
	}

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	application := body

	resp, err := m.client.requestExecutor.Do(ctx, req, &application)
	if err != nil {
		return nil, resp, err
	}

	return application, resp, nil
}

// Enumerates Certificate Signing Requests for an application
func (m *ApplicationResource) ListCsrsForApplication(ctx context.Context, appId string) ([]*Csr, *Response, error) {
	url := fmt.Sprintf("/api/v1/apps/%v/credentials/csrs", appId)

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var csr []*Csr

	resp, err := m.client.requestExecutor.Do(ctx, req, &csr)
	if err != nil {
		return nil, resp, err
	}

	return csr, resp, nil
}

// Generates a new key pair and returns the Certificate Signing Request for it.
func (m *ApplicationResource) GenerateCsrForApplication(ctx context.Context, appId string, body CsrMetadata) (*Csr, *Response, error) {
	url := fmt.Sprintf("/api/v1/apps/%v/credentials/csrs", appId)

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	var csr *Csr

	resp, err := m.client.requestExecutor.Do(ctx, req, &csr)
	if err != nil {
		return nil, resp, err
	}

	return csr, resp, nil
}

func (m *ApplicationResource) RevokeCsrFromApplication(ctx context.Context, appId string, csrId string) (*Response, error) {
	url := fmt.Sprintf("/api/v1/apps/%v/credentials/csrs/%v", appId, csrId)

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

func (m *ApplicationResource) GetCsrForApplication(ctx context.Context, appId string, csrId string) (*Csr, *Response, error) {
	url := fmt.Sprintf("/api/v1/apps/%v/credentials/csrs/%v", appId, csrId)

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var csr *Csr

	resp, err := m.client.requestExecutor.Do(ctx, req, &csr)
	if err != nil {
		return nil, resp, err
	}

	return csr, resp, nil
}

func (m *ApplicationResource) PublishCerCert(ctx context.Context, appId string, csrId string, body string) (*JsonWebKey, *Response, error) {
	url := fmt.Sprintf("/api/v1/apps/%v/credentials/csrs/%v/lifecycle/publish", appId, csrId)

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/x-x509-ca-cert").NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	var jsonWebKey *JsonWebKey

	resp, err := m.client.requestExecutor.Do(ctx, req, &jsonWebKey)
	if err != nil {
		return nil, resp, err
	}

	return jsonWebKey, resp, nil
}

func (m *ApplicationResource) PublishBinaryCerCert(ctx context.Context, appId string, csrId string, body string) (*JsonWebKey, *Response, error) {
	url := fmt.Sprintf("/api/v1/apps/%v/credentials/csrs/%v/lifecycle/publish", appId, csrId)

	req, err := m.client.requestExecutor.AsBinary().WithAccept("application/json").WithContentType("application/x-x509-ca-cert").NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	var jsonWebKey *JsonWebKey

	resp, err := m.client.requestExecutor.Do(ctx, req, &jsonWebKey)
	if err != nil {
		return nil, resp, err
	}

	return jsonWebKey, resp, nil
}

func (m *ApplicationResource) PublishDerCert(ctx context.Context, appId string, csrId string, body string) (*JsonWebKey, *Response, error) {
	url := fmt.Sprintf("/api/v1/apps/%v/credentials/csrs/%v/lifecycle/publish", appId, csrId)

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/pkix-cert").NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	var jsonWebKey *JsonWebKey

	resp, err := m.client.requestExecutor.Do(ctx, req, &jsonWebKey)
	if err != nil {
		return nil, resp, err
	}

	return jsonWebKey, resp, nil
}

func (m *ApplicationResource) PublishBinaryDerCert(ctx context.Context, appId string, csrId string, body string) (*JsonWebKey, *Response, error) {
	url := fmt.Sprintf("/api/v1/apps/%v/credentials/csrs/%v/lifecycle/publish", appId, csrId)

	req, err := m.client.requestExecutor.AsBinary().WithAccept("application/json").WithContentType("application/pkix-cert").NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	var jsonWebKey *JsonWebKey

	resp, err := m.client.requestExecutor.Do(ctx, req, &jsonWebKey)
	if err != nil {
		return nil, resp, err
	}

	return jsonWebKey, resp, nil
}

func (m *ApplicationResource) PublishBinaryPemCert(ctx context.Context, appId string, csrId string, body string) (*JsonWebKey, *Response, error) {
	url := fmt.Sprintf("/api/v1/apps/%v/credentials/csrs/%v/lifecycle/publish", appId, csrId)

	req, err := m.client.requestExecutor.AsBinary().WithAccept("application/json").WithContentType("application/x-pem-file").NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	var jsonWebKey *JsonWebKey

	resp, err := m.client.requestExecutor.Do(ctx, req, &jsonWebKey)
	if err != nil {
		return nil, resp, err
	}

	return jsonWebKey, resp, nil
}

// Enumerates key credentials for an application
func (m *ApplicationResource) ListApplicationKeys(ctx context.Context, appId string) ([]*JsonWebKey, *Response, error) {
	url := fmt.Sprintf("/api/v1/apps/%v/credentials/keys", appId)

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var jsonWebKey []*JsonWebKey

	resp, err := m.client.requestExecutor.Do(ctx, req, &jsonWebKey)
	if err != nil {
		return nil, resp, err
	}

	return jsonWebKey, resp, nil
}

// Generates a new X.509 certificate for an application key credential
func (m *ApplicationResource) GenerateApplicationKey(ctx context.Context, appId string, qp *query.Params) (*JsonWebKey, *Response, error) {
	url := fmt.Sprintf("/api/v1/apps/%v/credentials/keys/generate", appId)
	if qp != nil {
		url = url + qp.String()
	}

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("POST", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var jsonWebKey *JsonWebKey

	resp, err := m.client.requestExecutor.Do(ctx, req, &jsonWebKey)
	if err != nil {
		return nil, resp, err
	}

	return jsonWebKey, resp, nil
}

// Gets a specific application key credential by kid
func (m *ApplicationResource) GetApplicationKey(ctx context.Context, appId string, keyId string) (*JsonWebKey, *Response, error) {
	url := fmt.Sprintf("/api/v1/apps/%v/credentials/keys/%v", appId, keyId)

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var jsonWebKey *JsonWebKey

	resp, err := m.client.requestExecutor.Do(ctx, req, &jsonWebKey)
	if err != nil {
		return nil, resp, err
	}

	return jsonWebKey, resp, nil
}

// Clones a X.509 certificate for an application key credential from a source application to target application.
func (m *ApplicationResource) CloneApplicationKey(ctx context.Context, appId string, keyId string, qp *query.Params) (*JsonWebKey, *Response, error) {
	url := fmt.Sprintf("/api/v1/apps/%v/credentials/keys/%v/clone", appId, keyId)
	if qp != nil {
		url = url + qp.String()
	}

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("POST", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var jsonWebKey *JsonWebKey

	resp, err := m.client.requestExecutor.Do(ctx, req, &jsonWebKey)
	if err != nil {
		return nil, resp, err
	}

	return jsonWebKey, resp, nil
}

// Lists all scope consent grants for the application
func (m *ApplicationResource) ListScopeConsentGrants(ctx context.Context, appId string, qp *query.Params) ([]*OAuth2ScopeConsentGrant, *Response, error) {
	url := fmt.Sprintf("/api/v1/apps/%v/grants", appId)
	if qp != nil {
		url = url + qp.String()
	}

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var oAuth2ScopeConsentGrant []*OAuth2ScopeConsentGrant

	resp, err := m.client.requestExecutor.Do(ctx, req, &oAuth2ScopeConsentGrant)
	if err != nil {
		return nil, resp, err
	}

	return oAuth2ScopeConsentGrant, resp, nil
}

// Grants consent for the application to request an OAuth 2.0 Okta scope
func (m *ApplicationResource) GrantConsentToScope(ctx context.Context, appId string, body OAuth2ScopeConsentGrant) (*OAuth2ScopeConsentGrant, *Response, error) {
	url := fmt.Sprintf("/api/v1/apps/%v/grants", appId)

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	var oAuth2ScopeConsentGrant *OAuth2ScopeConsentGrant

	resp, err := m.client.requestExecutor.Do(ctx, req, &oAuth2ScopeConsentGrant)
	if err != nil {
		return nil, resp, err
	}

	return oAuth2ScopeConsentGrant, resp, nil
}

// Revokes permission for the application to request the given scope
func (m *ApplicationResource) RevokeScopeConsentGrant(ctx context.Context, appId string, grantId string) (*Response, error) {
	url := fmt.Sprintf("/api/v1/apps/%v/grants/%v", appId, grantId)

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

// Fetches a single scope consent grant for the application
func (m *ApplicationResource) GetScopeConsentGrant(ctx context.Context, appId string, grantId string, qp *query.Params) (*OAuth2ScopeConsentGrant, *Response, error) {
	url := fmt.Sprintf("/api/v1/apps/%v/grants/%v", appId, grantId)
	if qp != nil {
		url = url + qp.String()
	}

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var oAuth2ScopeConsentGrant *OAuth2ScopeConsentGrant

	resp, err := m.client.requestExecutor.Do(ctx, req, &oAuth2ScopeConsentGrant)
	if err != nil {
		return nil, resp, err
	}

	return oAuth2ScopeConsentGrant, resp, nil
}

// Enumerates group assignments for an application.
func (m *ApplicationResource) ListApplicationGroupAssignments(ctx context.Context, appId string, qp *query.Params) ([]*ApplicationGroupAssignment, *Response, error) {
	url := fmt.Sprintf("/api/v1/apps/%v/groups", appId)
	if qp != nil {
		url = url + qp.String()
	}

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var applicationGroupAssignment []*ApplicationGroupAssignment

	resp, err := m.client.requestExecutor.Do(ctx, req, &applicationGroupAssignment)
	if err != nil {
		return nil, resp, err
	}

	return applicationGroupAssignment, resp, nil
}

// Removes a group assignment from an application.
func (m *ApplicationResource) DeleteApplicationGroupAssignment(ctx context.Context, appId string, groupId string) (*Response, error) {
	url := fmt.Sprintf("/api/v1/apps/%v/groups/%v", appId, groupId)

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

// Fetches an application group assignment
func (m *ApplicationResource) GetApplicationGroupAssignment(ctx context.Context, appId string, groupId string, qp *query.Params) (*ApplicationGroupAssignment, *Response, error) {
	url := fmt.Sprintf("/api/v1/apps/%v/groups/%v", appId, groupId)
	if qp != nil {
		url = url + qp.String()
	}

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var applicationGroupAssignment *ApplicationGroupAssignment

	resp, err := m.client.requestExecutor.Do(ctx, req, &applicationGroupAssignment)
	if err != nil {
		return nil, resp, err
	}

	return applicationGroupAssignment, resp, nil
}

// Assigns a group to an application
func (m *ApplicationResource) CreateApplicationGroupAssignment(ctx context.Context, appId string, groupId string, body ApplicationGroupAssignment) (*ApplicationGroupAssignment, *Response, error) {
	url := fmt.Sprintf("/api/v1/apps/%v/groups/%v", appId, groupId)

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("PUT", url, body)
	if err != nil {
		return nil, nil, err
	}

	var applicationGroupAssignment *ApplicationGroupAssignment

	resp, err := m.client.requestExecutor.Do(ctx, req, &applicationGroupAssignment)
	if err != nil {
		return nil, resp, err
	}

	return applicationGroupAssignment, resp, nil
}

// Activates an inactive application.
func (m *ApplicationResource) ActivateApplication(ctx context.Context, appId string) (*Response, error) {
	url := fmt.Sprintf("/api/v1/apps/%v/lifecycle/activate", appId)

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := m.client.requestExecutor.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Deactivates an active application.
func (m *ApplicationResource) DeactivateApplication(ctx context.Context, appId string) (*Response, error) {
	url := fmt.Sprintf("/api/v1/apps/%v/lifecycle/deactivate", appId)

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := m.client.requestExecutor.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// Revokes all tokens for the specified application
func (m *ApplicationResource) RevokeOAuth2TokensForApplication(ctx context.Context, appId string) (*Response, error) {
	url := fmt.Sprintf("/api/v1/apps/%v/tokens", appId)

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

// Lists all tokens for the application
func (m *ApplicationResource) ListOAuth2TokensForApplication(ctx context.Context, appId string, qp *query.Params) ([]*OAuth2Token, *Response, error) {
	url := fmt.Sprintf("/api/v1/apps/%v/tokens", appId)
	if qp != nil {
		url = url + qp.String()
	}

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var oAuth2Token []*OAuth2Token

	resp, err := m.client.requestExecutor.Do(ctx, req, &oAuth2Token)
	if err != nil {
		return nil, resp, err
	}

	return oAuth2Token, resp, nil
}

// Revokes the specified token for the specified application
func (m *ApplicationResource) RevokeOAuth2TokenForApplication(ctx context.Context, appId string, tokenId string) (*Response, error) {
	url := fmt.Sprintf("/api/v1/apps/%v/tokens/%v", appId, tokenId)

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

// Gets a token for the specified application
func (m *ApplicationResource) GetOAuth2TokenForApplication(ctx context.Context, appId string, tokenId string, qp *query.Params) (*OAuth2Token, *Response, error) {
	url := fmt.Sprintf("/api/v1/apps/%v/tokens/%v", appId, tokenId)
	if qp != nil {
		url = url + qp.String()
	}

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var oAuth2Token *OAuth2Token

	resp, err := m.client.requestExecutor.Do(ctx, req, &oAuth2Token)
	if err != nil {
		return nil, resp, err
	}

	return oAuth2Token, resp, nil
}

// Enumerates all assigned [application users](#application-user-model) for an application.
func (m *ApplicationResource) ListApplicationUsers(ctx context.Context, appId string, qp *query.Params) ([]*AppUser, *Response, error) {
	url := fmt.Sprintf("/api/v1/apps/%v/users", appId)
	if qp != nil {
		url = url + qp.String()
	}

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var appUser []*AppUser

	resp, err := m.client.requestExecutor.Do(ctx, req, &appUser)
	if err != nil {
		return nil, resp, err
	}

	return appUser, resp, nil
}

// Assigns an user to an application with [credentials](#application-user-credentials-object) and an app-specific [profile](#application-user-profile-object). Profile mappings defined for the application are first applied before applying any profile properties specified in the request.
func (m *ApplicationResource) AssignUserToApplication(ctx context.Context, appId string, body AppUser) (*AppUser, *Response, error) {
	url := fmt.Sprintf("/api/v1/apps/%v/users", appId)

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	var appUser *AppUser

	resp, err := m.client.requestExecutor.Do(ctx, req, &appUser)
	if err != nil {
		return nil, resp, err
	}

	return appUser, resp, nil
}

// Removes an assignment for a user from an application.
func (m *ApplicationResource) DeleteApplicationUser(ctx context.Context, appId string, userId string, qp *query.Params) (*Response, error) {
	url := fmt.Sprintf("/api/v1/apps/%v/users/%v", appId, userId)
	if qp != nil {
		url = url + qp.String()
	}

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

// Fetches a specific user assignment for application by &#x60;id&#x60;.
func (m *ApplicationResource) GetApplicationUser(ctx context.Context, appId string, userId string, qp *query.Params) (*AppUser, *Response, error) {
	url := fmt.Sprintf("/api/v1/apps/%v/users/%v", appId, userId)
	if qp != nil {
		url = url + qp.String()
	}

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var appUser *AppUser

	resp, err := m.client.requestExecutor.Do(ctx, req, &appUser)
	if err != nil {
		return nil, resp, err
	}

	return appUser, resp, nil
}

// Updates a user&#x27;s profile for an application
func (m *ApplicationResource) UpdateApplicationUser(ctx context.Context, appId string, userId string, body AppUser) (*AppUser, *Response, error) {
	url := fmt.Sprintf("/api/v1/apps/%v/users/%v", appId, userId)

	req, err := m.client.requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("POST", url, body)
	if err != nil {
		return nil, nil, err
	}

	var appUser *AppUser

	resp, err := m.client.requestExecutor.Do(ctx, req, &appUser)
	if err != nil {
		return nil, resp, err
	}

	return appUser, resp, nil
}
