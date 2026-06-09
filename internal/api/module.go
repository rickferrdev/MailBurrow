package api

import (
	"github.com/rickferrdev/mail-burrow/internal/api/http/handlers"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"api",
	handlers.Module,
)
