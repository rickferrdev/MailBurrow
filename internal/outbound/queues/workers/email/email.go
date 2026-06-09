package email

import (
	"context"
	"encoding/json"
	"strconv"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rickferrdev/mail-burrow/internal/app/constants"
	"github.com/rickferrdev/mail-burrow/internal/app/domain"
	"github.com/rickferrdev/mail-burrow/internal/app/ports"
	"github.com/rickferrdev/mail-burrow/internal/config/env"
	"github.com/rickferrdev/mail-burrow/internal/config/rabbiting"
	"github.com/rickferrdev/mail-burrow/internal/outbound/queues/topology"
	"go.uber.org/fx"
)

var Invoke = fx.Invoke(func(worker *Worker) error {
	return worker.Start()
})

var Provide = fx.Provide(New)

type Worker struct {
	env          *env.Environment
	life         fx.Lifecycle
	mailer       ports.Mailer
	cancel       context.CancelFunc
	channel      *amqp.Channel
	amqpConn     *rabbiting.Rabbiting
	publisher    ports.QueuePublisher
	emailStorage ports.EmailStorage
}

type Params struct {
	fx.In

	Env            *env.Environment
	Life           fx.Lifecycle
	Mailer         ports.Mailer
	AMQPConn       *rabbiting.Rabbiting
	EmailStorage   ports.EmailStorage
	QueuePublisher ports.QueuePublisher
	Topology       *topology.Worker
}

func New(params Params) (*Worker, error) {
	worker := Worker{
		env:          params.Env,
		publisher:    params.QueuePublisher,
		mailer:       params.Mailer,
		life:         params.Life,
		emailStorage: params.EmailStorage,
		amqpConn:     params.AMQPConn,
	}

	return &worker, nil
}

func (worker *Worker) Start() error {
	worker.life.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			prefetch, err := strconv.Atoi(worker.env.RabbitMQPrefetch)
			if err != nil {
				return err
			}

			channel, err := worker.amqpConn.Conn.Channel()
			if err != nil {
				return err
			}
			worker.channel = channel

			if err := worker.channel.Qos(
				prefetch,
				0,
				false,
			); err != nil {
				return err
			}

			router := ports.MapQueueRoutes[ports.RoutingKeyEmailProcessing]

			_, err = worker.channel.QueueDeclare(
				string(router.QueueConfig.Name),
				router.QueueConfig.Durable,
				router.QueueConfig.AutoDelete,
				router.QueueConfig.Exclusive,
				router.QueueConfig.NoWait,
				amqp.Table(router.QueueConfig.Args),
			)
			if err != nil {
				return err
			}

			deliveries, err := worker.channel.Consume(
				string(ports.QueueEmailProcessing),
				"consumer.email.processing",
				false, false, false, false, nil,
			)

			if err != nil {
				return err
			}

			ctx, cancel := context.WithCancel(context.Background())
			worker.cancel = cancel

			go worker.process(ctx, deliveries)

			return nil
		},
		OnStop: func(ctx context.Context) error {
			if worker.cancel != nil {
				worker.cancel()
			}
			if worker.channel != nil {
				return worker.channel.Close()
			}

			return nil
		},
	})

	return nil
}

func (worker *Worker) process(ctx context.Context, deliveries <-chan amqp.Delivery) {
	for delivery := range deliveries {
		var payload domain.Email
		if err := json.Unmarshal(delivery.Body, &payload); err != nil {
			_ = worker.sendToDlqAndAck(ctx, delivery, ports.NewWorkerPayloadInvalid(err))
			continue
		}

		email, err := worker.emailStorage.SelectByID(ctx, payload.ID)
		if err != nil {
			if ports.IsCode(err, ports.CodeDatabaseRowsNotFound) {
				_ = worker.sendToDlqAndAck(ctx, delivery, err)
				continue
			}

			_ = delivery.Nack(false, false)
			continue
		}

		if email.Attempts >= constants.MaxRetry {
			if err := worker.emailStorage.MarkStatusByIDs(ctx, email.ID, domain.EmailStatusFailed); err != nil {
				_ = delivery.Nack(false, true)
				continue
			}
			_ = worker.sendToDlqAndAck(ctx, delivery, ports.NewLimitAttemptsExceeded(nil))
			continue
		}

		if email.Status == domain.EmailStatusSuccess {
			_ = delivery.Ack(false)
			continue
		}

		if err := worker.emailStorage.MarkStatusByIDs(ctx, email.ID, domain.EmailStatusProcessing); err != nil {
			_ = delivery.Nack(false, false)
			continue
		}

		if err := worker.mailer.Send(email.To, email.From, email.Subject, string(email.Body)); err != nil {
			if err := worker.emailStorage.IncrementAttempts(ctx, email.ID); err != nil {
				_ = delivery.Nack(false, true)
				continue
			}

			// Nack sem requeue: envia para DLX configurada na queue principal,
			// que roteia para a retry queue. Após o TTL, a retry queue envia de volta
			// para a fila principal.
			_ = delivery.Nack(false, false)
			continue
		}

		if err := worker.emailStorage.MarkStatusByIDs(ctx, email.ID, domain.EmailStatusSuccess); err != nil {
			if err := worker.emailStorage.IncrementAttempts(ctx, email.ID); err != nil {
				_ = delivery.Nack(false, false)
				continue
			}
			_ = delivery.Nack(false, false)
			continue
		}

		_ = delivery.Ack(false)
	}
}

func (worker *Worker) sendToDlqAndAck(ctx context.Context, delivery amqp.Delivery, cause error) error {
	payload, err := json.Marshal(map[string]string{
		"raw":   string(delivery.Body),
		"error": cause.Error(),
	})
	if err != nil {
		return cause
	}

	cause = worker.publisher.Publish(ctx, ports.PublisherMessage{
		RoutingKey: ports.RoutingKeyEmailProcessingDlq,
		Payload:    payload,
	})
	if cause != nil {
		return delivery.Nack(false, true)
	}

	return delivery.Ack(false)
}
