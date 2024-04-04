package pubsub

import (
	"errors"
	"sync"
)

var (
	ErrClosed             = errors.New("pubsub: closed")
	ErrNoTopic            = errors.New("pubsub: no such topic")
	ErrUnsubscribed       = errors.New("pubsub: not subscribed to the topic")
	ErrTopicAlreadyExists = errors.New("pubsub: topic already exists")
)

type PubSub interface {
	// Publish publishes a message to the topic.
	Publish(topic string, message interface{}) error

	// Subscribe subscribes to the topic and returns a channel that will receive messages.
	Subscribe(topic string) (<-chan interface{}, error)

	// Unsubscribe unsubscribes from the topic.
	Unsubscribe(topic string) error

	// Close closes the PubSub and unsubscribes all topics.
	Close() error

	// CreateTopic creates a new topic.
	CreateTopic(topic string) error
}

type Config struct {
	Topics []string
}

type pubsub struct {
	// mu protects the topics map.
	mu sync.RWMutex

	// subscribers contains the channels that will receive messages for each topic.
	subscribers map[string][]chan interface{}

	// closed is set to true when Close is called.
	closed bool
}

func NewPubSub(config Config) (PubSub, error) {
	ps := &pubsub{
		mu:          sync.RWMutex{},
		subscribers: make(map[string][]chan interface{}),
		closed:      false,
	}

	for _, topic := range config.Topics {
		if err := ps.CreateTopic(topic); err != nil {
			return nil, err
		}
	}

	return ps, nil
}

func WithTopics(topics ...string) Config {
	return Config{Topics: topics}
}

func (ps *pubsub) Publish(topic string, message interface{}) error {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	if ps.closed {
		return ErrClosed
	}

	channels, ok := ps.subscribers[topic]
	if !ok {
		return ErrNoTopic
	}

	for _, ch := range channels {
		go func(ch chan interface{}) {
			ch <- message
		}(ch)
	}

	return nil
}

func (ps *pubsub) Subscribe(topic string) (<-chan interface{}, error) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if ps.closed {
		return nil, ErrClosed
	}

	ch := make(chan interface{}, 1)
	ps.subscribers[topic] = append(ps.subscribers[topic], ch)
	return ch, nil
}

func (ps *pubsub) Unsubscribe(topic string) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if ps.closed {
		return ErrClosed
	}

	channels, ok := ps.subscribers[topic]
	if !ok {
		return ErrNoTopic
	}

	// Close all channels subscribed to the given topic
	for _, ch := range channels {
		close(ch)
	}

	// Clear subscribers for the given topic
	delete(ps.subscribers, topic)

	return nil
}

func (ps *pubsub) Close() error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if ps.closed {
		return ErrClosed
	}

	ps.closed = true
	for _, chs := range ps.subscribers {
		for _, ch := range chs {
			close(ch)
		}
	}

	return nil
}

func (ps *pubsub) CreateTopic(topic string) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if ps.closed {
		return ErrClosed
	}

	// Check if topic already exists
	if _, ok := ps.subscribers[topic]; ok {
		return ErrTopicAlreadyExists
	}

	// Create new topic
	ps.subscribers[topic] = make([]chan interface{}, 0)
	return nil
}
