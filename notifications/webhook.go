package notifications

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"
)

func NewWebhook(ctx context.Context, subscriber WebhookSubscriber, url string, opts ...WebhookOpts) *Webhook {
	options := NewWebhookOptions()
	for _, opt := range opts {
		opt(options)
	}

	wh := &Webhook{
		URL:        url,
		options:    options,
		buffer:     make(chan *RawEvent, options.BufferSize),
		subscriber: subscriber,
		handlers:   newEventsMap(),
	}
	for i := 0; i < options.Processors; i++ {
		go wh.process()
	}
	return wh
}

func (w *Webhook) Subscribe(ctx context.Context) error {
	return w.subscriber.AdminSubscribeWebhook(ctx, w.URL, w.options.TokenHeader, w.options.TokenValue)
}

func (w *Webhook) Unsubscribe(ctx context.Context) error {
	return w.subscriber.AdminUnsubscribeWebhook(ctx, w.URL)
}

func (w *Webhook) HTTPHandler() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if w.options.TokenHeader != "" && r.Header.Get(w.options.TokenHeader) != w.options.TokenValue {
			http.Error(rw, "Unauthorized", http.StatusUnauthorized)
			return
		}
		var events []*RawEvent
		if err := json.NewDecoder(r.Body).Decode(&events); err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Printf("Received: %v\n", events)
		for _, event := range events {
			select {
			case w.buffer <- event:
				// event sent
			case <-r.Context().Done():
				// request context cancelled
				return
			case <-w.options.RootContext.Done():
				// root context cancelled - the whole event processing has been stopped
				return
			case <-time.After(1 * time.Second):
				// timeout, most probably the channel is full
				// TODO: log this
			}
		}
		rw.WriteHeader(http.StatusOK)
	})
}

func (nd *Webhook) process() {
	for {
		select {
		case event := <-nd.buffer:
			handler, ok := nd.handlers.load(event.Type)
			if !ok {
				fmt.Printf("No handlers for %s event type", event.Type)
				continue
			}
			model := reflect.New(handler.ModelType).Interface()
			if err := json.Unmarshal(event.Content, model); err != nil {
				fmt.Println("Cannot unmarshall the content json")
				continue
			}
			handler.Caller.Call([]reflect.Value{reflect.ValueOf(model)})
		case <-nd.options.RootContext.Done():
			return
		}
	}
}

////////////////////// BELOW it should be imported from spv-wallet models

type RawEvent struct {
	Type    string          `json:"type"`
	Content json.RawMessage `json:"content"`
}

type StringEvent struct {
	Value string
}

type NumericEvent struct {
	Numeric int
}

type Events interface {
	StringEvent | NumericEvent
}
