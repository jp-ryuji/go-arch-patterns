//go:build integration

package testutil

import (
	"testing"

	"github.com/jp-ryuji/go-sample/internal/domain/model"
	"github.com/jp-ryuji/go-sample/internal/domain/model/factory"
	"github.com/jp-ryuji/go-sample/internal/infrastructure/postgres/dbmodel"
	"github.com/stretchr/testify/require"
)

// CreateTestTenant creates a tenant for testing and saves it to the database.
func CreateTestTenant(t *testing.T, code string) *model.Tenant {
	t.Helper()

	tenant := factory.NewTenantWithCode(code)

	// Convert to DB model and create in database
	tenantDB := dbmodel.FromDomainTenant(tenant)
	err := DBClient.DB.Model(&dbmodel.Tenant{}).Create(tenantDB).Error

	require.NoError(t, err)
	return tenant
}

// CreateRandomTestTenant creates a tenant with a random code for testing and saves it to the database.
func CreateRandomTestTenant(t *testing.T) *model.Tenant {
	t.Helper()

	tenant := factory.NewTenant()

	// Convert to DB model and create in database
	tenantDB := dbmodel.FromDomainTenant(tenant)
	err := DBClient.DB.Model(&dbmodel.Tenant{}).Create(tenantDB).Error

	require.NoError(t, err)
	return tenant
}
