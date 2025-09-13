//go:build integration

package testutil

import (
	"context"
	"testing"

	"github.com/jp-ryuji/go-ddd/internal/domain/model"
	"github.com/jp-ryuji/go-ddd/internal/domain/model/factory"
	"github.com/stretchr/testify/require"
)

// CreateTestTenant creates a tenant for testing and saves it to the database.
func CreateTestTenant(t *testing.T, code string) *model.Tenant {
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
	return &model.Tenant{
		ID:   tenantDB.ID,
		Code: tenantDB.Code,
	}
}

// CreateRandomTestTenant creates a tenant with a random code for testing and saves it to the database.
func CreateRandomTestTenant(t *testing.T) *model.Tenant {
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
	return &model.Tenant{
		ID:   tenantDB.ID,
		Code: tenantDB.Code,
	}
}
