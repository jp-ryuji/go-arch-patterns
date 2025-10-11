package rentalservice_test

import (
	"context"
	"testing"
	"time"

	sharedrepository "github.com/jp-ryuji/go-arch-patterns/internal/domain/repository"
	sharedrepositorymock "github.com/jp-ryuji/go-arch-patterns/internal/domain/repository/mock"
	rental "github.com/jp-ryuji/go-arch-patterns/internal/rental"
	rentalinput "github.com/jp-ryuji/go-arch-patterns/internal/rental/application/input"
	rentalservice "github.com/jp-ryuji/go-arch-patterns/internal/rental/application/service"
	rentalrepositorymock "github.com/jp-ryuji/go-arch-patterns/internal/rental/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

type fakeTransaction struct {
	raw any
}

func (f *fakeTransaction) Raw() any {
	return f.raw
}

var _ sharedrepository.Transaction = (*fakeTransaction)(nil)

func setupTest(t *testing.T) (*gomock.Controller, *rentalrepositorymock.MockCarRepository, *sharedrepositorymock.MockOutboxRepository, *sharedrepositorymock.MockTransactionManager, rentalservice.CarService) {
	t.Helper()

	ctrl := gomock.NewController(t)
	carRepo := rentalrepositorymock.NewMockCarRepository(ctrl)
	outboxRepo := sharedrepositorymock.NewMockOutboxRepository(ctrl)
	txManager := sharedrepositorymock.NewMockTransactionManager(ctrl)

	service := rentalservice.NewCarService(carRepo, outboxRepo, txManager)

	return ctrl, carRepo, outboxRepo, txManager, service
}

func TestCarService_Create_Success(t *testing.T) {
	t.Parallel()

	ctrl, carRepo, outboxRepo, txManager, service := setupTest(t)
	t.Cleanup(ctrl.Finish)

	ctx := context.Background()
	registerInput := rentalinput.CreateCar{
		TenantID: "tenant-123",
		Model:    "Toyota Prius",
	}

	mockTx := &fakeTransaction{raw: struct{}{}}
	txManager.EXPECT().BeginTx(ctx).Return(mockTx, nil)
	txManager.EXPECT().CommitTx(ctx, mockTx).Return(nil)

	var createdCar *rental.Car
	carRepo.EXPECT().CreateInTx(ctx, mockTx, gomock.Any()).DoAndReturn(func(_ context.Context, _ sharedrepository.Transaction, car *rental.Car) error {
		createdCar = car
		assert.Equal(t, registerInput.TenantID, car.TenantID)
		assert.Equal(t, registerInput.Model, car.Model)
		assert.NotEmpty(t, car.ID)
		assert.WithinDuration(t, time.Now(), car.CreatedAt, time.Second)
		assert.WithinDuration(t, time.Now(), car.UpdatedAt, time.Second)
		return nil
	})

	outboxRepo.EXPECT().CreateInTx(ctx, mockTx, gomock.Any()).Return(nil)

	result, err := service.Create(ctx, registerInput)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	assert.Equal(t, registerInput.TenantID, result.TenantID)
	assert.Equal(t, registerInput.Model, result.Model)
	assert.Equal(t, createdCar.ID, result.ID)
	assert.Equal(t, createdCar.TenantID, result.TenantID)
}

func TestCarService_Create_RepositoryError(t *testing.T) {
	t.Parallel()

	ctrl, carRepo, _, txManager, service := setupTest(t)
	t.Cleanup(ctrl.Finish)

	ctx := context.Background()
	registerInput := rentalinput.CreateCar{
		TenantID: "tenant-123",
		Model:    "Toyota Prius",
	}

	mockTx := &fakeTransaction{raw: struct{}{}}
	txManager.EXPECT().BeginTx(ctx).Return(mockTx, nil)
	txManager.EXPECT().RollbackTx(ctx, mockTx).Return(nil)

	expectedErr := assert.AnError
	carRepo.EXPECT().CreateInTx(ctx, mockTx, gomock.Any()).Return(expectedErr)

	result, err := service.Create(ctx, registerInput)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), expectedErr.Error())
	assert.Nil(t, result)
}

func TestCarService_GetByID_Success(t *testing.T) {
	t.Parallel()

	ctrl, carRepo, _, _, service := setupTest(t)
	t.Cleanup(ctrl.Finish)

	ctx := context.Background()
	getInput := rentalinput.GetCarByID{ID: "car-123"}
	expectedCar := rental.NewCar("tenant-123", "Toyota Prius", time.Now())

	carRepo.EXPECT().GetByID(ctx, getInput.ID).Return(expectedCar, nil)

	result, err := service.GetByID(ctx, getInput)
	assert.NoError(t, err)
	assert.Equal(t, expectedCar, result)
}

func TestCarService_GetByID_NotFound(t *testing.T) {
	t.Parallel()

	ctrl, carRepo, _, _, service := setupTest(t)
	t.Cleanup(ctrl.Finish)

	ctx := context.Background()
	getInput := rentalinput.GetCarByID{ID: "missing"}
	expectedErr := assert.AnError

	carRepo.EXPECT().GetByID(ctx, getInput.ID).Return(nil, expectedErr)

	result, err := service.GetByID(ctx, getInput)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, result)
}

func TestCarService_GetByIDWithTenant_Success(t *testing.T) {
	t.Parallel()

	ctrl, carRepo, _, _, service := setupTest(t)
	t.Cleanup(ctrl.Finish)

	ctx := context.Background()
	getInput := rentalinput.GetCarByID{ID: "car-123"}
	expectedCar := rental.NewCar("tenant-123", "Toyota Prius", time.Now())

	carRepo.EXPECT().GetByIDWithTenant(ctx, getInput.ID).Return(expectedCar, nil)

	result, err := service.GetByIDWithTenant(ctx, getInput)
	assert.NoError(t, err)
	assert.Equal(t, expectedCar, result)
}

func TestCarService_List_Success(t *testing.T) {
	t.Parallel()

	ctrl, carRepo, _, _, service := setupTest(t)
	t.Cleanup(ctrl.Finish)

	ctx := context.Background()
	listInput := rentalinput.ListCars{TenantID: "tenant-123"}

	expectedCars := []*rental.Car{
		rental.NewCar("tenant-123", "Car A", time.Now()),
		rental.NewCar("tenant-123", "Car B", time.Now()),
	}

	carRepo.EXPECT().ListByTenant(ctx, listInput.TenantID, gomock.Any(), gomock.Any()).Return(expectedCars, "", int32(len(expectedCars)), nil)

	result, err := service.List(ctx, listInput)
	assert.NoError(t, err)
	assert.Len(t, result.Cars, len(expectedCars))
	assert.Equal(t, expectedCars[0].Model, result.Cars[0].Model)
}
