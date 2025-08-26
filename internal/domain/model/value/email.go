package value

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Email represents a validated email address value object
type Email struct {
	value string
}

const (
	// MaxEmailLength is the maximum allowed length for an email address
	MaxEmailLength = 50
)

// NewEmail creates a new Email value object after validating the email format
func NewEmail(email string) (*Email, error) {
	if err := validateEmail(email); err != nil {
		return nil, err
	}

	return &Email{value: strings.ToLower(email)}, nil
}

// validateEmail checks if the email meets the required format and constraints
func validateEmail(email string) error {
	if email == "" {
		return errors.New("email cannot be empty")
	}

	if len(email) > MaxEmailLength {
		return errors.New("email exceeds maximum length")
	}

	validate := validator.New()
	if err := validate.Var(email, "email"); err != nil {
		return errors.New("invalid email format")
	}

	return nil
}

// String returns the string representation of the email
func (e *Email) String() string {
	return e.value
}

// Equals checks if two Email objects are equal
func (e *Email) Equals(other *Email) bool {
	if other == nil {
		return false
	}
	return e.value == other.value
}
