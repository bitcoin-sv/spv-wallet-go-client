package notifications

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"time"

	"github.com/bitcoin-sv/spv-wallet/models"
)

// Webhook - the webhook event receiver
type Webhook struct {
	URL        string
	options    *WebhookOptions
	buffer     chan *models.RawEvent
	subscriber WebhookSubscriber
	handlers   *eventsMap
}

// NewWebhook - creates a new webhook
func NewWebhook(subscriber WebhookSubscriber, url string, opts ...WebhookOpts) *Webhook {
	options := NewWebhookOptions()
	for _, opt := range opts {
		opt(options)
	}

	wh := &Webhook{
		URL:        url,
		options:    options,
		buffer:     make(chan *models.RawEvent, options.BufferSize),
		subscriber: subscriber,
		handlers:   newEventsMap(),
	}
	for i := 0; i < options.Processors; i++ {
		go wh.process()
	}
	return wh
}

// Subscribe - sends a subscription request to the spv-wallet
func (w *Webhook) Subscribe(ctx context.Context) error {
	return w.subscriber.AdminSubscribeWebhook(ctx, w.URL, w.options.TokenHeader, w.options.TokenValue)
}

// Unsubscribe - sends an unsubscription request to the spv-wallet
func (w *Webhook) Unsubscribe(ctx context.Context) error {
	return w.subscriber.AdminUnsubscribeWebhook(ctx, w.URL)
}

// HTTPHandler - returns an http handler for the webhook; it should be registered with the http server
func (w *Webhook) HTTPHandler() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if w.options.TokenHeader != "" && r.Header.Get(w.options.TokenHeader) != w.options.TokenValue {
			http.Error(rw, "Unauthorized", http.StatusUnauthorized)
			return
		}
		var events []*models.RawEvent
		if err := json.NewDecoder(r.Body).Decode(&events); err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		for _, event := range events {
			select {
			case w.buffer <- event:
				// event sent
			case <-r.Context().Done():
				// request context canceled
				return
			case <-w.options.RootContext.Done():
				// root context canceled - the whole event processing has been stopped
				return
			case <-time.After(1 * time.Second):
				// timeout, most probably the channel is full
			}
		}
		rw.WriteHeader(http.StatusOK)
	})
}

func (w *Webhook) process() {
	for {
		select {
		case event := <-w.buffer:
			handler, ok := w.handlers.load(event.Type)
			if !ok {
				continue
			}
			model := reflect.New(handler.ModelType).Interface()
			if err := json.Unmarshal(event.Content, model); err != nil {
				continue
			}
			handler.Caller.Call([]reflect.Value{reflect.ValueOf(model)})
		case <-w.options.RootContext.Done():
			return
		}
	}
}
