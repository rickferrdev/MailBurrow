package ports

import (
	"context"

	"github.com/rickferrdev/mail-burrow/internal/app/domain"
)

//go:generate mockgen -source=$GOFILE -destination=../../tests/mocks/services.go -package=mocks
type EmailDTO struct {
	To      string
	From    string
	Body    string
	Subject string
}

type SenderEmailOutput struct {
	ID string
}

type ObtainEmailOutput struct {
	ID       string
	Attempts int
	Status   domain.EmailStatus
}

type EmailService interface {
	Sender(ctx context.Context, email EmailDTO) (*SenderEmailOutput, error)
	Obtain(ctx context.Context, id string) (*ObtainEmailOutput, error)
}

type PublisherService interface {
	PublishEmail(ctx context.Context, payload []byte) error
}
