package processor_test

import (
	"testing"
	"time"

	mock_repository "github.com/jp-ryuji/go-arch-patterns/internal/domain/repository/mock"
	"github.com/jp-ryuji/go-arch-patterns/internal/infrastructure/opensearch/processor"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestOutboxProcessor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mocks
	mockOutboxRepo := mock_repository.NewMockOutboxRepository(ctrl)

	// Test that we can create a processor
	// Note: We're not testing the actual creation because it would try to connect to OpenSearch
	// In a real test, we would use dependency injection to mock the OpenSearch client
	p, err := processor.NewOutboxProcessor(mockOutboxRepo, 10, 3, time.Second)
	assert.NoError(t, err)
	assert.NotNil(t, p)

	// Test that processor was created successfully
	// We can't directly access the fields because they're not exported
	// But we can test that the function doesn't return an error
	assert.NotNil(t, p)
}
