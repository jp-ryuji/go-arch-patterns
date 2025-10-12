package renter

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

const MaxEmailLength = 50

type Email struct {
	value string
}

func NewEmail(email string) (*Email, error) {
	if err := validateEmail(email); err != nil {
		return nil, err
	}

	return &Email{value: strings.ToLower(email)}, nil
}

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

func (e *Email) String() string {
	return e.value
}

func (e *Email) Equals(other *Email) bool {
	if other == nil {
		return false
	}
	return e.value == other.value
}
