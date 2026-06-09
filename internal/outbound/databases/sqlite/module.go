package sqlite

import (
	"github.com/rickferrdev/mail-burrow/internal/outbound/databases/sqlite/handlers/email"
	"github.com/rickferrdev/mail-burrow/internal/outbound/databases/sqlite/migrations"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"sqlite-database",
	email.Provide,
	migrations.Invoke,
)
