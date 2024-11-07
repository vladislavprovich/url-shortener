package validator

import "github.com/go-playground/validator/v10"

func Validate(i interface{}) error {
	validator := validator.New()
	return validator.Struct(i)
}
