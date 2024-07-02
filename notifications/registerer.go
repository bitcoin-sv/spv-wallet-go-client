package notifications

import (
	"fmt"
	"reflect"

	"github.com/bitcoin-sv/spv-wallet/models"
)

type eventHandler struct {
	Caller    reflect.Value
	ModelType reflect.Type
}

func RegisterHandler[EventType models.Events](nd *Webhook, handlerFunction func(event *EventType)) error {
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
