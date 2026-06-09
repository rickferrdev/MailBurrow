package email

import (
	"context"

	"github.com/google/uuid"
	"github.com/rickferrdev/mail-burrow/internal/app/domain"
	"github.com/rickferrdev/mail-burrow/internal/app/ports"
	"github.com/rickferrdev/mail-burrow/internal/outbound/databases/sqlite/schemas"
	"github.com/rickferrdev/mail-burrow/internal/outbound/databases/sqlite/sqlutils"
	"github.com/uptrace/bun"
	"go.uber.org/fx"
)

var Provide = fx.Provide(New)

type Storage struct {
	database *bun.DB
}

type Params struct {
	fx.In

	Database *bun.DB
}

func New(params Params) (ports.EmailStorage, error) {
	storage := Storage{
		database: params.Database,
	}

	return &storage, nil
}

func (storage *Storage) Insert(ctx context.Context, email domain.Email) error {
	return storage.database.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		schema, err := schemas.FromEmailSchema(email)
		if err != nil {
			return ports.NewInvalidID(err)
		}

		res, err := tx.NewInsert().Model(schema).Exec(ctx)
		if err != nil {
			return sqlutils.NewError(err)
		}
		if err := sqlutils.NewResultError(res); err != nil {
			return err
		}

		return nil
	})
}

func (storage *Storage) SelectByID(ctx context.Context, id uuid.UUID) (*domain.Email, error) {
	var schema schemas.Email
	if id == uuid.Nil {
		return nil, ports.NewInvalidID(nil)
	}

	if err := storage.database.NewSelect().
		Model(&schema).
		Where("id = ?", id.String()).
		Scan(ctx); err != nil {
		return nil, sqlutils.NewError(err)
	}

	toDomain, err := schema.ToEmailDomain()
	if err != nil {
		return nil, err
	}

	return toDomain, err
}

func (storage *Storage) MarkStatusByIDs(ctx context.Context, id uuid.UUID, status domain.EmailStatus) error {
	return storage.database.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		res, err := tx.NewUpdate().
			Model((*schemas.Email)(nil)).
			Where("id = ?", id.String()).
			Set("status = ?", status).
			Exec(ctx)
		if err != nil {
			return sqlutils.NewError(err)
		}
		if err := sqlutils.NewResultError(res); err != nil {
			return err
		}

		return err
	})
}

func (storage *Storage) IncrementAttempts(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return ports.NewInvalidID(nil)
	}

	return storage.database.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		res, err := tx.NewUpdate().
			Model((*schemas.Email)(nil)).
			Where("id = ?", id.String()).
			Set("attempts = attempts + 1").
			Exec(ctx)
		if err != nil {
			return sqlutils.NewError(err)
		}
		if err := sqlutils.NewResultError(res); err != nil {
			return err
		}

		return nil
	})
}

func (storage *Storage) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return storage.database.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		res, err := tx.NewDelete().
			Model((*schemas.Email)(nil)).
			Where("id = ?", id.String()).
			Exec(ctx)
		if err != nil {
			return sqlutils.NewError(err)
		}
		if err := sqlutils.NewResultError(res); err != nil {
			return err
		}

		return err
	})
}
