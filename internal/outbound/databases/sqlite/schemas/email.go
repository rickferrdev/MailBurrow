package schemas

import (
	"time"

	"github.com/google/uuid"
	"github.com/rickferrdev/mail-burrow/internal/app/domain"
	"github.com/rickferrdev/mail-burrow/internal/app/ports"
	"github.com/uptrace/bun"
)

type Email struct {
	bun.BaseModel `bun:"table:emails,alias:e"`

	ID uuid.UUID `bun:"id,pk,type:text"`

	To      string `bun:"email_to,notnull"`
	From    string `bun:"email_from,notnull"`
	Subject string `bun:"email_subject,notnull"`
	Body    string `bun:"email_body,notnull"`

	Attempts int `bun:"attempts"`

	Status domain.EmailStatus `bun:"status,notnull,default:'pending'"`
	SentAt *time.Time         `bun:"sent_at,nullzero"`

	CreatedAt time.Time `bun:"created_at,nullzero,default:current_timestamp"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,default:current_timestamp"`
}

func FromEmailSchema(email domain.Email) (*Email, error) {
	if err := uuid.Validate(email.ID.String()); err != nil {
		return nil, ports.NewInvalidID(err)
	}

	emailSchema := Email{
		ID:        email.ID,
		To:        email.To,
		From:      email.From,
		Subject:   email.Subject,
		Body:      email.Body,
		Attempts:  email.Attempts,
		Status:    email.Status,
		SentAt:    email.SentAt,
		CreatedAt: email.CreatedAt,
		UpdatedAt: email.UpdatedAt,
	}

	return &emailSchema, nil
}

func (email *Email) ToEmailDomain() (*domain.Email, error) {
	if err := uuid.Validate(email.ID.String()); err != nil {
		return nil, ports.NewInvalidID(err)
	}

	emailDomain := domain.Email{
		ID:        email.ID,
		To:        email.To,
		From:      email.From,
		Subject:   email.Subject,
		Body:      email.Body,
		Attempts:  email.Attempts,
		Status:    email.Status,
		SentAt:    email.SentAt,
		CreatedAt: email.CreatedAt,
		UpdatedAt: email.UpdatedAt,
	}

	return &emailDomain, nil
}
