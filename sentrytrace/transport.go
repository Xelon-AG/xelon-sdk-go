// Package sentrytrace provides an http.RoundTripper that instruments outgoing
// requests with Sentry performance spans and distributed tracing headers
// (sentry-trace and baggage), so that services called through the Xelon SDK can
// continue the trace.
//
// It is fully optional: importing this package is the only thing that pulls
// sentry-go into a consumer's build. Wire it into the SDK client via
// xelon.WithHTTPClient:
//
//	client := xelon.NewClient(token, xelon.WithHTTPClient(sentrytrace.NewHTTPClient()))
//
// Headers are only injected when the request context carries an active Sentry
// span (or a hub with a valid client); otherwise the transport is a no-op
// passthrough.
package sentrytrace

import (
	"context"
	"net/http"

	"github.com/getsentry/sentry-go"
)

// Transport is an http.RoundTripper that starts a Sentry child span for each
// outgoing request and injects sentry-trace/baggage headers.
type Transport struct {
	// Base is the underlying RoundTripper used to execute the request.
	// http.DefaultTransport is used when nil.
	Base http.RoundTripper
}

var _ http.RoundTripper = (*Transport)(nil)

// NewHTTPClient returns an *http.Client using Transport with the default base,
// suitable for passing to xelon.WithHTTPClient.
func NewHTTPClient() *http.Client {
	return &http.Client{Transport: &Transport{}}
}

// RoundTrip implements http.RoundTripper.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	base := t.Base
	if base == nil {
		base = http.DefaultTransport
	}

	ctx := req.Context()

	span := sentry.SpanFromContext(ctx)
	if span == nil {
		return t.roundTripWithoutSpan(ctx, req, base)
	}

	child := span.StartChild("http.client", sentry.WithDescription(req.Method+" "+req.URL.Path))
	defer child.Finish()

	req = req.Clone(ctx)
	req.Header.Set(sentry.SentryTraceHeader, child.ToSentryTrace())
	if baggage := child.ToBaggage(); baggage != "" {
		req.Header.Set(sentry.SentryBaggageHeader, baggage)
	}

	resp, err := base.RoundTrip(req)
	if err != nil {
		child.Status = sentry.SpanStatusInternalError
		return resp, err
	}
	child.Status = sentry.HTTPtoSpanStatus(resp.StatusCode)
	return resp, err
}

// roundTripWithoutSpan propagates the trace from the hub's propagation context
// ("tracing without performance") when no span is active.
func (t *Transport) roundTripWithoutSpan(ctx context.Context, req *http.Request, base http.RoundTripper) (*http.Response, error) {
	hub := sentry.GetHubFromContext(ctx)
	if hub == nil {
		hub = sentry.CurrentHub()
	}
	if hub.Client() == nil {
		return base.RoundTrip(req)
	}

	req = req.Clone(ctx)
	req.Header.Set(sentry.SentryTraceHeader, hub.GetTraceparent())
	if baggage := hub.GetBaggage(); baggage != "" {
		req.Header.Set(sentry.SentryBaggageHeader, baggage)
	}
	return base.RoundTrip(req)
}
