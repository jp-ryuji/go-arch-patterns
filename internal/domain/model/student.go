package model

import (
	"time"

	"github.com/jp-ryuji/go-sample/internal/domain/model/value"
	"github.com/jp-ryuji/go-sample/internal/pkg/id"
)

type Student struct {
	ID               string
	AuthUID          string
	Email            *value.Email
	DisplayName      string
	FirstName        string
	LastName         string
	FirstNameKana    string
	LastNameKana     string
	ProfileImagePath string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type Students []*Student

func NewStudent(
	authUID string,
	email string,
	displayName string,
	firstName string,
	lastName string,
	firstNameKana string,
	lastNameKana string,
	profileImagePath string,
	t time.Time,
) (*Student, error) {
	emailVO, err := value.NewEmail(email)
	if err != nil {
		return nil, err
	}

	return &Student{
		ID:               id.New(),
		AuthUID:          authUID,
		Email:            emailVO,
		DisplayName:      displayName,
		FirstName:        firstName,
		LastName:         lastName,
		FirstNameKana:    firstNameKana,
		LastNameKana:     lastNameKana,
		ProfileImagePath: profileImagePath,
		CreatedAt:        t,
		UpdatedAt:        t,
	}, nil
}
