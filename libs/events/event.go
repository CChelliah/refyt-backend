package events

import (
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
)

type Event struct {
	EventName EventName   `json:"eventName"`
	Data      interface{} `json:"data"`
}

func ToEventPayload(event interface{}, eventName string) (msg *message.Message, err error) {

	eventJSON, err := json.Marshal(event)

	if err != nil {
		return msg, err
	}

	msg = message.NewMessage(watermill.NewUUID(), eventJSON)
	msg.Metadata.Set("eventType", eventName)

	return msg, err
}
