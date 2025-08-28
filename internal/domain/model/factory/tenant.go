package factory

import (
	"time"

	"github.com/jp-ryuji/go-sample/internal/domain/model"
	"github.com/jp-ryuji/go-sample/internal/pkg/id"
)

// NewTenant creates a new Tenant for testing purposes
func NewTenant() *model.Tenant {
	return model.NewTenant(
		"test-"+id.New(),
		time.Now(),
	)
}

// NewTenantWithCode creates a new Tenant with a specific code for testing purposes
func NewTenantWithCode(code string) *model.Tenant {
	return model.NewTenant(
		code,
		time.Now(),
	)
}
