package tfe

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/hashicorp/go-cleanhttp"
	retryablehttp "github.com/hashicorp/go-retryablehttp"
	"github.com/svanharmelen/jsonapi"
	"golang.org/x/time/rate"
)

const (
	userAgent        = "go-tfe"
	headerRateLimit  = "X-RateLimit-Limit"
	headerRateReset  = "X-RateLimit-Reset"
	headerAPIVersion = "TFP-API-Version"

	// DefaultAddress of Terraform Enterprise.
	DefaultAddress = "https://app.terraform.io"
	// DefaultBasePath on which the API is served.
	DefaultBasePath = "/api/v2/"
	// PingEndpoint is a no-op API endpoint used to configure the rate limiter
	PingEndpoint = "ping"
)

var (
	// ErrWorkspaceLocked is returned when trying to lock a
	// locked workspace.
	ErrWorkspaceLocked = errors.New("workspace already locked")
	// ErrWorkspaceNotLocked is returned when trying to unlock
	// a unlocked workspace.
	ErrWorkspaceNotLocked = errors.New("workspace already unlocked")

	// ErrUnauthorized is returned when a receiving a 401.
	ErrUnauthorized = errors.New("unauthorized")
	// ErrResourceNotFound is returned when a receiving a 404.
	ErrResourceNotFound = errors.New("resource not found")
)

// RetryLogHook allows a function to run before each retry.
type RetryLogHook func(attemptNum int, resp *http.Response)

// Config provides configuration details to the API client.
type Config struct {
	// The address of the Terraform Enterprise API.
	Address string

	// The base path on which the API is served.
	BasePath string

	// API token used to access the Terraform Enterprise API.
	Token string

	// Headers that will be added to every request.
	Headers http.Header

	// A custom HTTP client to use.
	HTTPClient *http.Client

	// RetryLogHook is invoked each time a request is retried.
	RetryLogHook RetryLogHook
}

// DefaultConfig returns a default config structure.
func DefaultConfig() *Config {
	config := &Config{
		Address:    os.Getenv("TFE_ADDRESS"),
		BasePath:   DefaultBasePath,
		Token:      os.Getenv("TFE_TOKEN"),
		Headers:    make(http.Header),
		HTTPClient: cleanhttp.DefaultPooledClient(),
	}

	// Set the default address if none is given.
	if config.Address == "" {
		config.Address = DefaultAddress
	}

	// Set the default user agent.
	config.Headers.Set("User-Agent", userAgent)

	return config
}

// Client is the Terraform Enterprise API client. It provides the basic
// connectivity and configuration for accessing the TFE API.
type Client struct {
	baseURL           *url.URL
	token             string
	headers           http.Header
	http              *retryablehttp.Client
	limiter           *rate.Limiter
	retryLogHook      RetryLogHook
	retryServerErrors bool
	remoteAPIVersion  string

	AgentPools                 AgentPools
	AgentTokens                AgentTokens
	Applies                    Applies
	ConfigurationVersions      ConfigurationVersions
	CostEstimates              CostEstimates
	NotificationConfigurations NotificationConfigurations
	OAuthClients               OAuthClients
	OAuthTokens                OAuthTokens
	Organizations              Organizations
	OrganizationMemberships    OrganizationMemberships
	OrganizationTokens         OrganizationTokens
	Plans                      Plans
	PlanExports                PlanExports
	Policies                   Policies
	PolicyChecks               PolicyChecks
	PolicySetParameters        PolicySetParameters
	PolicySets                 PolicySets
	RegistryModules            RegistryModules
	Runs                       Runs
	RunTriggers                RunTriggers
	SSHKeys                    SSHKeys
	StateVersionOutputs        StateVersionOutputs
	StateVersions              StateVersions
	Teams                      Teams
	TeamAccess                 TeamAccesses
	TeamMembers                TeamMembers
	TeamTokens                 TeamTokens
	Users                      Users
	UserTokens                 UserTokens
	Variables                  Variables
	Workspaces                 Workspaces

	Meta Meta
}

// Meta contains any Terraform Cloud APIs which provide data about the API itself.
type Meta struct {
	IPRanges IPRanges
}

// NewClient creates a new Terraform Enterprise API client.
func NewClient(cfg *Config) (*Client, error) {
	config := DefaultConfig()

	// Layer in the provided config for any non-blank values.
	if cfg != nil {
		if cfg.Address != "" {
			config.Address = cfg.Address
		}
		if cfg.BasePath != "" {
			config.BasePath = cfg.BasePath
		}
		if cfg.Token != "" {
			config.Token = cfg.Token
		}
		for k, v := range cfg.Headers {
			config.Headers[k] = v
		}
		if cfg.HTTPClient != nil {
			config.HTTPClient = cfg.HTTPClient
		}
		if cfg.RetryLogHook != nil {
			config.RetryLogHook = cfg.RetryLogHook
		}
	}

	// Parse the address to make sure its a valid URL.
	baseURL, err := url.Parse(config.Address)
	if err != nil {
		return nil, fmt.Errorf("invalid address: %v", err)
	}

	baseURL.Path = config.BasePath
	if !strings.HasSuffix(baseURL.Path, "/") {
		baseURL.Path += "/"
	}

	// This value must be provided by the user.
	if config.Token == "" {
		return nil, fmt.Errorf("missing API token")
	}

	// Create the client.
	client := &Client{
		baseURL:      baseURL,
		token:        config.Token,
		headers:      config.Headers,
		retryLogHook: config.RetryLogHook,
	}

	client.http = &retryablehttp.Client{
		Backoff:      client.retryHTTPBackoff,
		CheckRetry:   client.retryHTTPCheck,
		ErrorHandler: retryablehttp.PassthroughErrorHandler,
		HTTPClient:   config.HTTPClient,
		RetryWaitMin: 100 * time.Millisecond,
		RetryWaitMax: 400 * time.Millisecond,
		RetryMax:     30,
	}

	meta, err := client.getRawAPIMetadata()
	if err != nil {
		return nil, err
	}

	// Configure the rate limiter.
	client.configureLimiter(meta.RateLimit)

	// Save the API version so we can return it from the RemoteAPIVersion
	// method later.
	client.remoteAPIVersion = meta.APIVersion

	// Create the services.
	client.AgentPools = &agentPools{client: client}
	client.AgentTokens = &agentTokens{client: client}
	client.Applies = &applies{client: client}
	client.ConfigurationVersions = &configurationVersions{client: client}
	client.CostEstimates = &costEstimates{client: client}
	client.NotificationConfigurations = &notificationConfigurations{client: client}
	client.OAuthClients = &oAuthClients{client: client}
	client.OAuthTokens = &oAuthTokens{client: client}
	client.Organizations = &organizations{client: client}
	client.OrganizationMemberships = &organizationMemberships{client: client}
	client.OrganizationTokens = &organizationTokens{client: client}
	client.Plans = &plans{client: client}
	client.PlanExports = &planExports{client: client}
	client.Policies = &policies{client: client}
	client.PolicyChecks = &policyChecks{client: client}
	client.PolicySetParameters = &policySetParameters{client: client}
	client.PolicySets = &policySets{client: client}
	client.RegistryModules = &registryModules{client: client}
	client.Runs = &runs{client: client}
	client.RunTriggers = &runTriggers{client: client}
	client.SSHKeys = &sshKeys{client: client}
	client.StateVersionOutputs = &stateVersionOutputs{client: client}
	client.StateVersions = &stateVersions{client: client}
	client.Teams = &teams{client: client}
	client.TeamAccess = &teamAccesses{client: client}
	client.TeamMembers = &teamMembers{client: client}
	client.TeamTokens = &teamTokens{client: client}
	client.Users = &users{client: client}
	client.UserTokens = &userTokens{client: client}
	client.Variables = &variables{client: client}
	client.Workspaces = &workspaces{client: client}

	client.Meta = Meta{
		IPRanges: &ipRanges{client: client},
	}

	return client, nil
}

// RemoteAPIVersion returns the server's declared API version string.
//
// A Terraform Cloud or Enterprise API server returns its API version in an
// HTTP header field in all responses. The NewClient function saves the
// version number returned in its initial setup request and RemoteAPIVersion
// returns that cached value.
//
// The API protocol calls for this string to be a dotted-decimal version number
// like 2.3.0, where the first number indicates the API major version while the
// second indicates a minor version which may have introduced some
// backward-compatible additional features compared to its predecessor.
//
// Explicit API versioning was added to the Terraform Cloud and Enterprise
// APIs as a later addition, so older servers will not return version
// information. In that case, this function returns an empty string as the
// version.
func (c *Client) RemoteAPIVersion() string {
	return c.remoteAPIVersion
}

// SetFakeRemoteAPIVersion allows setting a given string as the client's remoteAPIVersion,
// overriding the value pulled from the API header during client initialization.
//
// This is intended for use in tests, when you may want to configure your TFE client to
// return something different than the actual API version in order to test error handling.
func (c *Client) SetFakeRemoteAPIVersion(fakeAPIVersion string) {
	c.remoteAPIVersion = fakeAPIVersion
}

// RetryServerErrors configures the retry HTTP check to also retry
// unexpected errors or requests that failed with a server error.
func (c *Client) RetryServerErrors(retry bool) {
	c.retryServerErrors = retry
}

// retryHTTPCheck provides a callback for Client.CheckRetry which
// will retry both rate limit (429) and server (>= 500) errors.
func (c *Client) retryHTTPCheck(ctx context.Context, resp *http.Response, err error) (bool, error) {
	if ctx.Err() != nil {
		return false, ctx.Err()
	}
	if err != nil {
		return c.retryServerErrors, err
	}
	if resp.StatusCode == 429 || (c.retryServerErrors && resp.StatusCode >= 500) {
		return true, nil
	}
	return false, nil
}

// retryHTTPBackoff provides a generic callback for Client.Backoff which
// will pass through all calls based on the status code of the response.
func (c *Client) retryHTTPBackoff(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
	if c.retryLogHook != nil {
		c.retryLogHook(attemptNum, resp)
	}

	// Use the rate limit backoff function when we are rate limited.
	if resp != nil && resp.StatusCode == 429 {
		return rateLimitBackoff(min, max, attemptNum, resp)
	}

	// Set custom duration's when we experience a service interruption.
	min = 700 * time.Millisecond
	max = 900 * time.Millisecond

	return retryablehttp.LinearJitterBackoff(min, max, attemptNum, resp)
}

// rateLimitBackoff provides a callback for Client.Backoff which will use the
// X-RateLimit_Reset header to determine the time to wait. We add some jitter
// to prevent a thundering herd.
//
// min and max are mainly used for bounding the jitter that will be added to
// the reset time retrieved from the headers. But if the final wait time is
// less then min, min will be used instead.
func rateLimitBackoff(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
	// rnd is used to generate pseudo-random numbers.
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	// First create some jitter bounded by the min and max durations.
	jitter := time.Duration(rnd.Float64() * float64(max-min))

	if resp != nil {
		if v := resp.Header.Get(headerRateReset); v != "" {
			if reset, _ := strconv.ParseFloat(v, 64); reset > 0 {
				// Only update min if the given time to wait is longer.
				if wait := time.Duration(reset * 1e9); wait > min {
					min = wait
				}
			}
		}
	}

	return min + jitter
}

type rawAPIMetadata struct {
	// APIVersion is the raw API version string reported by the server in the
	// TFP-API-Version response header, or an empty string if that header
	// field was not included in the response.
	APIVersion string

	// RateLimit is the raw API version string reported by the server in the
	// X-RateLimit-Limit response header, or an empty string if that header
	// field was not included in the response.
	RateLimit string
}

func (c *Client) getRawAPIMetadata() (rawAPIMetadata, error) {
	var meta rawAPIMetadata

	// Create a new request.
	u, err := c.baseURL.Parse(PingEndpoint)
	if err != nil {
		return meta, err
	}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return meta, err
	}

	// Attach the default headers.
	for k, v := range c.headers {
		req.Header[k] = v
	}
	req.Header.Set("Accept", "application/vnd.api+json")
	req.Header.Set("Authorization", "Bearer "+c.token)

	// Make a single request to retrieve the rate limit headers.
	resp, err := c.http.HTTPClient.Do(req)
	if err != nil {
		return meta, err
	}
	resp.Body.Close()

	meta.APIVersion = resp.Header.Get(headerAPIVersion)
	meta.RateLimit = resp.Header.Get(headerRateLimit)

	return meta, nil
}

// configureLimiter configures the rate limiter.
func (c *Client) configureLimiter(rawLimit string) {

	// Set default values for when rate limiting is disabled.
	limit := rate.Inf
	burst := 0

	if v := rawLimit; v != "" {
		if rateLimit, _ := strconv.ParseFloat(v, 64); rateLimit > 0 {
			// Configure the limit and burst using a split of 2/3 for the limit and
			// 1/3 for the burst. This enables clients to burst 1/3 of the allowed
			// calls before the limiter kicks in. The remaining calls will then be
			// spread out evenly using intervals of time.Second / limit which should
			// prevent hitting the rate limit.
			limit = rate.Limit(rateLimit * 0.66)
			burst = int(rateLimit * 0.33)
		}
	}

	// Create a new limiter using the calculated values.
	c.limiter = rate.NewLimiter(limit, burst)
}

// newRequest creates an API request. A relative URL path can be provided in
// path, in which case it is resolved relative to the apiVersionPath of the
// Client. Relative URL paths should always be specified without a preceding
// slash.
// If v is supplied, the value will be JSONAPI encoded and included as the
// request body. If the method is GET, the value will be parsed and added as
// query parameters.
func (c *Client) newRequest(method, path string, v interface{}) (*retryablehttp.Request, error) {
	u, err := c.baseURL.Parse(path)
	if err != nil {
		return nil, err
	}

	// Create a request specific headers map.
	reqHeaders := make(http.Header)
	reqHeaders.Set("Authorization", "Bearer "+c.token)

	var body interface{}
	switch method {
	case "GET":
		reqHeaders.Set("Accept", "application/vnd.api+json")

		if v != nil {
			q, err := query.Values(v)
			if err != nil {
				return nil, err
			}
			u.RawQuery = q.Encode()
		}
	case "DELETE", "PATCH", "POST":
		reqHeaders.Set("Accept", "application/vnd.api+json")
		reqHeaders.Set("Content-Type", "application/vnd.api+json")

		if v != nil {
			if body, err = serializeRequestBody(v); err != nil {
				return nil, err
			}
		}
	case "PUT":
		reqHeaders.Set("Accept", "application/json")
		reqHeaders.Set("Content-Type", "application/octet-stream")
		body = v
	}

	req, err := retryablehttp.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	// Set the default headers.
	for k, v := range c.headers {
		req.Header[k] = v
	}

	// Set the request specific headers.
	for k, v := range reqHeaders {
		req.Header[k] = v
	}

	return req, nil
}

// Helper method that serializes the given ptr or ptr slice into a JSON
// request. It automatically uses jsonapi or json serialization, depending
// on the body type's tags.
func serializeRequestBody(v interface{}) (interface{}, error) {
	// The body can be a slice of pointers or a pointer. In either
	// case we want to choose the serialization type based on the
	// individual record type. To determine that type, we need
	// to either follow the pointer or examine the slice element type.
	// There are other theoretical possiblities (e. g. maps,
	// non-pointers) but they wouldn't work anyway because the
	// json-api library doesn't support serializing other things.
	var modelType reflect.Type
	bodyType := reflect.TypeOf(v)
	invalidBodyError := errors.New("go-tfe bug: DELETE/PATCH/POST body must be nil, ptr, or ptr slice")
	switch bodyType.Kind() {
	case reflect.Slice:
		sliceElem := bodyType.Elem()
		if sliceElem.Kind() != reflect.Ptr {
			return nil, invalidBodyError
		}
		modelType = sliceElem.Elem()
	case reflect.Ptr:
		modelType = reflect.ValueOf(v).Elem().Type()
	default:
		return nil, invalidBodyError
	}

	// Infer whether the request uses jsonapi or regular json
	// serialization based on how the fields are tagged.
	jsonApiFields := 0
	jsonFields := 0
	for i := 0; i < modelType.NumField(); i++ {
		structField := modelType.Field(i)
		if structField.Tag.Get("jsonapi") != "" {
			jsonApiFields++
		}
		if structField.Tag.Get("json") != "" {
			jsonFields++
		}
	}
	if jsonApiFields > 0 && jsonFields > 0 {
		// Defining a struct with both json and jsonapi tags doesn't
		// make sense, because a struct can only be serialized
		// as one or another. If this does happen, it's a bug
		// in the library that should be fixed at development time
		return nil, errors.New("go-tfe bug: struct can't use both json and jsonapi attributes")
	}

	if jsonFields > 0 {
		return json.Marshal(v)
	} else {
		buf := bytes.NewBuffer(nil)
		if err := jsonapi.MarshalPayloadWithoutIncluded(buf, v); err != nil {
			return nil, err
		}
		return buf, nil
	}
}

// do sends an API request and returns the API response. The API response
// is JSONAPI decoded and the document's primary data is stored in the value
// pointed to by v, or returned as an error if an API error has occurred.

// If v implements the io.Writer interface, the raw response body will be
// written to v, without attempting to first decode it.
//
// The provided ctx must be non-nil. If it is canceled or times out, ctx.Err()
// will be returned.
func (c *Client) do(ctx context.Context, req *retryablehttp.Request, v interface{}) error {
	// Wait will block until the limiter can obtain a new token
	// or returns an error if the given context is canceled.
	if err := c.limiter.Wait(ctx); err != nil {
		return err
	}

	// Add the context to the request.
	req = req.WithContext(ctx)

	// Execute the request and check the response.
	resp, err := c.http.Do(req)
	if err != nil {
		// If we got an error, and the context has been canceled,
		// the context's error is probably more useful.
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			return err
		}
	}
	defer resp.Body.Close()

	// Basic response checking.
	if err := checkResponseCode(resp); err != nil {
		return err
	}

	// Return here if decoding the response isn't needed.
	if v == nil {
		return nil
	}

	// If v implements io.Writer, write the raw response body.
	if w, ok := v.(io.Writer); ok {
		_, err = io.Copy(w, resp.Body)
		return err
	}

	// Get the value of v so we can test if it's a struct.
	dst := reflect.Indirect(reflect.ValueOf(v))

	// Return an error if v is not a struct or an io.Writer.
	if dst.Kind() != reflect.Struct {
		return fmt.Errorf("v must be a struct or an io.Writer")
	}

	// Try to get the Items and Pagination struct fields.
	items := dst.FieldByName("Items")
	pagination := dst.FieldByName("Pagination")

	// Unmarshal a single value if v does not contain the
	// Items and Pagination struct fields.
	if !items.IsValid() || !pagination.IsValid() {
		return jsonapi.UnmarshalPayload(resp.Body, v)
	}

	// Return an error if v.Items is not a slice.
	if items.Type().Kind() != reflect.Slice {
		return fmt.Errorf("v.Items must be a slice")
	}

	// Create a temporary buffer and copy all the read data into it.
	body := bytes.NewBuffer(nil)
	reader := io.TeeReader(resp.Body, body)

	// Unmarshal as a list of values as v.Items is a slice.
	raw, err := jsonapi.UnmarshalManyPayload(reader, items.Type().Elem())
	if err != nil {
		return err
	}

	// Make a new slice to hold the results.
	sliceType := reflect.SliceOf(items.Type().Elem())
	result := reflect.MakeSlice(sliceType, 0, len(raw))

	// Add all of the results to the new slice.
	for _, v := range raw {
		result = reflect.Append(result, reflect.ValueOf(v))
	}

	// Pointer-swap the result.
	items.Set(result)

	// As we are getting a list of values, we need to decode
	// the pagination details out of the response body.
	p, err := parsePagination(body)
	if err != nil {
		return err
	}

	// Pointer-swap the decoded pagination details.
	pagination.Set(reflect.ValueOf(p))

	return nil
}

// ListOptions is used to specify pagination options when making API requests.
// Pagination allows breaking up large result sets into chunks, or "pages".
type ListOptions struct {
	// The page number to request. The results vary based on the PageSize.
	PageNumber int `url:"page[number],omitempty"`

	// The number of elements returned in a single page.
	PageSize int `url:"page[size],omitempty"`
}

// Pagination is used to return the pagination details of an API request.
type Pagination struct {
	CurrentPage  int `json:"current-page"`
	PreviousPage int `json:"prev-page"`
	NextPage     int `json:"next-page"`
	TotalPages   int `json:"total-pages"`
	TotalCount   int `json:"total-count"`
}

func parsePagination(body io.Reader) (*Pagination, error) {
	var raw struct {
		Meta struct {
			Pagination Pagination `json:"pagination"`
		} `json:"meta"`
	}

	// JSON decode the raw response.
	if err := json.NewDecoder(body).Decode(&raw); err != nil {
		return &Pagination{}, err
	}

	return &raw.Meta.Pagination, nil
}

// checkResponseCode can be used to check the status code of an HTTP request.
func checkResponseCode(r *http.Response) error {
	if r.StatusCode >= 200 && r.StatusCode <= 299 {
		return nil
	}

	switch r.StatusCode {
	case 401:
		return ErrUnauthorized
	case 404:
		return ErrResourceNotFound
	case 409:
		switch {
		case strings.HasSuffix(r.Request.URL.Path, "actions/lock"):
			return ErrWorkspaceLocked
		case strings.HasSuffix(r.Request.URL.Path, "actions/unlock"):
			return ErrWorkspaceNotLocked
		case strings.HasSuffix(r.Request.URL.Path, "actions/force-unlock"):
			return ErrWorkspaceNotLocked
		}
	}

	// Decode the error payload.
	errPayload := &jsonapi.ErrorsPayload{}
	err := json.NewDecoder(r.Body).Decode(errPayload)
	if err != nil || len(errPayload.Errors) == 0 {
		return fmt.Errorf(r.Status)
	}

	// Parse and format the errors.
	var errs []string
	for _, e := range errPayload.Errors {
		if e.Detail == "" {
			errs = append(errs, e.Title)
		} else {
			errs = append(errs, fmt.Sprintf("%s\n\n%s", e.Title, e.Detail))
		}
	}

	return fmt.Errorf(strings.Join(errs, "\n"))
}
