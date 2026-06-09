package ports

import (
	"context"
)

const (
	QueueEmailProcessing      Queue         = "queue.email.processing"
	ExchangeEmailProcessing   ExchangeValue = "email.processing"
	RoutingKeyEmailProcessing RoutingKey    = "email.processing"

	QueueEmailProcessingRetry      Queue         = "queue.emails.processing.retry"
	ExchangeEmailProcessingRetry   ExchangeValue = "emails.processing.retry"
	RoutingKeyEmailProcessingRetry RoutingKey    = "emails.processing.retry"

	QueueEmailProcessingDlq      Queue         = "queue.emails.processing.dlq"
	ExchangeEmailProcessingDlx   ExchangeValue = "emails.processing.dlx"
	RoutingKeyEmailProcessingDlq RoutingKey    = "emails.processing.dlq"

	ExchangeDirect ExchangeType = "direct"
)

type (
	RoutingKey    string
	ExchangeType  string
	ExchangeValue string
	Queue         string
	Table         map[string]any
)

type ExchangeConfig struct {
	Name       ExchangeValue
	Type       ExchangeType
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       Table
}

type QueueConfig struct {
	Name       Queue
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       Table
}

type BindingConfig struct {
	Queue      Queue
	Exchange   ExchangeValue
	RoutingKey RoutingKey
	NoWait     bool
	Args       Table
}

type RouterConfig = struct {
	BindingConfig  BindingConfig
	QueueConfig    QueueConfig
	ExchangeConfig ExchangeConfig
}

var MapQueueRoutes = map[RoutingKey]RouterConfig{
	RoutingKeyEmailProcessing: {
		BindingConfig: BindingConfig{
			Queue:      QueueEmailProcessing,
			Exchange:   ExchangeEmailProcessing,
			RoutingKey: RoutingKeyEmailProcessing,
		},
		QueueConfig: QueueConfig{
			Name:    QueueEmailProcessing,
			Durable: true,
			Args: Table{
				"x-dead-letter-exchange":    string(ExchangeEmailProcessingRetry),
				"x-dead-letter-routing-key": string(RoutingKeyEmailProcessingRetry),
			},
		},
		ExchangeConfig: ExchangeConfig{
			Name:    ExchangeEmailProcessing,
			Type:    ExchangeDirect,
			Durable: true,
		},
	},
	RoutingKeyEmailProcessingDlq: {
		BindingConfig: BindingConfig{
			Queue:      QueueEmailProcessingDlq,
			Exchange:   ExchangeEmailProcessingDlx,
			RoutingKey: RoutingKeyEmailProcessingDlq,
			NoWait:     false,
		},
		QueueConfig: QueueConfig{
			Name:    QueueEmailProcessingDlq,
			Durable: true,
		},
		ExchangeConfig: ExchangeConfig{
			Name:    ExchangeEmailProcessingDlx,
			Type:    ExchangeDirect,
			Durable: true,
		},
	},
	RoutingKeyEmailProcessingRetry: {
		BindingConfig: BindingConfig{
			Queue:      QueueEmailProcessingRetry,
			Exchange:   ExchangeEmailProcessingRetry,
			RoutingKey: RoutingKeyEmailProcessingRetry,
		},
		QueueConfig: QueueConfig{
			Name:    QueueEmailProcessingRetry,
			Durable: true,
			Args: Table{
				"x-message-ttl":             int32(15_000),
				"x-dead-letter-exchange":    string(ExchangeEmailProcessing),
				"x-dead-letter-routing-key": string(RoutingKeyEmailProcessing),
			},
		},
		ExchangeConfig: ExchangeConfig{
			Name:    ExchangeEmailProcessingRetry,
			Type:    ExchangeDirect,
			Durable: true,
		},
	},
}

type PublisherMessage struct {
	RoutingKey RoutingKey
	Payload    []byte
}

//go:generate mockgen -source=$GOFILE -destination=../../tests/mocks/queues.go -package=mocks
type QueuePublisher interface {
	Publish(ctx context.Context, message PublisherMessage) error
}
