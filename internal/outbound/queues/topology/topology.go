package topology

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rickferrdev/mail-burrow/internal/app/ports"
	"github.com/rickferrdev/mail-burrow/internal/config/rabbiting"
	"go.uber.org/fx"
)

var Provide = fx.Provide(New)
var Invoke = fx.Invoke(func(worker *Worker) error {
	return worker.Start()

})

type Worker struct {
	channel   *amqp.Channel
	rabbiting *rabbiting.Rabbiting
	life      fx.Lifecycle
}

type Params struct {
	fx.In

	Life      fx.Lifecycle
	Rabbiting *rabbiting.Rabbiting
}

func New(params Params) (*Worker, error) {
	worker := Worker{
		rabbiting: params.Rabbiting,
		life:      params.Life,
	}

	return &worker, nil
}

func (worker *Worker) Start() error {
	worker.life.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			channel, err := worker.rabbiting.Conn.Channel()
			if err != nil {
				return err
			}

			worker.channel = channel

			if err := worker.Topology(); err != nil {
				return err
			}

			return nil
		},
		OnStop: func(ctx context.Context) error {
			if worker.channel != nil {
				return worker.channel.Close()
			}

			return nil
		},
	})

	return nil
}

func (worker *Worker) Topology() error {
	for _, router := range ports.MapQueueRoutes {
		if err := worker.channel.ExchangeDeclare(
			string(router.ExchangeConfig.Name),
			string(router.ExchangeConfig.Type),
			router.ExchangeConfig.Durable,
			router.ExchangeConfig.AutoDelete,
			router.ExchangeConfig.Internal,
			router.ExchangeConfig.NoWait,
			amqp.Table(router.ExchangeConfig.Args),
		); err != nil {
			return err
		}

		if _, err := worker.channel.QueueDeclare(
			string(router.QueueConfig.Name),
			router.QueueConfig.Durable,
			router.QueueConfig.AutoDelete,
			router.QueueConfig.Exclusive,
			router.QueueConfig.NoWait,
			amqp.Table(router.QueueConfig.Args),
		); err != nil {
			return err
		}

		if err := worker.channel.QueueBind(
			string(router.BindingConfig.Queue),
			string(router.BindingConfig.RoutingKey),
			string(router.BindingConfig.Exchange),
			router.BindingConfig.NoWait,
			amqp.Table(router.BindingConfig.Args),
		); err != nil {
			return err
		}
	}

	return nil
}
