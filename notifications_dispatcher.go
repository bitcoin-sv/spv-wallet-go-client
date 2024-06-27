package walletclient

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
)

type Handler struct {
	HandlerFunc any // any as function to handle the event
	Model       EventContent
}

type NotificationsDispatcher struct {
	ctx      context.Context
	input    chan *RawEvent
	handlers map[string][]Handler
}

func NewNotificationsDispatcher(ctx context.Context, inputChannel chan *RawEvent, providedHandlers []Handler) *NotificationsDispatcher {
	dispatcher := &NotificationsDispatcher{
		ctx:      ctx,
		input:    inputChannel,
		handlers: make(map[string][]Handler, len(providedHandlers)),
	}

	for _, handler := range providedHandlers {
		dispatcher.handlers[handler.Model.GetType()] = append(dispatcher.handlers[handler.Model.GetType()], handler)
	}

	go dispatcher.process()

	return dispatcher
}

func (nd *NotificationsDispatcher) process() {
	for {
		select {
		case event := <-nd.input:
			handlers, ok := nd.handlers[event.Type]
			if !ok {
				fmt.Printf("No handlers for %s event type", event.Type)
				continue
			}
			for _, handler := range handlers {
				modelSource := handler.Model
				// copy the event to the model, use reflection
				model := reflect.New(reflect.TypeOf(modelSource).Elem()).Interface()
				if err := json.Unmarshal(event.Content, model); err != nil {
					fmt.Println("Cannot unmarshall the content json")
					continue
				}
				// use reflect
				handlerValue := reflect.ValueOf(handler.HandlerFunc)
				if handlerValue.Kind() != reflect.Func {
					fmt.Println("Not a function")
					continue
				}
				if handlerValue.Type().NumIn() != 1 {
					fmt.Println("Wrong number of arguments")
					continue
				}
				handlerValue.Call([]reflect.Value{reflect.ValueOf(model)})
			}
		case <-nd.ctx.Done():
			return
		}
	}
}

////////////////////// BELOW it should be imported from spv-wallet models

// RawEvent - event type
type RawEvent struct {
	Type    string          `json:"type"`
	Content json.RawMessage `json:"content"`
}

type EventContent interface {
	GetType() string
}

type GeneralPurposeEvent struct {
	Value string
}

func (*GeneralPurposeEvent) GetType() string {
	return "general-purpose-event"
}
