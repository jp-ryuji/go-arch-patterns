package factory

import (
	"time"

	"github.com/jp-ryuji/go-arch-patterns/internal/domain/entity"
	"github.com/jp-ryuji/go-arch-patterns/internal/pkg/id"
)

// NewTenant creates a new Tenant for testing purposes
func NewTenant() *entity.Tenant {
	return entity.NewTenant(
		"test-"+id.New(),
		time.Now(),
	)
}

// NewTenantWithCode creates a new Tenant with a specific code for testing purposes
func NewTenantWithCode(code string) *entity.Tenant {
	return entity.NewTenant(
		code,
		time.Now(),
	)
}
