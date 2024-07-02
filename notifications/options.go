package notifications

import (
	"context"

	"github.com/bitcoin-sv/spv-wallet/models"
)

type WebhookOptions struct {
	TokenHeader string
	TokenValue  string
	BufferSize  int
	RootContext context.Context
	Processors  int
}

func NewWebhookOptions() *WebhookOptions {
	return &WebhookOptions{
		TokenHeader: "",
		TokenValue:  "",
		BufferSize:  100,
		Processors:  1,
		RootContext: context.Background(),
	}
}

type WebhookOpts = func(*WebhookOptions)

func WithToken(tokenHeader, tokenValue string) WebhookOpts {
	return func(w *WebhookOptions) {
		w.TokenHeader = tokenHeader
		w.TokenValue = tokenValue
	}
}

func WithBufferSize(size int) WebhookOpts {
	return func(w *WebhookOptions) {
		w.BufferSize = size
	}
}

func WithRootContext(ctx context.Context) WebhookOpts {
	return func(w *WebhookOptions) {
		w.RootContext = ctx
	}
}

func WithProcessors(count int) WebhookOpts {
	return func(w *WebhookOptions) {
		w.Processors = count
	}
}

type Webhook struct {
	URL        string
	options    *WebhookOptions
	buffer     chan *models.RawEvent
	subscriber WebhookSubscriber
	handlers   *eventsMap
}
