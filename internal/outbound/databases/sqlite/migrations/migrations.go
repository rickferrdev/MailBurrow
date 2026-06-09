package migrations

import (
	"context"

	"github.com/rickferrdev/mail-burrow/internal/outbound/databases/sqlite/schemas"
	"github.com/rickferrdev/mail-burrow/internal/outbound/databases/sqlite/sqlutils"
	"github.com/uptrace/bun"
	"go.uber.org/fx"
)

var Invoke = fx.Invoke(func(lc fx.Lifecycle, db *bun.DB) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			_, err := db.NewCreateTable().
				Model((*schemas.Email)(nil)).
				IfNotExists().
				Exec(ctx)

			if err != nil {
				return sqlutils.NewError(err)
			}

			return nil
		},
	})
})
