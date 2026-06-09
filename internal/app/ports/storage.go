package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/rickferrdev/mail-burrow/internal/app/domain"
)

//go:generate mockgen -source=$GOFILE -destination=../../tests/mocks/storage.go -package=mocks
type EmailStorage interface {
	Insert(ctx context.Context, email domain.Email) error
	SelectByID(ctx context.Context, id uuid.UUID) (*domain.Email, error)
	MarkStatusByIDs(ctx context.Context, id uuid.UUID, status domain.EmailStatus) error
	IncrementAttempts(ctx context.Context, id uuid.UUID) error
	DeleteByID(ctx context.Context, id uuid.UUID) error
}
