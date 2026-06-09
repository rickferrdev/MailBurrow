package services

import (
	"github.com/rickferrdev/mail-burrow/internal/app/services/emails"
	"github.com/rickferrdev/mail-burrow/internal/app/services/publisher"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"services",
	emails.Provide,
	publisher.Provide,
)
