package pubsub_test

import (
	"testing"

	"github.com/hitesh22rana/eventhub/pkg/pubsub"
)

const (
	topic = "test"
)

func TestPubSub(t *testing.T) {
	// Test initializing the pubsub system
	config := pubsub.WithTopics(topic)
	ps, err := pubsub.NewPubSub(config)
	if err != nil {
		t.Errorf("Failed to initialize pubsub: %v", err)
	}

	// Test subscribing to a topic
	ch, err := ps.Subscribe("test")
	if err != nil {
		t.Errorf("Failed to subscribe to topic: %v", err)
	}
	defer ps.Unsubscribe("test")

	// Test publishing a message to the topic
	expectedMessage := "test message"
	go func() {
		if err := ps.Publish("test", expectedMessage); err != nil {
			t.Errorf("Failed to publish message: %v", err)
		}
	}()
	receivedMessage := <-ch
	if receivedMessage != expectedMessage {
		t.Errorf("Received message differs from expected. Got: %s, Expected: %s", receivedMessage, expectedMessage)
	}

	// Test unsubscribing from the topic
	err = ps.Unsubscribe("test")
	if err != nil {
		t.Errorf("Failed to unsubscribe from topic: %v", err)
	}

	// Test closing the pubsub system
	err = ps.Close()
	if err != nil {
		t.Errorf("Failed to close pubsub: %v", err)
	}
}
