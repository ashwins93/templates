package utils

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/jaevor/go-nanoid"
)

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

var validate *validator.Validate
var GenID func() string

func init() {
	var err error
	validate = validator.New()
	GenID, err = nanoid.Standard(21)
	if err != nil {
		log.Fatal(err)
	}
}

func ValidateStruct(s interface{}) []*ErrorResponse {
	var errors []*ErrorResponse
	err := validate.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}
