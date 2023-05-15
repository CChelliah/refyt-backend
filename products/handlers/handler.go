package handlers

import (
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill/message"
	"refyt-backend/libs/events/evdata"
)

func Handler(msg *message.Message) ([]*message.Message, error) {

	var event evdata.ProductEvent

	err := json.Unmarshal(msg.Payload, &event)

	if err != nil {
		return nil, err
	}

	return nil, nil
}
