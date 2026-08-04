package sentrytrace

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/getsentry/sentry-go"
)

func newTestHub(t *testing.T) *sentry.Hub {
	t.Helper()

	client, err := sentry.NewClient(sentry.ClientOptions{
		Dsn:              "",
		EnableTracing:    true,
		TracesSampleRate: 1.0,
		Transport:        &sentry.HTTPSyncTransport{},
	})
	if err != nil {
		t.Fatalf("sentry.NewClient: %v", err)
	}
	return sentry.NewHub(client, sentry.NewScope())
}

func headerRecordingServer(t *testing.T, headers *http.Header) *httptest.Server {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*headers = r.Header.Clone()
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)
	return server
}

func TestTransport_InjectsHeadersFromActiveSpan(t *testing.T) {
	var receivedHeaders http.Header
	server := headerRecordingServer(t, &receivedHeaders)

	hub := newTestHub(t)
	ctx := sentry.SetHubOnContext(t.Context(), hub)
	transaction := sentry.StartTransaction(ctx, "test-transaction")
	defer transaction.Finish()

	httpClient := NewHTTPClient()
	req, err := http.NewRequestWithContext(transaction.Context(), http.MethodGet, server.URL, nil)
	if err != nil {
		t.Fatalf("http.NewRequestWithContext: %v", err)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		t.Fatalf("httpClient.Do: %v", err)
	}
	defer resp.Body.Close()

	sentryTrace := receivedHeaders.Get(sentry.SentryTraceHeader)
	if sentryTrace == "" {
		t.Fatal("expected sentry-trace header to be set")
	}
	traceID := transaction.TraceID.String()
	if got := sentryTrace[:len(traceID)]; got != traceID {
		t.Errorf("sentry-trace header trace id = %q, want %q", got, traceID)
	}
}

func TestTransport_InjectsHeadersWithoutSpan(t *testing.T) {
	var receivedHeaders http.Header
	server := headerRecordingServer(t, &receivedHeaders)

	hub := newTestHub(t)
	ctx := sentry.SetHubOnContext(t.Context(), hub)

	httpClient := NewHTTPClient()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, server.URL, nil)
	if err != nil {
		t.Fatalf("http.NewRequestWithContext: %v", err)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		t.Fatalf("httpClient.Do: %v", err)
	}
	defer resp.Body.Close()

	if receivedHeaders.Get(sentry.SentryTraceHeader) == "" {
		t.Error("expected sentry-trace header from propagation context to be set")
	}
}

func TestTransport_NoopWithoutHub(t *testing.T) {
	var receivedHeaders http.Header
	server := headerRecordingServer(t, &receivedHeaders)

	httpClient := NewHTTPClient()
	req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, server.URL, nil)
	if err != nil {
		t.Fatalf("http.NewRequestWithContext: %v", err)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		t.Fatalf("httpClient.Do: %v", err)
	}
	defer resp.Body.Close()

	if got := receivedHeaders.Get(sentry.SentryTraceHeader); got != "" {
		t.Errorf("expected no sentry-trace header without a configured hub, got %q", got)
	}
}
