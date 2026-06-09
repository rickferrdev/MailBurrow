package publisher

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/rickferrdev/mail-burrow/internal/app/ports"
	"go.uber.org/fx"
)

var Provide = fx.Provide(New)

type Service struct {
	publisher ports.QueuePublisher
}

type Params struct {
	fx.In

	QueuePublisher ports.QueuePublisher
}

func New(params Params) (ports.PublisherService, error) {
	service := Service{
		publisher: params.QueuePublisher,
	}

	return &service, nil
}

func (service *Service) PublishEmail(ctx context.Context, payload []byte) error {
	if !json.Valid(payload) {
		return ports.NewInternal(errors.New("payload for sending the publisher email must be a valid JSON"))
	}

	if err := service.publisher.Publish(ctx, ports.PublisherMessage{
		RoutingKey: ports.RoutingKeyEmailProcessing,
		Payload:    payload,
	}); err != nil {
		return ports.NewRabbitMQPublish(err)
	}

	return nil
}
