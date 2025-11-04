package validatorutil

import (
	"time"

	"github.com/go-playground/validator/v10"
)

const durationValidatorTag string = "duration"

func durationValidator(fl validator.FieldLevel) bool {
	value := fl.Field().String();

	_, err := time.ParseDuration(value);
	return err == nil;
}

