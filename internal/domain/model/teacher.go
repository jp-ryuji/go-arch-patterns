package model

import (
	"time"

	"github.com/jp-ryuji/go-sample/internal/pkg/id"
)

type Teacher struct {
	ID        string
	SchoolID  string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time

	ReadonlyReference *struct {
		School *School
	}
}

type Teachers []*Teacher

func NewTeacher(
	schoolID string,
	name string,
	t time.Time,
) *Teacher {
	return &Teacher{
		ID:        id.New(),
		SchoolID:  schoolID,
		Name:      name,
		CreatedAt: t,
		UpdatedAt: t,

		ReadonlyReference: nil,
	}
}
