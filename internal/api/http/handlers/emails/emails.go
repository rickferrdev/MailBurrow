package emails

import (
	"errors"

	"github.com/gofiber/fiber/v3"
	"github.com/rickferrdev/mail-burrow/internal/app/ports"
	"go.uber.org/fx"
)

var Invoke = fx.Invoke(New)

type Handler struct {
	service ports.EmailService
	router  fiber.Router
}

type Params struct {
	fx.In

	Service ports.EmailService
	Router  fiber.Router
}

type Interface interface {
	Publish(c fiber.Ctx) error
	Obtain(c fiber.Ctx) error
}

func New(params Params) (Interface, error) {
	handler := Handler{
		service: params.Service,
		router:  params.Router,
	}

	handler.router.Post("/emails/publish", handler.Publish)
	handler.router.Get("/emails", handler.Obtain)

	return &handler, nil
}

type RequestEmailPublishDTO struct {
	To      string `json:"to" validate:"required,email"`
	From    string `json:"from" validate:"required,email"`
	Body    string `json:"body" validate:"required,max=998"`
	Subject string `json:"subject" validate:"required,max=998"`
}

type ResponseEmailPublishDTO struct {
	ID string `json:"id"`
}

type RequestEmailObtainDTO struct {
	ID string `json:"id" validate:"required,uuid"`
}

type ResponseEmailObtainDTO struct {
	ID       string `json:"id"`
	Attempts int    `json:"attempts"`
	Status   string `json:"status"`
}

func (handler *Handler) Publish(c fiber.Ctx) error {
	var body RequestEmailPublishDTO
	if err := c.Bind().JSON(&body); err != nil {
		return ports.NewBadRequest(err)
	}

	output, err := handler.service.Sender(c.RequestCtx(), ports.EmailDTO{
		To:      body.To,
		From:    body.From,
		Body:    body.Body,
		Subject: body.Subject,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(ResponseEmailPublishDTO{
		ID: output.ID,
	})
}

func (handler *Handler) Obtain(c fiber.Ctx) error {
	id := c.Query("id")
	if id == "" {
		return ports.NewBadRequest(errors.New("should provide an ID"))
	}

	output, err := handler.service.Obtain(c.RequestCtx(), id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(ResponseEmailObtainDTO{
		ID:       output.ID,
		Attempts: output.Attempts,
		Status:   string(output.Status),
	})
}
