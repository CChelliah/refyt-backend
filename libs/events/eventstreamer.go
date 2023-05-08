package events

import (
	"context"
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
)

type IEventStreamer interface {
	Publish(topic string, messages ...*message.Message) error
	Subscribe(ctx context.Context, topic string) (<-chan *message.Message, error)
	Close() error
	PublishEvent(topic Topic, event Event) (err error)
}

type EventStreamer struct {
	goChannel *gochannel.GoChannel
}

func NewEventStreamer(logger watermill.LoggerAdapter) *EventStreamer {
	return &EventStreamer{
		goChannel: gochannel.NewGoChannel(gochannel.Config{}, logger),
	}
}

func (e *EventStreamer) Publish(topic string, messages ...*message.Message) error {
	return e.goChannel.Publish(topic, messages...)
}

func (e *EventStreamer) Subscribe(ctx context.Context, topic string) (<-chan *message.Message, error) {
	return e.goChannel.Subscribe(ctx, topic)
}

func (e *EventStreamer) Close() error {
	return e.goChannel.Close()
}

func (e *EventStreamer) PublishEvent(topic Topic, event Event) (err error) {

	eventJSON, err := json.Marshal(event)

	if err != nil {
		return err
	}

	msg := message.NewMessage(watermill.NewUUID(), eventJSON)

	err = e.Publish(string(topic), msg)

	if err != nil {
		return err
	}

	return err
}
