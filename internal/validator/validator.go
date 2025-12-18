package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// Create a single validator instance
var validate = validator.New()

// ValidateStruct validates a struct and returns formatted error messages
func ValidateStruct(s interface{}) map[string]string {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	errors := make(map[string]string)
	for _, err := range err.(validator.ValidationErrors) {
		fieldName := err.Field()
		var msg string

		switch err.Tag() {
		case "required":
			msg = fmt.Sprintf("%s is required", fieldName)
		case "datetime":
			msg = fmt.Sprintf("%s must be in YYYY-MM-DD format", fieldName)
		default:
			msg = fmt.Sprintf("%s is invalid", fieldName)
		}

		errors[fieldName] = msg
	}

	return errors
}
