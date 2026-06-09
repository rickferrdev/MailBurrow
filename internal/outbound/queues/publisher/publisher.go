package publisher

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rickferrdev/mail-burrow/internal/app/ports"
	"github.com/rickferrdev/mail-burrow/internal/config/rabbiting"
	"go.uber.org/fx"
)

var Provide = fx.Provide(
	fx.Annotate(
		New,
		fx.As(new(ports.QueuePublisher)),
	),
)

type Queue struct {
	rabbiting *rabbiting.Rabbiting
}

type Params struct {
	fx.In

	Rabbiting *rabbiting.Rabbiting
}

func New(params Params) (ports.QueuePublisher, error) {
	queue := Queue{
		rabbiting: params.Rabbiting,
	}

	return &queue, nil
}

func (queue *Queue) Publish(ctx context.Context, message ports.PublisherMessage) error {
	router, ok := ports.MapQueueRoutes[message.RoutingKey]
	if !ok {
		return fmt.Errorf("there is no registered routing key with the name \"%v\"", message.RoutingKey)
	}

	if !json.Valid(message.Payload) {
		return errors.New("the received payload is not a valid JSON")
	}

	channel, err := queue.rabbiting.Conn.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	if err := channel.PublishWithContext(
		ctx,
		string(router.ExchangeConfig.Name),
		string(router.BindingConfig.RoutingKey),
		false, false, amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         message.Payload,
		}); err != nil {
		return err
	}

	return nil
}
