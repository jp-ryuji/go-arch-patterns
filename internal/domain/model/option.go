package model

import (
	"time"

	"github.com/jp-ryuji/go-sample/internal/pkg/id"
)

// OptionRefs holds references to related entities for Option
type OptionRefs struct {
	Tenant *Tenant
}

type Option struct {
	ID        string
	TenantID  string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time

	Refs *OptionRefs
}

type Options []*Option

func NewOption(
	tenantID string,
	name string,
	t time.Time,
) *Option {
	return &Option{
		ID:        id.New(),
		TenantID:  tenantID,
		Name:      name,
		CreatedAt: t,
		UpdatedAt: t,

		Refs: nil,
	}
}
