package pubsub

import (
	"errors"
)

// Subscriber represents a subscriber to a specific topic in the pubsub system.
type Subscriber struct {
	topic string
	ch    <-chan interface{}
}

// NewSubscriber creates a new subscriber for the specified topic.
func NewSubscriber(ps PubSub, topic string) (*Subscriber, error) {
	ch, err := ps.Subscribe(topic)
	if err != nil {
		return nil, err
	}
	return &Subscriber{topic: topic, ch: ch}, nil
}

// ReceiveMessage blocks and waits for a new message from the subscribed topic.
func (s *Subscriber) ReceiveMessage() (interface{}, error) {
	msg, ok := <-s.ch
	if !ok {
		return nil, errors.New("channel closed")
	}
	return msg, nil
}

// Unsubscribe unsubscribes the subscriber from the topic.
func (s *Subscriber) Unsubscribe(ps PubSub) error {
	return ps.Unsubscribe(s.topic)
}
