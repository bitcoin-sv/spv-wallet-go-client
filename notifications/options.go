package notifications

import (
	"context"
	"runtime"
)

// WebhookOptions - options for the webhook
type WebhookOptions struct {
	TokenHeader string
	TokenValue  string
	BufferSize  int
	RootContext context.Context
	Processors  int
}

// NewWebhookOptions - creates a new webhook options
func NewWebhookOptions() *WebhookOptions {
	return &WebhookOptions{
		TokenHeader: "",
		TokenValue:  "",
		BufferSize:  100,
		Processors:  runtime.NumCPU(),
		RootContext: context.Background(),
	}
}

// WebhookOpts - functional options for the webhook
type WebhookOpts = func(*WebhookOptions)

// WithToken - sets the token header and value
func WithToken(tokenHeader, tokenValue string) WebhookOpts {
	return func(w *WebhookOptions) {
		w.TokenHeader = tokenHeader
		w.TokenValue = tokenValue
	}
}

// WithBufferSize - sets the buffer size
func WithBufferSize(size int) WebhookOpts {
	return func(w *WebhookOptions) {
		w.BufferSize = size
	}
}

// WithRootContext - sets the root context
func WithRootContext(ctx context.Context) WebhookOpts {
	return func(w *WebhookOptions) {
		w.RootContext = ctx
	}
}

// WithProcessors - sets the number of concurrent loops which will process the events
func WithProcessors(count int) WebhookOpts {
	return func(w *WebhookOptions) {
		w.Processors = count
	}
}
