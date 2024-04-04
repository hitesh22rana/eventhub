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
	config := pubsub.WithTopics(topic)
	ps, err := pubsub.NewPubSub(config)
	if err != nil {
		fmt.Println("Error initializing pubsub:", err)
		return
	}

	publisher := pubsub.NewPublisher(ps)
	go func() {
		defer ps.Close()
		for i := 0; i < 10; i++ {
			if err := publisher.PublishMessage(topic, "Hello, world!"); err != nil {
				fmt.Println("Error publishing message:", err)
				return
			}
			time.Sleep(time.Millisecond * 100)
		}
	}()

	subscriber, err := pubsub.NewSubscriber(ps, topic)
	if err != nil {
		fmt.Println("Error creating subscriber:", err)
		return
	}

	defer subscriber.Unsubscribe(ps)

	for {
		message, err := subscriber.ReceiveMessage()
		if err != nil {
			fmt.Println("Error receiving message:", err)
			return
		}
		fmt.Println("Received message:", message)
	}
}
