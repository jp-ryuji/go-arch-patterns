package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/jp-ryuji/go-ddd/internal/domain/model"
	mock_repository "github.com/jp-ryuji/go-ddd/internal/domain/repository/mock"
	"github.com/jp-ryuji/go-ddd/internal/infrastructure/usecase"
	"github.com/jp-ryuji/go-ddd/internal/infrastructure/usecase/input"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// setupTest creates a new mock controller and car usecase for testing
func setupTest(t *testing.T) (*gomock.Controller, *mock_repository.MockCarRepository, usecase.CarUsecase) {
	t.Helper()
	ctrl := gomock.NewController(t)
	mockCarRepo := mock_repository.NewMockCarRepository(ctrl)
	carUsecase := usecase.NewCarUsecase(mockCarRepo)
	return ctrl, mockCarRepo, carUsecase
}

// TestCarUsecase_Register_Success tests the successful registration of a car
func TestCarUsecase_Register_Success(t *testing.T) {
	t.Parallel()

	// Setup
	ctrl, mockCarRepo, carUsecase := setupTest(t)
	defer ctrl.Finish()

	// Test data
	ctx := context.Background()
	registerInput := input.RegisterCarInput{
		TenantID: "tenant-123",
		Model:    "Toyota Prius",
	}

	// Set up expectations for registering a car
	var createdCar *model.Car
	mockCarRepo.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(
		func(ctx context.Context, car *model.Car) error {
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

	// Execute - Register a new car using the usecase
	registeredCar, err := carUsecase.Register(ctx, registerInput)
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

// TestCarUsecase_Register_RepositoryError tests registration when repository returns an error
func TestCarUsecase_Register_RepositoryError(t *testing.T) {
	t.Parallel()

	// Setup
	ctrl, mockCarRepo, carUsecase := setupTest(t)
	defer ctrl.Finish()

	// Test data
	ctx := context.Background()
	registerInput := input.RegisterCarInput{
		TenantID: "tenant-123",
		Model:    "Toyota Prius",
	}

	// Set up expectations for repository error
	expectedError := assert.AnError
	mockCarRepo.EXPECT().Create(ctx, gomock.Any()).Return(expectedError)

	// Execute - Try to register a car when repository fails
	registeredCar, err := carUsecase.Register(ctx, registerInput)
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, registeredCar)
}

// TestCarUsecase_GetByID_Success tests successful retrieval of a car by ID
func TestCarUsecase_GetByID_Success(t *testing.T) {
	t.Parallel()

	// Setup
	ctrl, mockCarRepo, carUsecase := setupTest(t)
	defer ctrl.Finish()

	// Test data
	ctx := context.Background()
	carID := "car-123"
	getInput := input.GetCarByIDInput{
		ID: carID,
	}

	// Create expected car using the factory method (similar to car_test.go)
	now := time.Now()
	expectedCar := model.NewCar("tenant-123", "Toyota Prius", now)

	// Set up expectations for retrieving the car
	mockCarRepo.EXPECT().GetByID(ctx, carID).Return(expectedCar, nil)

	// Execute - Retrieve the car using the usecase
	retrievedCar, err := carUsecase.GetByID(ctx, getInput)
	assert.NoError(t, err)
	assert.NotNil(t, retrievedCar)

	// Verify the returned car
	assert.Equal(t, expectedCar.ID, retrievedCar.ID)
	assert.Equal(t, expectedCar.TenantID, retrievedCar.TenantID)
	assert.Equal(t, expectedCar.Model, retrievedCar.Model)
}

// TestCarUsecase_GetByID_NotFound tests retrieval when car doesn't exist
func TestCarUsecase_GetByID_NotFound(t *testing.T) {
	t.Parallel()

	// Setup
	ctrl, mockCarRepo, carUsecase := setupTest(t)
	defer ctrl.Finish()

	// Test data
	ctx := context.Background()
	carID := "non-existent-car"
	getInput := input.GetCarByIDInput{
		ID: carID,
	}

	// Set up expectations for not found error
	expectedError := assert.AnError
	mockCarRepo.EXPECT().GetByID(ctx, carID).Return(nil, expectedError)

	// Execute - Try to retrieve a non-existent car
	retrievedCar, err := carUsecase.GetByID(ctx, getInput)
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, retrievedCar)
}

// TestCarUsecase_GetByIDWithTenant_Success tests successful retrieval of a car with tenant by ID
func TestCarUsecase_GetByIDWithTenant_Success(t *testing.T) {
	t.Parallel()

	// Setup
	ctrl, mockCarRepo, carUsecase := setupTest(t)
	defer ctrl.Finish()

	// Test data
	ctx := context.Background()
	carID := "car-123"
	getInput := input.GetCarByIDInput{
		ID: carID,
	}

	// Create expected car with tenant using the factory method
	now := time.Now()
	expectedCar := model.NewCar("tenant-123", "Toyota Prius", now)
	expectedTenant := model.NewTenant("tenant-code", now)
	expectedCar.Refs = &model.CarRefs{
		Tenant: expectedTenant,
	}

	// Set up expectations for retrieving the car with tenant
	mockCarRepo.EXPECT().GetByIDWithTenant(ctx, carID).Return(expectedCar, nil)

	// Execute - Retrieve the car with tenant using the usecase
	retrievedCar, err := carUsecase.GetByIDWithTenant(ctx, getInput)
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

// TestCarUsecase_GetByIDWithTenant_NotFound tests retrieval when car with tenant doesn't exist
func TestCarUsecase_GetByIDWithTenant_NotFound(t *testing.T) {
	t.Parallel()

	// Setup
	ctrl, mockCarRepo, carUsecase := setupTest(t)
	defer ctrl.Finish()

	// Test data
	ctx := context.Background()
	carID := "non-existent-car"
	getInput := input.GetCarByIDInput{
		ID: carID,
	}

	// Set up expectations for not found error
	expectedError := assert.AnError
	mockCarRepo.EXPECT().GetByIDWithTenant(ctx, carID).Return(nil, expectedError)

	// Execute - Try to retrieve a non-existent car with tenant
	retrievedCar, err := carUsecase.GetByIDWithTenant(ctx, getInput)
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, retrievedCar)
}

// TestCarUsecase_Register_Validation tests validation failures for Register
func TestCarUsecase_Register_Validation(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input   input.RegisterCarInput
		wantErr string
	}{
		"empty tenant ID": {
			input: input.RegisterCarInput{
				TenantID: "", // Missing required field
				Model:    "Toyota Prius",
			},
			wantErr: "validation failed",
		},
		"empty model": {
			input: input.RegisterCarInput{
				TenantID: "tenant-123",
				Model:    "", // Missing required field
			},
			wantErr: "validation failed",
		},
		"both fields empty": {
			input: input.RegisterCarInput{
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
			ctrl, _, carUsecase := setupTest(t)
			defer ctrl.Finish()

			// Test data
			ctx := context.Background()

			// Execute - Try to register a car with invalid input
			registeredCar, err := carUsecase.Register(ctx, tt.input)

			// Assert
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
			assert.Nil(t, registeredCar)
		})
	}
}

// TestCarUsecase_GetByID_Validation tests validation failures for GetByID
func TestCarUsecase_GetByID_Validation(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input   input.GetCarByIDInput
		wantErr string
	}{
		"empty ID": {
			input: input.GetCarByIDInput{
				ID: "", // Missing required field
			},
			wantErr: "validation failed",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Setup
			ctrl, _, carUsecase := setupTest(t)
			defer ctrl.Finish()

			// Test data
			ctx := context.Background()

			// Execute - Try to retrieve a car with invalid input
			retrievedCar, err := carUsecase.GetByID(ctx, tt.input)

			// Assert
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
			assert.Nil(t, retrievedCar)
		})
	}
}
