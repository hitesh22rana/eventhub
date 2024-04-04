package main

import (
	"fmt"
	"time"

	"github.com/hitesh22rana/eventhub/pkg/pubsub"
)

const (
	topic = "example"
)

func main() {
	// Create a new pubsub instance
	ps := pubsub.NewPubSub()

	// Subscribe to a topic
	ch, err := ps.Subscribe(topic)
	if err != nil {
		fmt.Println("Error subscribing to topic:", err)
		return
	}

	// Listen for messages on the subscribed channel
	go func() {
		for message := range ch {
			fmt.Printf("Received message from topic '%s': %v\n", topic, message)
		}
	}()

	// Publish some messages to the topic
	for i := 1; i <= 3; i++ {
		message := fmt.Sprintf("Message %d", i)
		if err := ps.Publish(topic, message); err != nil {
			fmt.Println("Error publishing message:", err)
			return
		}
		fmt.Printf("Published: %s\n", message)
		time.Sleep(time.Second)
	}

	// Unsubscribe from the topic
	if err := ps.Unsubscribe(topic); err != nil {
		fmt.Println("Error unsubscribing from topic:", err)
		return
	}
	fmt.Println("Unsubscribed from topic")

	// Closing pubsub
	if err := ps.Close(); err != nil {
		fmt.Println("Error closing pubsub:", err)
		return
	}
	fmt.Println("PubSub closed")
}
