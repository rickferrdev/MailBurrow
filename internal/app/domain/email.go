package domain

import (
	"time"

	"github.com/google/uuid"
)

type EmailStatus string

var (
	EmailStatusPending    EmailStatus = "pending"
	EmailStatusProcessing EmailStatus = "processing"
	EmailStatusSuccess    EmailStatus = "success"
	EmailStatusFailed     EmailStatus = "failed"
)

type Email struct {
	ID uuid.UUID

	To      string
	From    string
	Subject string
	Body    string

	Attempts int

	Status EmailStatus
	SentAt *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}
