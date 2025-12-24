package eventbus

import "context"

// Producer defines the message producer interface
type Producer interface {
	// Send sends a single message
	Send(ctx context.Context, body []byte, opts ...SendOpt) error
	// BatchSend sends multiple messages in batch
	BatchSend(ctx context.Context, bodyArr [][]byte, opts ...SendOpt) error
}

// ConsumerService defines the consumer service interface
type ConsumerService interface {
	// RegisterConsumer registers a message consumer
	RegisterConsumer(nameServer, topic, group string, handler ConsumerHandler, opts ...ConsumerOpt) error
}

// ConsumerHandler defines the message handler interface
type ConsumerHandler interface {
	// HandleMessage handles the received message
	HandleMessage(ctx context.Context, msg *Message) error
}

// Message represents a standard message format
type Message struct {
	Topic string // Message topic
	Group string // Consumer group name
	Body  []byte // Message body
}

// SendOpt is the function type for send options
type SendOpt func(option *SendOption)

// SendOption contains optional configuration for sending messages
type SendOption struct {
	ShardingKey *string // Sharding key for message routing
}

// WithShardingKey creates a send option with sharding key
func WithShardingKey(key string) SendOpt {
	return func(o *SendOption) {
		o.ShardingKey = &key
	}
}

// ConsumerOpt is the function type for consumer options
type ConsumerOpt func(option *ConsumerOption)

// ConsumerOption contains optional configuration for consuming messages
type ConsumerOption struct {
	Orderly *bool // Whether to consume in order
}

// WithConsumerOrderly creates a consumer option for ordered consumption
func WithConsumerOrderly(orderly bool) ConsumerOpt {
	return func(o *ConsumerOption) {
		o.Orderly = &orderly
	}
}

// Global default service instance
var defaultSVC ConsumerService

// SetDefaultSVC sets the default consumer service
func SetDefaultSVC(svc ConsumerService) {
	defaultSVC = svc
}

// GetDefaultSVC returns the default consumer service
func GetDefaultSVC() ConsumerService {
	return defaultSVC
}
