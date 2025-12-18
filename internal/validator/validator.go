package validator

import (

	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/SamrudhiTM/user_api/internal/models"
)

// Create a single validator instance
var validate = validator.New()

// ValidateStruct validates a struct and returns formatted error messages
func ValidateStruct(s interface{}) map[string]string {
	err := validate.Struct(s)
	if err == nil {
		// Extra custom checks
		if req, ok := s.(models.CreateUserRequest); ok {
			errors := make(map[string]string)
			dob, parseErr := time.Parse("2006-01-02", req.Dob)
			if parseErr == nil {
				if dob.After(time.Now()) {
					errors["Dob"] = "Dob cannot be in the future"
					return errors
				}
			}
		}
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
		case "min":
			msg = fmt.Sprintf("%s must be at least %s characters", fieldName, err.Param())
		case "max":
			msg = fmt.Sprintf("%s must be at most %s characters", fieldName, err.Param())
		default:
			msg = fmt.Sprintf("%s is invalid", fieldName)
		}

		errors[fieldName] = msg
	}

	// ✅ Custom future DOB check
	if req, ok := s.(models.CreateUserRequest); ok {
		dob, parseErr := time.Parse("2006-01-02", req.Dob)
		if parseErr == nil && dob.After(time.Now()) {
			errors["Dob"] = "Dob cannot be in the future"
		}
	}

	return errors
}
