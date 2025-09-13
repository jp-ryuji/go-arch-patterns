//go:build integration

package repository_test

import (
	"os"
	"testing"

	"github.com/jp-ryuji/go-arch-patterns/internal/infrastructure/postgres/repository/testutil"
)

func TestMain(m *testing.M) {
	// Set up the shared test environment
	if err := testutil.SetupTestEnvironment(); err != nil {
		os.Exit(1)
	}

	// Run the tests
	code := m.Run()

	// Clean up the shared test environment
	if err := testutil.TeardownTestEnvironment(); err != nil {
		os.Exit(1)
	}

	os.Exit(code)
}
