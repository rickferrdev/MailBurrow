package emails

import (
	"context"
	"encoding/json"
	"net/mail"

	"github.com/google/uuid"
	"github.com/rickferrdev/mail-burrow/internal/app/domain"
	"github.com/rickferrdev/mail-burrow/internal/app/ports"
	"go.uber.org/fx"
)

var Provide = fx.Provide(New)

type Service struct {
	publisher ports.PublisherService
	storage   ports.EmailStorage
}

type Params struct {
	fx.In

	EmailStorage     ports.EmailStorage
	PublisherService ports.PublisherService
}

func New(params Params) (ports.EmailService, error) {
	service := Service{
		storage:   params.EmailStorage,
		publisher: params.PublisherService,
	}

	return &service, nil
}

func (service *Service) Sender(ctx context.Context, email ports.EmailDTO) (*ports.SenderEmailOutput, error) {
	if _, err := mail.ParseAddress(email.To); err != nil {
		return nil, ports.NewBadRequest(err)
	}
	if _, err := mail.ParseAddress(email.From); err != nil {
		return nil, ports.NewBadRequest(err)
	}

	domainEmail := domain.Email{
		ID:      uuid.New(),
		To:      email.To,
		From:    email.From,
		Subject: email.Subject,
		Body:    email.Body,
		Status:  domain.EmailStatusPending,
	}

	payload, err := json.Marshal(domainEmail)
	if err != nil {
		return nil, ports.NewInternal(err)
	}

	if err := service.storage.Insert(ctx, domainEmail); err != nil {
		if ports.IsCode(err, ports.CodeDatabaseDuplicate) {
			return nil, ports.NewInternal(err)
		}

		return nil, ports.NewInternal(err)
	}

	if err := service.publisher.PublishEmail(ctx, payload); err != nil {
		return nil, ports.NewInternal(err)
	}

	return &ports.SenderEmailOutput{ID: domainEmail.ID.String()}, nil
}

func (service *Service) Obtain(ctx context.Context, id string) (*ports.ObtainEmailOutput, error) {
	parse, err := uuid.Parse(id)
	if err != nil {
		return nil, ports.NewInvalidID(err)
	}

	email, err := service.storage.SelectByID(ctx, parse)
	if err != nil {
		if ports.IsCode(err, ports.CodeDatabaseRowsNotFound) {
			return nil, ports.NewNotFound(err)
		}

		return nil, ports.NewInternal(err)
	}

	obtainEmailDTO := ports.ObtainEmailOutput{
		ID:       email.ID.String(),
		Attempts: email.Attempts,
		Status:   email.Status,
	}

	return &obtainEmailDTO, nil
}
