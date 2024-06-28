package walletclient

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"sync"
)

type Handler struct {
	Caller    reflect.Value
	ModelType reflect.Type
}

type NotificationsDispatcher struct {
	ctx      context.Context
	input    chan *RawEvent
	handlers *sync.Map
}

func NewNotificationsDispatcher(ctx context.Context, inputChannel chan *RawEvent) *NotificationsDispatcher {
	dispatcher := &NotificationsDispatcher{
		ctx:      ctx,
		input:    inputChannel,
		handlers: &sync.Map{},
	}

	go dispatcher.process()

	return dispatcher
}

func RegisterHandler[EventType Events](nd *NotificationsDispatcher, handlerFunction func(event *EventType)) error {
	handlerValue := reflect.ValueOf(handlerFunction)
	if handlerValue.Kind() != reflect.Func {
		return fmt.Errorf("Not a function")
	}

	modelType := handlerValue.Type().In(0)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}
	name := modelType.Name()

	nd.handlers.Store(name, Handler{
		Caller:    handlerValue,
		ModelType: modelType,
	})
	return nil
}

func (nd *NotificationsDispatcher) process() {
	for {
		select {
		case event := <-nd.input:
			h, ok := nd.handlers.Load(event.Type)
			if !ok {
				fmt.Printf("No handlers for %s event type", event.Type)
				continue
			}
			handler := h.(Handler)

			model := reflect.New(handler.ModelType).Interface()
			if err := json.Unmarshal(event.Content, model); err != nil {
				fmt.Println("Cannot unmarshall the content json")
				continue
			}
			handler.Caller.Call([]reflect.Value{reflect.ValueOf(model)})
		case <-nd.ctx.Done():
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
