package handlers

import (
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill/message"
	"log"
	"refyt-backend/libs/events/evdata"
)

func Handler(msg *message.Message) ([]*message.Message, error) {

	var event evdata.CustomerEvent

	err := json.Unmarshal(msg.Payload, &event)

	if err != nil {
		return nil, err
	}

	log.Printf("recieved message in handler : %+v\n\n", event)

	return nil, nil
}
