package main

import (
	"github.com/rickferrdev/mail-burrow/internal/api"
	"github.com/rickferrdev/mail-burrow/internal/app/services"
	"github.com/rickferrdev/mail-burrow/internal/config"
	"github.com/rickferrdev/mail-burrow/internal/infra"
	"github.com/rickferrdev/mail-burrow/internal/outbound"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		config.Module,
		infra.Module,
		outbound.Module,
		services.Module,
		api.Module,
	).Run()
}
