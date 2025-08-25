package validator

import (
	"time"

	"github.com/go-playground/validator/v10"
)

var V = validator.New(validator.WithRequiredStructEnabled())

func init() {
	V.RegisterValidation("duration", validateDuration)
}

func validateDuration(fl validator.FieldLevel) bool {
	_, err := time.ParseDuration(fl.Field().String())
	return err == nil
}
