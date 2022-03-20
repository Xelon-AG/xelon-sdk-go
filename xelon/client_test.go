package xelon

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

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
	client.BaseURL, _ = url.Parse(fmt.Sprintf("%v/", server.URL))
}

func teardown() {
	server.Close()
}

func TestClient_NewClient(t *testing.T) {
	setup()
	defer teardown()

	assert.NotNil(t, client.BaseURL)
	assert.Equal(t, fmt.Sprintf("%v/", server.URL), client.BaseURL.String())
	assert.Equal(t, "auth-token", client.Token)
}

func TestClient_SetUserAgent(t *testing.T) {
	setup()
	defer teardown()

	client.SetUserAgent("custom-user-agent")

	assert.Equal(t, "custom-user-agent", client.UserAgent)
}
