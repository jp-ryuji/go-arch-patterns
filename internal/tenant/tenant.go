package tenant

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type Tenant struct {
	ID        string
	Code      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func New(code string, createdAt time.Time) *Tenant {
	return &Tenant{
		ID:        ulid.Make().String(),
		Code:      code,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	}
}

func (t *Tenant) WithID(id string) *Tenant {
	t.ID = id
	return t
}
