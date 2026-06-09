package handlers

import (
	"github.com/rickferrdev/mail-burrow/internal/api/http/handlers/emails"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"handlers",
	emails.Invoke,
)
