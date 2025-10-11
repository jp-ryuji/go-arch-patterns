//go:build integration

package testutil

import (
	"context"
	"testing"

	tenant "github.com/jp-ryuji/go-arch-patterns/internal/tenant"
	"github.com/stretchr/testify/require"
)

func CreateTestTenant(t *testing.T, code string) *tenant.Tenant {
	t.Helper()

	tenantModel := tenant.NewWithCode(code)

	tenantDB, err := DBClient.Tenant.
		Create().
		SetID(tenantModel.ID).
		SetCode(tenantModel.Code).
		Save(context.Background())

	require.NoError(t, err)

	return &tenant.Tenant{
		ID:        tenantDB.ID,
		Code:      tenantDB.Code,
		CreatedAt: tenantDB.CreatedAt,
		UpdatedAt: tenantDB.UpdatedAt,
	}
}

func CreateRandomTestTenant(t *testing.T) *tenant.Tenant {
	t.Helper()

	tenantModel := tenant.NewFixture()

	tenantDB, err := DBClient.Tenant.
		Create().
		SetID(tenantModel.ID).
		SetCode(tenantModel.Code).
		Save(context.Background())

	require.NoError(t, err)

	return &tenant.Tenant{
		ID:        tenantDB.ID,
		Code:      tenantDB.Code,
		CreatedAt: tenantDB.CreatedAt,
		UpdatedAt: tenantDB.UpdatedAt,
	}
}
