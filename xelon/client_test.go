package xelon

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	mux    *http.ServeMux
	ctx    = context.TODO()
	client *Client
	server *httptest.Server
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	client = NewClient("auth-token")
	client.baseURL, _ = url.Parse(fmt.Sprintf("%v/", server.URL))
}

func teardown() {
	server.Close()
}

func TestClient_NewClient(t *testing.T) {
	setup()
	defer teardown()

	assert.NotNil(t, client.baseURL)
	assert.Equal(t, fmt.Sprintf("%v/", server.URL), client.baseURL.String())
	assert.Equal(t, "auth-token", client.token)
}

func TestClient_SetUserAgent(t *testing.T) {
	client := &Client{}

	client.SetUserAgent("custom-user-agent")

	assert.Equal(t, "custom-user-agent", client.userAgent)
}

func TestClient_Defaults(t *testing.T) {
	client := NewClient("auth-token")

	assert.Equal(t, "https://hq.xelon.ch/api/service/", client.baseURL.String())
	assert.Contains(t, client.userAgent, "xelon-sdk-go/")
	assert.Equal(t, 60*time.Second, client.httpClient.Timeout)
}

func TestClient_WithBaseURL(t *testing.T) {
	expectedBaseURL, _ := url.Parse("https://testing.xelon.ch/")
	client := NewClient("auth-token",
		WithBaseURL("https://testing.xelon.ch/"),
	)

	assert.Equal(t, expectedBaseURL, client.baseURL)
}

func TestClient_WithClientID(t *testing.T) {
	client := NewClient("auth-token",
		WithClientID("custom-client-id"),
	)

	assert.Equal(t, "custom-client-id", client.clientID)
}

func TestClient_WithHTTPClient(t *testing.T) {
	httpClient := &http.Client{Timeout: 2 * time.Second}
	client := NewClient("auth-token",
		WithHTTPClient(httpClient),
	)

	assert.Equal(t, httpClient, client.httpClient)
	assert.Equal(t, 2*time.Second, client.httpClient.Timeout)
}

func TestClient_WithUserAgent(t *testing.T) {
	client := NewClient("auth-token",
		WithUserAgent("custom-user-agent"),
	)

	assert.Equal(t, "custom-user-agent", client.userAgent)
}
