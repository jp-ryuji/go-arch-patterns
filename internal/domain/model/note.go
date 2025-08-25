package model

import (
	"time"

	"github.com/jp-ryuji/go-sample/internal/pkg/id"
)

type Note struct {
	ID        string
	SchoolID  string
	OwnableID string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Notes []*Note

func NewNote(
	schoolID string,
	ownableID string,
	name string,
	t time.Time,
) *Note {
	return &Note{
		ID:        id.New(),
		SchoolID:  schoolID,
		OwnableID: ownableID,
		Name:      name,
		CreatedAt: t,
		UpdatedAt: t,
	}
}
