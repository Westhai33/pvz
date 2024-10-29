package validation

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func InitValidator() {
	validate = validator.New()
}

func ValidateStruct(s interface{}) error {
	if validate == nil {
		InitValidator()
	}
	return validate.Struct(s)
}
