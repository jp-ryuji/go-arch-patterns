package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/jp-ryuji/go-arch-patterns/internal/application/input"
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

// TestCarService_Create_Success tests the successful creation of a car
func TestCarService_Create_Success(t *testing.T) {
	t.Parallel()

	// Setup
	ctrl, mockCarRepo, mockOutboxRepo, mockTxManager, carService := setupTest(t)
	defer ctrl.Finish()

	// Test data
	ctx := context.Background()
	registerInput := input.CreateCar{
		TenantID: "tenant-123",
		Model:    "Toyota Prius",
	}

	// Create a mock transaction
	mockTx := &entgen.Tx{}

	// Set up expectations for transaction management
	mockTxManager.EXPECT().BeginTx(ctx).Return(mockTx, nil)
	mockTxManager.EXPECT().CommitTx(ctx, mockTx).Return(nil)

	// Set up expectations for creating a car
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

	// Execute - Create a new car using the service
	createdCarOutput, err := carService.Create(ctx, registerInput)
	assert.NoError(t, err)
	assert.NotNil(t, createdCarOutput)

	// Verify the returned car DTO
	assert.Equal(t, registerInput.TenantID, createdCarOutput.TenantID)
	assert.Equal(t, registerInput.Model, createdCarOutput.Model)
	assert.NotEmpty(t, createdCarOutput.ID)
	assert.NotZero(t, createdCarOutput.CreatedAt)

	// Verify that the car sent to the repository matches what was used to create the DTO
	assert.Equal(t, createdCarOutput.ID, createdCar.ID)
	assert.Equal(t, createdCarOutput.TenantID, createdCar.TenantID)
	assert.Equal(t, createdCarOutput.Model, createdCar.Model)
}

// TestCarService_Create_RepositoryError tests creation when repository returns an error
func TestCarService_Create_RepositoryError(t *testing.T) {
	t.Parallel()

	// Setup
	ctrl, mockCarRepo, _, mockTxManager, carService := setupTest(t)
	defer ctrl.Finish()

	// Test data
	ctx := context.Background()
	registerInput := input.CreateCar{
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

	// Execute - Try to create a car when repository fails
	createdCar, err := carService.Create(ctx, registerInput)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), expectedError.Error())
	assert.Nil(t, createdCar)
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
	getInput := input.GetCarByID{
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

	// Verify the returned car DTO
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
	getInput := input.GetCarByID{
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
	getInput := input.GetCarByID{
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

	// Verify the returned car DTO
	assert.Equal(t, expectedCar.ID, retrievedCar.ID)
	assert.Equal(t, expectedCar.TenantID, retrievedCar.TenantID)
	assert.Equal(t, expectedCar.Model, retrievedCar.Model)
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
	getInput := input.GetCarByID{
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

// TestCarService_List_Success tests successful listing of cars
func TestCarService_List_Success(t *testing.T) {
	t.Parallel()

	// Setup
	ctrl, mockCarRepo, _, _, carService := setupTest(t)
	defer ctrl.Finish()

	// Test data
	ctx := context.Background()
	tenantID := "tenant-123"
	listInput := input.ListCars{
		TenantID: tenantID,
		PageSize: 10,
	}

	// Create expected cars using the factory method
	now := time.Now()
	expectedCars := []*entity.Car{
		entity.NewCar(tenantID, "Toyota Prius", now),
		entity.NewCar(tenantID, "Honda Civic", now),
	}
	expectedNextPageToken := "next-page-token"
	expectedTotalCount := int32(2)

	// Set up expectations for listing cars
	mockCarRepo.EXPECT().ListByTenant(ctx, tenantID, 10, 0).Return(expectedCars, expectedNextPageToken, expectedTotalCount, nil)

	// Execute - List cars using the service
	listOutput, err := carService.List(ctx, listInput)
	assert.NoError(t, err)
	assert.NotNil(t, listOutput)

	// Verify the returned list DTO
	assert.Equal(t, expectedNextPageToken, listOutput.NextPageToken)
	assert.Equal(t, expectedTotalCount, listOutput.TotalCount)
	assert.Len(t, listOutput.Cars, 2)
	assert.Equal(t, expectedCars[0].ID, listOutput.Cars[0].ID)
	assert.Equal(t, expectedCars[0].Model, listOutput.Cars[0].Model)
	assert.Equal(t, expectedCars[1].ID, listOutput.Cars[1].ID)
	assert.Equal(t, expectedCars[1].Model, listOutput.Cars[1].Model)
}

// TestCarService_Create_Validation tests validation failures for Create
func TestCarService_Register_Validation(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input   input.CreateCar
		wantErr string
	}{
		"empty tenant ID": {
			input: input.CreateCar{
				TenantID: "", // Missing required field
				Model:    "Toyota Prius",
			},
			wantErr: "validation failed",
		},
		"empty model": {
			input: input.CreateCar{
				TenantID: "tenant-123",
				Model:    "", // Missing required field
			},
			wantErr: "validation failed",
		},
		"both fields empty": {
			input: input.CreateCar{
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

			// Execute - Try to create a car with invalid input
			createdCar, err := carService.Create(ctx, tt.input)

			// Assert
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
			assert.Nil(t, createdCar)
		})
	}
}

// TestCarService_GetByID_Validation tests validation failures for GetByID
func TestCarService_GetByID_Validation(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input   input.GetCarByID
		wantErr string
	}{
		"empty ID": {
			input: input.GetCarByID{
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
