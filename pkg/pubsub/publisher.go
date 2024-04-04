package pubsub

// Publisher represents a publisher that publishes messages to a specific topic in the pubsub system.
type Publisher struct {
	ps PubSub
}

// NewPublisher creates a new publisher.
func NewPublisher(ps PubSub) *Publisher {
	return &Publisher{ps: ps}
}

// PublishMessage publishes a message to the specified topic.
func (p *Publisher) PublishMessage(topic string, message interface{}) error {
	return p.ps.Publish(topic, message)
}
