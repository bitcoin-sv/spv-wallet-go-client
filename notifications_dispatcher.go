package walletclient

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

type Handlers struct {
	GeneralPurposeEvent []func(*GeneralPurposeEvent)
}

type NotificationsDispatcher struct {
	ctx      context.Context
	input    chan *RawEvent
	handlers Handlers
}

func NewNotificationsDispatcher(ctx context.Context, inputChannel chan *RawEvent) *NotificationsDispatcher {
	obj := &NotificationsDispatcher{
		ctx:   ctx,
		input: inputChannel,
	}

	return obj
}

func (nd *NotificationsDispatcher) process() {
	for {
		select {
		case event := <-nd.input:
			switch event.Type {
			case "general-purpose-event":
				content, err := GetEventContent[GeneralPurposeEvent](event)
				if err != nil {
					fmt.Println("Error getting event content")
					continue
				}
				for _, handler := range nd.handlers.GeneralPurposeEvent {
					handler(content)
				}
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

func (GeneralPurposeEvent) GetType() string {
	return "general-purpose-event"
}

func GetEventContent[modelType EventContent](raw *RawEvent) (*modelType, error) {
	model := *new(modelType)
	if raw.Type != model.GetType() {
		return nil, fmt.Errorf("Wrong type")
	}

	if err := json.Unmarshal(raw.Content, &model); err != nil {
		return nil, errors.Wrap(err, "Cannot unmarshall the content json")
	}
	return &model, nil
}

func NewRawEvent(namedEvent EventContent) *RawEvent {
	asJson, _ := json.Marshal(namedEvent)
	return &RawEvent{
		Type:    namedEvent.GetType(),
		Content: asJson,
	}
}
