package model

import (
	"time"

	"github.com/jp-ryuji/go-sample/internal/pkg/id"
)

type NoteOwnable struct {
	ID         string
	SchoolID   string
	OwnerID    string
	OwnerModel NoteOwnerModel
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type NoteOwnables []*NoteOwnables

func NewOwnable(
	schoolID string,
	ownerID string,
	ownerModel NoteOwnerModel,
	t time.Time,
) *NoteOwnable {
	return &NoteOwnable{
		ID:         id.New(),
		SchoolID:   schoolID,
		OwnerID:    ownerID,
		OwnerModel: ownerModel,
		CreatedAt:  t,
		UpdatedAt:  t,
	}
}

type NoteOwnerModel string

const (
	NoteOwnerModelTeacher NoteOwnerModel = "teacher"
	NoteOwnerModelStudent NoteOwnerModel = "student"
)
