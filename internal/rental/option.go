package rental

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type Options []*Option

type Option struct {
	ID        string
	TenantID  string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewOption(tenantID, name string) *Option {
	now := time.Now()
	return &Option{
		ID:        ulid.Make().String(),
		TenantID:  tenantID,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (o *Option) WithID(id string) *Option {
	o.ID = id
	return o
}
