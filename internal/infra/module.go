package infra

import (
	"github.com/rickferrdev/mail-burrow/internal/infra/log"
	"github.com/rickferrdev/mail-burrow/internal/infra/mailer"
	"github.com/rickferrdev/mail-burrow/internal/infra/server"
	"github.com/rickferrdev/mail-burrow/internal/infra/server/validator"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"infra",
	log.Provide,
	mailer.Provide,
	validator.Provide,
	server.Provide,
	server.Invoke,
)
