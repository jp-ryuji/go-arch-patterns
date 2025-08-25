package model

import (
	"time"

	"github.com/aarondl/null/v9"
	"github.com/jp-ryuji/go-sample/internal/pkg/id"
)

type SchoolStudent struct {
	ID        string
	SchoolID  string
	StudentID string
	CreatedAt time.Time
	UpdatedAt time.Time
	LeftAt    null.Time

	ReadonlyReference *struct {
		Student *Student
	}

	ProfileImageURL null.String
}

type SchoolStudents []*SchoolStudent

func NewSchoolStudent(
	schoolID string,
	studentID string,
	name string,
	firstName string,
	lastName string,
	firstNameKana string,
	lastNameKana string,
	profileImagePath string,
	t time.Time,
) *SchoolStudent {
	return &SchoolStudent{
		ID:        id.New(),
		SchoolID:  schoolID,
		StudentID: studentID,
		CreatedAt: t,
		UpdatedAt: t,
		LeftAt:    null.Time{},

		ReadonlyReference: nil,
		ProfileImageURL:   null.String{},
	}
}
