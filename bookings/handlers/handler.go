package handlers

import (
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill/message"
	"refyt-backend/libs/events/evdata"
)

func CustomerHandler(msg *message.Message) ([]*message.Message, error) {

	var event evdata.CustomerEvent

	err := json.Unmarshal(msg.Payload, &event)

	if err != nil {
		return nil, err
	}

	return nil, nil
}
