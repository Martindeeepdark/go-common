package impl

import (
	"fmt"
	"os"

	"common/eventbus"
)

const (
	// MQTypeKey is the environment variable key for message queue type
	MQTypeKey = "COMMON_MQ_TYPE"

	// MQ types
	MQTypeNSQ    = "nsq"
	MQTypeKafka  = "kafka"
	MQTypeRMQ    = "rmq"
	MQTypePulsar = "pulsar"
	MQTypeNATS   = "nats"
)

type consumerServiceImpl struct{}

// NewConsumerService creates a new consumer service
func NewConsumerService() eventbus.ConsumerService {
	return &consumerServiceImpl{}
}

// RegisterConsumer registers a consumer based on the MQ type from environment
// Note: Implementations should be provided by scaffold packages
func (s *consumerServiceImpl) RegisterConsumer(nameServer, topic, group string, handler eventbus.ConsumerHandler, opts ...eventbus.ConsumerOpt) error {
	mqType := os.Getenv(MQTypeKey)
	if mqType == "" {
		mqType = MQTypeNSQ // Default to NSQ
	}

	// Return error indicating implementation should be in scaffold
	return fmt.Errorf("message queue implementation for '%s' should be provided by scaffold package. "+
		"See eventbus/impl/README.md for implementation examples", mqType)
}

// NewProducer creates a new producer based on the MQ type from environment
// Note: Implementations should be provided by scaffold packages
func NewProducer(nameServer, topic, group string, retries int) (eventbus.Producer, error) {
	mqType := os.Getenv(MQTypeKey)
	if mqType == "" {
		mqType = MQTypeNSQ // Default to NSQ
	}

	// Return error indicating implementation should be in scaffold
	return nil, fmt.Errorf("message queue implementation for '%s' should be provided by scaffold package. "+
		"See eventbus/impl/README.md for implementation examples", mqType)
}
