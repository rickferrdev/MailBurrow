package databases

import (
	"github.com/rickferrdev/mail-burrow/internal/outbound/databases/sqlite"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"databases",
	sqlite.Module,
)
