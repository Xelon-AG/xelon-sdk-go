package xelon

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
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

func newTestClient(t testing.TB, options ...ClientOption) *Client {
	t.Helper()

	rejectingHTTPClient := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			err := fmt.Errorf("unexpected HTTP request: %s %s", req.Method, req.URL)
			t.Error(err)
			return nil, err
		}),
	}
	defaults := []ClientOption{
		WithBaseURL("https://example.invalid/"),
		WithHTTPClient(rejectingHTTPClient),
	}
	defaults = append(defaults, options...)

	return NewClient("auth-token", defaults...)
}

func setup(t testing.TB) {
	t.Helper()

	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	client = newTestClient(
		t,
		WithBaseURL(server.URL+"/"),
		WithHTTPClient(server.Client()),
	)
}

func teardown() {
	server.Close()
}

func TestClient_NewClient(t *testing.T) {
	setup(t)
	defer teardown()

	assert.NotNil(t, client.baseURL)
	assert.Equal(t, fmt.Sprintf("%v/", server.URL), client.baseURL.String())
	assert.Equal(t, "auth-token", client.token)
}

func TestClient_Defaults(t *testing.T) {
	client := NewClient("auth-token")

	assert.Equal(t, defaultBaseURL, client.baseURL.String())
	assert.Contains(t, client.userAgent, "xelon-sdk-go/")
	assert.Zero(t, client.httpClient.Timeout)
	assert.False(t, client.customHTTPClient)
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
	assert.True(t, client.customHTTPClient)
}

func TestClient_WithUserAgent(t *testing.T) {
	client := NewClient("auth-token",
		WithUserAgent("custom-user-agent"),
	)

	assert.Equal(t, "custom-user-agent", client.userAgent)
}

func TestClient_NewRequest(t *testing.T) {
	client := NewClient(
		"auth-token",
		WithBaseURL("https://example.test/api/v2/"),
		WithClientID("client-id"),
		WithUserAgent("custom-user-agent"),
	)

	req, err := client.NewRequest(http.MethodPost, "resources", struct {
		Name string `json:"name"`
	}{Name: "example"})
	assert.NoError(t, err)

	body, err := io.ReadAll(req.Body)
	assert.NoError(t, err)
	assert.Equal(t, "{\"name\":\"example\"}\n", string(body))
	assert.Equal(t, int64(len(body)), req.ContentLength)
	assert.Equal(t, "https://example.test/api/v2/resources", req.URL.String())
	assert.Equal(t, "Bearer auth-token", req.Header.Get("Authorization"))
	assert.Equal(t, defaultMediaType, req.Header.Get("Accept"))
	assert.Equal(t, defaultMediaType, req.Header.Get("Content-Type"))
	assert.Equal(t, "custom-user-agent", req.Header.Get("User-Agent"))
	assert.Equal(t, "client-id", req.Header.Get("X-User-Id"))
}

func TestClient_NewRequestWithContentType(t *testing.T) {
	client := NewClient(
		"auth-token",
		WithBaseURL("https://example.test/api/v2/"),
		WithClientID("client-id"),
		WithUserAgent("custom-user-agent"),
	)
	contentType := "multipart/form-data; boundary=test-boundary"

	req, err := client.newRequest(
		http.MethodPost,
		"isos/upload",
		bytes.NewReader([]byte("multipart body")),
		contentType,
	)
	assert.NoError(t, err)

	assert.Equal(t, "Bearer auth-token", req.Header.Get("Authorization"))
	assert.Equal(t, defaultMediaType, req.Header.Get("Accept"))
	assert.Equal(t, contentType, req.Header.Get("Content-Type"))
	assert.Equal(t, "custom-user-agent", req.Header.Get("User-Agent"))
	assert.Equal(t, "client-id", req.Header.Get("X-User-Id"))
}

func TestClient_DoTimeoutPolicy(t *testing.T) {
	t.Run("SDK-owned client uses ordinary fallback", func(t *testing.T) {
		client := NewClient("auth-token")
		client.httpClient.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
			deadline, ok := req.Context().Deadline()
			assert.True(t, ok)
			assert.InDelta(t, defaultRequestTimeout.Seconds(), time.Until(deadline).Seconds(), 0.5)
			return newTestResponse(req, http.StatusNoContent, ""), nil
		})
		req, err := client.NewRequest(http.MethodGet, "resources", nil)
		assert.NoError(t, err)

		_, err = client.Do(context.Background(), req, nil)
		assert.NoError(t, err)
	})

	t.Run("caller deadline is preserved", func(t *testing.T) {
		client := NewClient("auth-token")
		deadline := time.Now().Add(5 * time.Minute)
		ctx, cancel := context.WithDeadline(context.Background(), deadline)
		defer cancel()

		client.httpClient.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
			got, ok := req.Context().Deadline()
			assert.True(t, ok)
			assert.True(t, got.Equal(deadline), "deadline = %v, want %v", got, deadline)
			return newTestResponse(req, http.StatusNoContent, ""), nil
		})
		req, err := client.NewRequest(http.MethodGet, "resources", nil)
		assert.NoError(t, err)

		_, err = client.Do(ctx, req, nil)
		assert.NoError(t, err)
	})

	t.Run("custom zero-timeout client suppresses fallback", func(t *testing.T) {
		httpClient := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			_, ok := req.Context().Deadline()
			assert.False(t, ok)
			return newTestResponse(req, http.StatusNoContent, ""), nil
		})}
		client := NewClient("auth-token", WithHTTPClient(httpClient))
		req, err := client.NewRequest(http.MethodGet, "resources", nil)
		assert.NoError(t, err)

		_, err = client.Do(context.Background(), req, nil)
		assert.NoError(t, err)
	})

	t.Run("custom client timeout remains caller-owned", func(t *testing.T) {
		const customTimeout = 30 * time.Second
		httpClient := &http.Client{
			Timeout: customTimeout,
			Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				deadline, ok := req.Context().Deadline()
				assert.True(t, ok)
				assert.InDelta(t, customTimeout.Seconds(), time.Until(deadline).Seconds(), 0.5)
				return newTestResponse(req, http.StatusNoContent, ""), nil
			}),
		}
		client := NewClient("auth-token", WithHTTPClient(httpClient))
		req, err := client.NewRequest(http.MethodGet, "resources", nil)
		assert.NoError(t, err)

		_, err = client.Do(context.Background(), req, nil)
		assert.NoError(t, err)
		assert.Equal(t, customTimeout, client.httpClient.Timeout)
	})

	t.Run("already-cancelled context returns promptly", func(t *testing.T) {
		client := NewClient("auth-token")
		client.httpClient.Transport = roundTripFunc(func(req *http.Request) (*http.Response, error) {
			return nil, req.Context().Err()
		})
		req, err := client.NewRequest(http.MethodGet, "resources", nil)
		assert.NoError(t, err)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		started := time.Now()
		_, err = client.Do(ctx, req, nil)
		assert.ErrorIs(t, err, context.Canceled)
		assert.Less(t, time.Since(started), time.Second)
	})
}

func loadFixture(t *testing.T, fixtureName string) []byte {
	t.Helper()

	data, err := os.ReadFile(filepath.Join("testdata", fixtureName))
	if err != nil {
		t.Fatalf("failed to load fixture %q: %v", fixtureName, err)
	}

	return data
}

func mustTime(t *testing.T, timestamp string) *time.Time {
	t.Helper()

	ts, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		t.Fatalf("failed to parse timestamp %q: %v", timestamp, err)
	}

	return &ts
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func newTestResponse(req *http.Request, statusCode int, body string) *http.Response {
	return &http.Response{
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
		Request:    req,
		Status:     fmt.Sprintf("%d %s", statusCode, http.StatusText(statusCode)),
		StatusCode: statusCode,
	}
}
