//go:build integration

package testutil

import (
	"context"
	"testing"

	"github.com/jp-ryuji/go-arch-patterns/internal/domain/entity"
	"github.com/jp-ryuji/go-arch-patterns/internal/domain/entity/factory"
	"github.com/stretchr/testify/require"
)

// CreateTestTenant creates a tenant for testing and saves it to the database.
func CreateTestTenant(t *testing.T, code string) *entity.Tenant {
	t.Helper()

	tenant := factory.NewTenantWithCode(code)

	// Create tenant in database using Ent
	tenantDB, err := DBClient.Tenant.
		Create().
		SetID(tenant.ID).
		SetCode(tenant.Code).
		Save(context.Background())

	require.NoError(t, err)

	// Convert back to domain model
	return &entity.Tenant{
		ID:   tenantDB.ID,
		Code: tenantDB.Code,
	}
}

// CreateRandomTestTenant creates a tenant with a random code for testing and saves it to the database.
func CreateRandomTestTenant(t *testing.T) *entity.Tenant {
	t.Helper()

	tenant := factory.NewTenant()

	// Create tenant in database using Ent
	tenantDB, err := DBClient.Tenant.
		Create().
		SetID(tenant.ID).
		SetCode(tenant.Code).
		Save(context.Background())

	require.NoError(t, err)

	// Convert back to domain model
	return &entity.Tenant{
		ID:   tenantDB.ID,
		Code: tenantDB.Code,
	}
}
