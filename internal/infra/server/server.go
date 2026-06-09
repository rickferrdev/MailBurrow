package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/rickferrdev/mail-burrow/internal/app/ports"
	"github.com/rickferrdev/mail-burrow/internal/config/env"
	"go.uber.org/fx"
)

var Provide = fx.Provide(New)
var Invoke = fx.Invoke(
	func(life fx.Lifecycle, env *env.Environment, log *slog.Logger, app *fiber.App) error {
		life.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				go func() {
					addr := fmt.Sprintf("%v:%v", env.ServerHost, env.ServerPort)
					if err := app.Listen(addr); err != nil {
						log.Error("error starting the HTTP server", "error", err.Error())
					}
				}()

				return nil
			},
			OnStop: app.ShutdownWithContext,
		})

		return nil
	},
)

func New(validator ports.Validator) (*fiber.App, fiber.Router, error) {
	app := fiber.New(fiber.Config{
		StrictRouting:   true,
		CaseSensitive:   true,
		AppName:         "mailburrow",
		ErrorHandler:    ErrorHandler,
		StructValidator: validator,
	})

	app.Use(logger.New())

	return app, app.Group("/api/v1"), nil
}

func ErrorHandler(c fiber.Ctx, err error) error {
	if err == nil {
		return nil
	}

	var portErr *ports.Error
	if errors.As(err, &portErr) {
		return c.Status(portErr.Status).JSON(ports.Error{
			Message: portErr.Message,
			Code:    portErr.Code,
			Status:  portErr.Status,
		})
	}

	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		code := ports.CodeHTTPRequest
		message := ports.Message(fiberErr.Message)

		if fiberErr.Code == fiber.StatusNotFound {
			code = ports.CodeNotFound
			message = ports.MessageNotFound
		}

		return c.Status(fiberErr.Code).JSON(ports.Error{
			Message: message,
			Code:    code,
			Status:  fiberErr.Code,
		})
	}

	internal := ports.NewInternal(err)

	return c.Status(internal.Status).JSON(ports.Error{
		Message: internal.Message,
		Code:    internal.Code,
		Status:  internal.Status,
	})
}
