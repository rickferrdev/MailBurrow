package rabbiting

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rickferrdev/mail-burrow/internal/config/env"
	"go.uber.org/fx"
)

var Provide = fx.Provide(New)

type Rabbiting struct {
	Conn *amqp.Connection
}

type Params struct {
	fx.In

	Env  *env.Environment
	Life fx.Lifecycle
}

func New(params Params) (*Rabbiting, error) {
	dial, err := amqp.Dial(params.Env.AmqpUrl)
	if err != nil {
		return nil, err
	}

	params.Life.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			if dial != nil {
				return dial.Close()
			}

			return nil
		},
	})

	rabbiting := Rabbiting{
		Conn: dial,
	}

	return &rabbiting, nil
}
