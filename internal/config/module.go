package config

import (
	"github.com/rickferrdev/mail-burrow/internal/config/env"
	"github.com/rickferrdev/mail-burrow/internal/config/rabbiting"
	"github.com/rickferrdev/mail-burrow/internal/config/sqlconn"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"config",
	env.Provide,
	rabbiting.Provide,
	sqlconn.Provide,
)
