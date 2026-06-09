package outbound

import (
	"github.com/rickferrdev/mail-burrow/internal/outbound/databases"
	"github.com/rickferrdev/mail-burrow/internal/outbound/queues"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"outbound",
	databases.Module,
	queues.Module,
)
