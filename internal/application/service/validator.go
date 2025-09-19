package service

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// validate is a package-level validator instance
var validate *validator.Validate

func init() {
	validate = validator.New()
}

// Validate performs validation on the provided struct
func Validate(s interface{}) error {
	if err := validate.Struct(s); err != nil {
		// You can customize the error handling here
		// For now, we'll just return the first error
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			if len(validationErrors) > 0 {
				// Return the first validation error with a descriptive message
				field := validationErrors[0].Field()
				tag := validationErrors[0].Tag()
				return fmt.Errorf("validation failed for field '%s': %s", field, tag)
			}
		}
		return err
	}
	return nil
}
