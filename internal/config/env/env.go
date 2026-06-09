package env

import (
	"github.com/rickferrdev/dotenv"
	_ "github.com/rickferrdev/dotenv/auto"
	"go.uber.org/fx"
)

var Provide = fx.Provide(New)

type Environment struct {
	AmqpUrl          string `env:"AMQP_URL" required:"true"`
	RabbitMQPrefetch string `env:"RABBIT_MQ_PREFETCH" default:"5"`
	ServerPort       string `env:"SERVER_PORT" default:"8080"`
	ServerHost       string `env:"SERVER_HOST" default:"0.0.0.0"`
	MailerHost       string `env:"MAILER_HOST" required:"true"`
	MailerPort       string `env:"MAILER_PORT" required:"true"`
	MailerUsername   string `env:"MAILER_USERNAME" required:"true"`
	MailerPassword   string `env:"MAILER_PASSWORD" required:"true"`
	DatabaseUrl      string `env:"DATABASE_URL" required:"true"`
}

func New() (*Environment, error) {
	var env Environment
	if err := dotenv.Unmarshal(&env); err != nil {
		return nil, err
	}

	return &env, nil
}
