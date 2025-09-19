package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/jp-ryuji/go-arch-patterns/internal/application/dto"
	"github.com/jp-ryuji/go-arch-patterns/internal/application/service"
	"github.com/jp-ryuji/go-arch-patterns/internal/domain/entity"
	mock_repository "github.com/jp-ryuji/go-arch-patterns/internal/domain/repository/mock"
	"github.com/jp-ryuji/go-arch-patterns/internal/infrastructure/postgres/entgen"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// setupTest creates a new mock controller and car service for testing
func setupTest(t *testing.T) (*gomock.Controller, *mock_repository.MockCarRepository, *mock_repository.MockOutboxRepository, *mock_repository.MockTransactionManager, service.CarService) {
	t.Helper()
	ctrl := gomock.NewController(t)
	mockCarRepo := mock_repository.NewMockCarRepository(ctrl)
	mockOutboxRepo := mock_repository.NewMockOutboxRepository(ctrl)
	mockTxManager := mock_repository.NewMockTransactionManager(ctrl)
	carService := service.NewCarService(mockCarRepo, mockOutboxRepo, mockTxManager)
	return ctrl, mockCarRepo, mockOutboxRepo, mockTxManager, carService
}

// TestCarService_Register_Success tests the successful registration of a car
func TestCarService_Register_Success(t *testing.T) {
	t.Parallel()

	// Setup
	ctrl, mockCarRepo, mockOutboxRepo, mockTxManager, carService := setupTest(t)
	defer ctrl.Finish()

	// Test data
	ctx := context.Background()
	registerInput := dto.RegisterCarInput{
		TenantID: "tenant-123",
		Model:    "Toyota Prius",
	}

	// Create a mock transaction
	mockTx := &entgen.Tx{}

	// Set up expectations for transaction management
	mockTxManager.EXPECT().BeginTx(ctx).Return(mockTx, nil)
	mockTxManager.EXPECT().CommitTx(ctx, mockTx).Return(nil)

	// Set up expectations for registering a car
	var createdCar *entity.Car
	mockCarRepo.EXPECT().CreateInTx(ctx, mockTx, gomock.Any()).DoAndReturn(
		func(ctx context.Context, tx *entgen.Tx, car *entity.Car) error {
			// Capture the car that was created for later verification
			createdCar = car

			// Verify that the car has the correct properties (similar to car_test.go)
			assert.Equal(t, registerInput.TenantID, car.TenantID)
			assert.Equal(t, registerInput.Model, car.Model)
			assert.NotEmpty(t, car.ID)
			assert.WithinDuration(t, time.Now(), car.CreatedAt, time.Second)
			assert.WithinDuration(t, time.Now(), car.UpdatedAt, time.Second)
			return nil
		},
	)

	// Set up expectations for outbox message creation
	mockOutboxRepo.EXPECT().CreateInTx(ctx, mockTx, gomock.Any()).Return(nil)

	// Execute - Register a new car using the service
	registeredCar, err := carService.Register(ctx, registerInput)
	assert.NoError(t, err)
	assert.NotNil(t, registeredCar)

	// Verify the returned car
	assert.Equal(t, registerInput.TenantID, registeredCar.TenantID)
	assert.Equal(t, registerInput.Model, registeredCar.Model)
	assert.NotEmpty(t, registeredCar.ID)
	assert.NotZero(t, registeredCar.CreatedAt)
	assert.NotZero(t, registeredCar.UpdatedAt)

	// Verify that the car sent to the repository matches what was returned
	assert.Equal(t, registeredCar, createdCar)
}

// TestCarService_Register_RepositoryError tests registration when repository returns an error
func TestCarService_Register_RepositoryError(t *testing.T) {
	t.Parallel()

	// Setup
	ctrl, mockCarRepo, _, mockTxManager, carService := setupTest(t)
	defer ctrl.Finish()

	// Test data
	ctx := context.Background()
	registerInput := dto.RegisterCarInput{
		TenantID: "tenant-123",
		Model:    "Toyota Prius",
	}

	// Create a mock transaction
	mockTx := &entgen.Tx{}

	// Set up expectations for transaction management
	mockTxManager.EXPECT().BeginTx(ctx).Return(mockTx, nil)
	mockTxManager.EXPECT().RollbackTx(ctx, mockTx).Return(nil)

	// Set up expectations for repository error
	expectedError := assert.AnError
	mockCarRepo.EXPECT().CreateInTx(ctx, mockTx, gomock.Any()).Return(expectedError)

	// Execute - Try to register a car when repository fails
	registeredCar, err := carService.Register(ctx, registerInput)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), expectedError.Error())
	assert.Nil(t, registeredCar)
}

// TestCarService_GetByID_Success tests successful retrieval of a car by ID
func TestCarService_GetByID_Success(t *testing.T) {
	t.Parallel()

	// Setup
	ctrl, mockCarRepo, _, _, carService := setupTest(t)
	defer ctrl.Finish()

	// Test data
	ctx := context.Background()
	carID := "car-123"
	getInput := dto.GetCarByIDInput{
		ID: carID,
	}

	// Create expected car using the factory method (similar to car_test.go)
	now := time.Now()
	expectedCar := entity.NewCar("tenant-123", "Toyota Prius", now)

	// Set up expectations for retrieving the car
	mockCarRepo.EXPECT().GetByID(ctx, carID).Return(expectedCar, nil)

	// Execute - Retrieve the car using the service
	retrievedCar, err := carService.GetByID(ctx, getInput)
	assert.NoError(t, err)
	assert.NotNil(t, retrievedCar)

	// Verify the returned car
	assert.Equal(t, expectedCar.ID, retrievedCar.ID)
	assert.Equal(t, expectedCar.TenantID, retrievedCar.TenantID)
	assert.Equal(t, expectedCar.Model, retrievedCar.Model)
}

// TestCarService_GetByID_NotFound tests retrieval when car doesn't exist
func TestCarService_GetByID_NotFound(t *testing.T) {
	t.Parallel()

	// Setup
	ctrl, mockCarRepo, _, _, carService := setupTest(t)
	defer ctrl.Finish()

	// Test data
	ctx := context.Background()
	carID := "non-existent-car"
	getInput := dto.GetCarByIDInput{
		ID: carID,
	}

	// Set up expectations for not found error
	expectedError := assert.AnError
	mockCarRepo.EXPECT().GetByID(ctx, carID).Return(nil, expectedError)

	// Execute - Try to retrieve a non-existent car
	retrievedCar, err := carService.GetByID(ctx, getInput)
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, retrievedCar)
}

// TestCarService_GetByIDWithTenant_Success tests successful retrieval of a car with tenant by ID
func TestCarService_GetByIDWithTenant_Success(t *testing.T) {
	t.Parallel()

	// Setup
	ctrl, mockCarRepo, _, _, carService := setupTest(t)
	defer ctrl.Finish()

	// Test data
	ctx := context.Background()
	carID := "car-123"
	getInput := dto.GetCarByIDInput{
		ID: carID,
	}

	// Create expected car with tenant using the factory method
	now := time.Now()
	expectedCar := entity.NewCar("tenant-123", "Toyota Prius", now)
	expectedTenant := entity.NewTenant("tenant-code", now)
	expectedCar.Refs = &entity.CarRefs{
		Tenant: expectedTenant,
	}

	// Set up expectations for retrieving the car with tenant
	mockCarRepo.EXPECT().GetByIDWithTenant(ctx, carID).Return(expectedCar, nil)

	// Execute - Retrieve the car with tenant using the service
	retrievedCar, err := carService.GetByIDWithTenant(ctx, getInput)
	assert.NoError(t, err)
	assert.NotNil(t, retrievedCar)

	// Verify the returned car
	assert.Equal(t, expectedCar.ID, retrievedCar.ID)
	assert.Equal(t, expectedCar.TenantID, retrievedCar.TenantID)
	assert.Equal(t, expectedCar.Model, retrievedCar.Model)
	assert.NotNil(t, retrievedCar.Refs)
	assert.NotNil(t, retrievedCar.Refs.Tenant)
	assert.Equal(t, expectedTenant.ID, retrievedCar.Refs.Tenant.ID)
	assert.Equal(t, expectedTenant.Code, retrievedCar.Refs.Tenant.Code)
}

// TestCarService_GetByIDWithTenant_NotFound tests retrieval when car with tenant doesn't exist
func TestCarService_GetByIDWithTenant_NotFound(t *testing.T) {
	t.Parallel()

	// Setup
	ctrl, mockCarRepo, _, _, carService := setupTest(t)
	defer ctrl.Finish()

	// Test data
	ctx := context.Background()
	carID := "non-existent-car"
	getInput := dto.GetCarByIDInput{
		ID: carID,
	}

	// Set up expectations for not found error
	expectedError := assert.AnError
	mockCarRepo.EXPECT().GetByIDWithTenant(ctx, carID).Return(nil, expectedError)

	// Execute - Try to retrieve a non-existent car with tenant
	retrievedCar, err := carService.GetByIDWithTenant(ctx, getInput)
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, retrievedCar)
}

// TestCarService_Register_Validation tests validation failures for Register
func TestCarService_Register_Validation(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input   dto.RegisterCarInput
		wantErr string
	}{
		"empty tenant ID": {
			input: dto.RegisterCarInput{
				TenantID: "", // Missing required field
				Model:    "Toyota Prius",
			},
			wantErr: "validation failed",
		},
		"empty model": {
			input: dto.RegisterCarInput{
				TenantID: "tenant-123",
				Model:    "", // Missing required field
			},
			wantErr: "validation failed",
		},
		"both fields empty": {
			input: dto.RegisterCarInput{
				TenantID: "", // Missing required field
				Model:    "", // Missing required field
			},
			wantErr: "validation failed",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Setup
			ctrl, _, _, _, carService := setupTest(t)
			defer ctrl.Finish()

			// Test data
			ctx := context.Background()

			// Execute - Try to register a car with invalid input
			registeredCar, err := carService.Register(ctx, tt.input)

			// Assert
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
			assert.Nil(t, registeredCar)
		})
	}
}

// TestCarService_GetByID_Validation tests validation failures for GetByID
func TestCarService_GetByID_Validation(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input   dto.GetCarByIDInput
		wantErr string
	}{
		"empty ID": {
			input: dto.GetCarByIDInput{
				ID: "", // Missing required field
			},
			wantErr: "validation failed",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Setup
			ctrl, _, _, _, carService := setupTest(t)
			defer ctrl.Finish()

			// Test data
			ctx := context.Background()

			// Execute - Try to retrieve a car with invalid input
			retrievedCar, err := carService.GetByID(ctx, tt.input)

			// Assert
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
			assert.Nil(t, retrievedCar)
		})
	}
}
