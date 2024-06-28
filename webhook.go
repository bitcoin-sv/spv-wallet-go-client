package walletclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"sync"
	"time"
)

const (
	eventBufferLength = 100
)

type Webhook struct {
	URL         string
	TokenHeader string
	TokenValue  string
	buffer      chan *RawEvent

	client   *WalletClient
	rootCtx  context.Context
	handlers *eventsMap
}

func NewWebhook(ctx context.Context, client *WalletClient, url, tokenHeader, tokenValue string, processors int) *Webhook {
	wh := &Webhook{
		URL:         url,
		TokenHeader: tokenHeader,
		TokenValue:  tokenValue,
		buffer:      make(chan *RawEvent, eventBufferLength),
		client:      client,
		rootCtx:     ctx,
		handlers:    newEventsMap(),
	}
	for i := 0; i < processors; i++ {
		go wh.process()
	}
	return wh
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
			case w.buffer <- event:
				// event sent
			case <-r.Context().Done():
				// request context cancelled
				return
			case <-w.rootCtx.Done():
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

func RegisterHandler[EventType Events](nd *Webhook, handlerFunction func(event *EventType)) error {
	handlerValue := reflect.ValueOf(handlerFunction)
	if handlerValue.Kind() != reflect.Func {
		return fmt.Errorf("Not a function")
	}

	modelType := handlerValue.Type().In(0)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}
	name := modelType.Name()

	nd.handlers.store(name, &eventHandler{
		Caller:    handlerValue,
		ModelType: modelType,
	})

	return nil
}

type eventHandler struct {
	Caller    reflect.Value
	ModelType reflect.Type
}

type eventsMap struct {
	registered *sync.Map
}

func newEventsMap() *eventsMap {
	return &eventsMap{
		registered: &sync.Map{},
	}
}

func (em *eventsMap) store(name string, handler *eventHandler) {
	em.registered.Store(name, handler)
}

func (em *eventsMap) load(name string) (*eventHandler, bool) {
	h, ok := em.registered.Load(name)
	if !ok {
		return nil, false
	}
	return h.(*eventHandler), true
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
		case <-nd.rootCtx.Done():
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
