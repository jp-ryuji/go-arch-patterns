//go:build integration

package dbmodel_test

import (
	"testing"

	"github.com/jp-ryuji/go-sample/internal/domain/model/factory"
	"github.com/jp-ryuji/go-sample/internal/infrastructure/postgres/dbmodel"
	"github.com/jp-ryuji/go-sample/internal/infrastructure/postgres/repository/testutil"
	"github.com/stretchr/testify/require"
)

func TestTenantToDomain(t *testing.T) {
	testutil.SkipIfShort(t)

	db := testutil.DBClient.DB

	// Create a tenant
	tenant := factory.NewTenantWithCode("test-tenant-to-domain")
	tenantDB := dbmodel.FromDomainTenant(tenant)
	err := db.Create(tenantDB).Error
	require.NoError(t, err)

	// Fetch the tenant from DB
	var fetchedTenantDB dbmodel.Tenant
	err = db.Where("id = ?", tenant.ID).First(&fetchedTenantDB).Error
	require.NoError(t, err)

	// Convert to domain model
	domainTenant := fetchedTenantDB.ToDomain()
	require.Equal(t, tenant.ID, domainTenant.ID)
	require.Equal(t, tenant.Code, domainTenant.Code)
	require.Equal(t, tenant.CreatedAt.Unix(), domainTenant.CreatedAt.Unix())
	require.Equal(t, tenant.UpdatedAt.Unix(), domainTenant.UpdatedAt.Unix())
}
