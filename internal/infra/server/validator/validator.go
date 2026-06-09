package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/rickferrdev/mail-burrow/internal/app/ports"
	"go.uber.org/fx"
)

var Provide = fx.Provide(New)

type Validator struct {
	validate *validator.Validate
}

func New() (ports.Validator, error) {
	return &Validator{
		validate: validator.New(),
	}, nil
}

func (val *Validator) Validate(out any) error {
	return val.validate.Struct(out)
}
