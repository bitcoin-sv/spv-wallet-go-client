package notifications

import (
	"reflect"

	"github.com/bitcoin-sv/spv-wallet/models"
)

type eventHandler struct {
	Caller    reflect.Value
	ModelType reflect.Type
}

// RegisterHandler - registers a handler for a specific event type
func RegisterHandler[EventType models.Events](nd *Webhook, handlerFunction func(event *EventType)) error {
	handlerValue := reflect.ValueOf(handlerFunction)

	modelType := handlerValue.Type().In(0).Elem()
	name := modelType.Name()

	nd.handlers.store(name, &eventHandler{
		Caller:    handlerValue,
		ModelType: modelType,
	})

	return nil
}
