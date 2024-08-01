package xelon

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
)

const (
	libraryVersion = "0.14.1"

	defaultBaseURL   = "https://hq.xelon.ch/api/service/"
	defaultMediaType = "application/json"
	defaultUserAgent = "xelon-sdk-go/" + libraryVersion

	headerStackifyID = "X-StackifyID"
)

// A Client manages communication with the Xelon API.
type Client struct {
	httpClient *http.Client // HTTP client used to communicate with the API.

	BaseURL   *url.URL // Base URL for API requests. BaseURL should always be specified with a trailing slash.
	ClientID  string   // ClientID for IP ranges.
	token     string   // token for Xelon API.
	UserAgent string   // User agent used when communicating with Xelon API.

	common service // Reuse a single struct instead of allocating one for each service on the heap.

	Clouds               *CloudsService
	Devices              *DevicesService
	LoadBalancerClusters *LoadBalancerClustersService
	LoadBalancers        *LoadBalancersService
	Networks             *NetworksService
	PersistentStorages   *PersistentStoragesService
	SSHKeys              *SSHKeysService
	Templates            *TemplatesService
	Tenants              *TenantsService
}

type service struct {
	client *Client
}

// ListOptions specifies the optional parameters to various List methods that
// support pagination.
type ListOptions struct {
	// Page of results to retrieve.
	Page int `url:"page,omitempty"`

	// PerPage specifies the number of results to include per page.
	PerPage int `url:"per_page,omitempty"`
}

// addOptions adds the parameters in opts as URL query parameters to s. opts
// must be a struct whose fields may contain "url" tags.
func addOptions(s string, opts interface{}) (string, error) {
	v := reflect.ValueOf(opts)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opts)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}

type ClientOption func(*Client)

// WithBaseURL configures Client to use a specific API endpoint.
func WithBaseURL(baseURL string) ClientOption {
	return func(client *Client) {
		parsedURL, _ := url.Parse(baseURL)
		client.BaseURL = parsedURL
	}
}

// WithClientID configures Client to use "X-User-Id" http header by all API requests.
func WithClientID(clientID string) ClientOption {
	return func(client *Client) {
		client.ClientID = clientID
	}
}

// WithHTTPClient configures Client to use a specific http client for communication.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(client *Client) {
		client.httpClient = httpClient
	}
}

// WithUserAgent configures Client to use a specific user agent.
func WithUserAgent(ua string) ClientOption {
	return func(client *Client) {
		client.UserAgent = ua
	}
}

// NewClient returns a new Xelon API client.
func NewClient(token string, opts ...ClientOption) *Client {
	baseUrl, _ := url.Parse(defaultBaseURL)
	httpClient := &http.Client{
		Timeout: 60 * time.Second,
	}

	c := &Client{
		BaseURL:    baseUrl,
		httpClient: httpClient,
		token:      token,
		UserAgent:  defaultUserAgent,
	}
	for _, opt := range opts {
		opt(c)
	}

	c.common.client = c

	c.Clouds = (*CloudsService)(&c.common)
	c.Devices = (*DevicesService)(&c.common)
	c.LoadBalancerClusters = (*LoadBalancerClustersService)(&c.common)
	c.LoadBalancers = (*LoadBalancersService)(&c.common)
	c.Networks = (*NetworksService)(&c.common)
	c.PersistentStorages = (*PersistentStoragesService)(&c.common)
	c.SSHKeys = (*SSHKeysService)(&c.common)
	c.Templates = (*TemplatesService)(&c.common)
	c.Tenants = (*TenantsService)(&c.common)

	// Notify user if no ClientID is set
	if c.ClientID == "" {
		fmt.Printf("ClientID is not set, please update your credentials\nUsing the HQ-API without the ClientID-Header will be deprecated in 2024\n")
	}

	return c
}

// Deprecated: SetBaseURL overrides the default BaseURL. Use WithBaseURL instead.
func (c *Client) SetBaseURL(baseURL string) {
	parsedURL, _ := url.Parse(baseURL)
	c.BaseURL = parsedURL
}

// Deprecated: SetUserAgent overrides the default UserAgent. Use WithUserAgent instead.
func (c *Client) SetUserAgent(ua string) {
	c.UserAgent = ua
}

// NewRequest creates an API request. A relative URL can be provided in urlStr, in which case it is resolved
// relative to the BaseURL of the Client. Relative URLs should always be specified without a preceding slash.
// If specified, the value pointed to by body is JSON encoded and included as the request body.
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	if !strings.HasSuffix(c.BaseURL.Path, "/") {
		return nil, fmt.Errorf("BaseURL must have a traling slash, but %q does not", c.BaseURL)
	}
	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if body != nil {
		err = json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if req.Header.Get("Authorization") == "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	}
	req.Header.Set("Accept", defaultMediaType)
	req.Header.Set("Content-Type", defaultMediaType)
	req.Header.Set("User-Agent", c.UserAgent)

	if c.ClientID != "" {
		req.Header.Set("X-User-Id", c.ClientID)
	}

	return req, nil
}

// Response is a Xelon response. This wraps the standard http.Response.
type Response struct {
	*http.Response

	Meta *Meta // Meta describes generic information about the response.

	StackifyID string // StackifyID returned from the API, useful to contact support.
}

// Do sends an API request and returns the API response. The API response is JSON decoded and stored in
// the value pointed to by v, or returned as an error if an API error has occurred.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
	req = req.WithContext(ctx)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		// if we got an error, and the context has been canceled, the context's error is more useful.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		// if the error type is *url.Error, sanitize its URL before returning.
		if e, ok := err.(*url.Error); ok {
			if uri, err := url.Parse(e.URL); err == nil {
				e.URL = sanitizeURL(uri).String()
				return nil, e
			}
		}

		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	response := newResponse(resp)
	err = CheckResponse(response)
	if err != nil {
		return response, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
			if err != nil {
				return nil, err
			}
		} else {
			decodedErr := json.NewDecoder(resp.Body).Decode(v)
			if decodedErr == io.EOF {
				// ignore EOF errors caused by empty response body
				decodedErr = nil
			}
			if decodedErr != nil {
				err = decodedErr
			}
		}
	}

	return response, err
}

// newResponse creates a new Response for the provided http.Response. r must be not nil.
func newResponse(r *http.Response) *Response {
	response := &Response{Response: r}
	response.populateStackifyID()
	return response
}

// populateStackifyID parses the request headers and populates the response stackify id.
func (r *Response) populateStackifyID() {
	if stackifyID := r.Header.Get(headerStackifyID); stackifyID != "" {
		r.StackifyID = stackifyID
	}
}

// CheckResponse checks the API response for errors, and returns them if present. A response is considered
// an error if it has a status code outside the 200 range.
func CheckResponse(resp *Response) error {
	if code := resp.StatusCode; code >= 200 && code <= 299 {
		return nil
	}

	errorResponse := &ErrorResponse{Response: resp}
	data, err := io.ReadAll(resp.Body)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, &errorResponse.ErrorElement)
		if err != nil {
			return err
		}
	}
	return errorResponse
}

// sanitizeURL redacts the password parameter from the URL which may be exposed by the user.
func sanitizeURL(uri *url.URL) *url.URL {
	if uri == nil {
		return nil
	}
	params := uri.Query()
	if len(params.Get("password")) > 0 {
		params.Set("password", "REDACTED")
		uri.RawQuery = params.Encode()
	}

	return uri
}
