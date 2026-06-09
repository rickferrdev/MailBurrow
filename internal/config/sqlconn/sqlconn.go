package sqlconn

import (
	"context"
	"database/sql"

	_ "modernc.org/sqlite"

	"github.com/rickferrdev/mail-burrow/internal/app/ports"
	"github.com/rickferrdev/mail-burrow/internal/config/env"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"go.uber.org/fx"
)

var Provide = fx.Provide(New)

type Params struct {
	fx.In

	Env  *env.Environment
	Life fx.Lifecycle
}

func New(params Params) (*bun.DB, error) {
	sqldb, err := sql.Open(sqliteshim.ShimName, params.Env.DatabaseUrl)
	if err != nil {
		return nil, ports.NewDatabaseConnection(err)
	}
	database := bun.NewDB(sqldb, sqlitedialect.New())

	params.Life.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			if database != nil {
				return database.Close()
			}

			return nil
		},
	})

	return database, nil
}
