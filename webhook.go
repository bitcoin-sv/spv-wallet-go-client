package walletclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	eventBufferLength = 100
)

type Webhook struct {
	URL         string
	TokenHeader string
	TokenValue  string
	Channel     chan *RawEvent

	client *WalletClient
}

func NewWebhook(client *WalletClient, url, tokenHeader, tokenValue string) *Webhook {
	return &Webhook{
		URL:         url,
		TokenHeader: tokenHeader,
		TokenValue:  tokenValue,
		Channel:     make(chan *RawEvent, eventBufferLength),
		client:      client,
	}
}

func (w *Webhook) Subscribe(ctx context.Context) ResponseError {
	return w.client.AdminSubscribeWebhook(ctx, w.URL, w.TokenHeader, w.TokenValue)
}

func (w *Webhook) Unsubscribe(ctx context.Context) ResponseError {
	return w.client.AdminUnsubscribeWebhook(ctx, w.URL)
}

func (w *Webhook) HTTPHandler() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if w.TokenHeader != "" && r.Header.Get(w.TokenHeader) != w.TokenValue {
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
			case w.Channel <- event:
				// event sent
			case <-r.Context().Done():
				// context cancelled
				return
			case <-time.After(1 * time.Second):
				// timeout, most probably the channel is full
				// TODO: log this
			}
		}
		rw.WriteHeader(http.StatusOK)
	})
}
