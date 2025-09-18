package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/jp-ryuji/go-arch-patterns/internal/domain/model"
	"github.com/jp-ryuji/go-arch-patterns/internal/domain/repository"
	"github.com/jp-ryuji/go-arch-patterns/internal/infrastructure/postgres/entgen"
	"github.com/jp-ryuji/go-arch-patterns/internal/pkg/id"
	"github.com/jp-ryuji/go-arch-patterns/internal/usecase/input"
)

// carUsecase implements CarUsecase interface
type carUsecase struct {
	carRepo    repository.CarRepository
	outboxRepo repository.OutboxRepository
	txManager  repository.TransactionManager
}

// NewCarUsecase creates a new car usecase
func NewCarUsecase(
	carRepo repository.CarRepository,
	outboxRepo repository.OutboxRepository,
	txManager repository.TransactionManager,
) CarUsecase {
	return &carUsecase{
		carRepo:    carRepo,
		outboxRepo: outboxRepo,
		txManager:  txManager,
	}
}

// Register registers a new car using the outbox pattern with transactional guarantees
func (uc *carUsecase) Register(ctx context.Context, input input.RegisterCarInput) (*model.Car, error) {
	// Validate input
	if err := Validate(input); err != nil {
		return nil, err
	}

	now := time.Now()
	car := model.NewCar(input.TenantID, input.Model, now)

	// Start a transaction
	tx, err := uc.txManager.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}

	// Ensure rollback in case of error
	defer func() {
		if r := recover(); r != nil {
			if err := uc.txManager.RollbackTx(ctx, tx); err != nil {
				// Log the error but continue with the panic
				// In a production system, you might want to use a proper logger
				fmt.Printf("Failed to rollback transaction: %v\n", err)
			}
			panic(r) // re-panic
		}
	}()

	// Step 1: Save to PostgreSQL within transaction
	if err := uc.carRepo.CreateInTx(ctx, tx, car); err != nil {
		if rollbackErr := uc.txManager.RollbackTx(ctx, tx); rollbackErr != nil {
			return nil, fmt.Errorf("failed to create car in database: %w; also failed to rollback transaction: %v", err, rollbackErr)
		}
		return nil, fmt.Errorf("failed to create car in database: %w", err)
	}

	// Step 2: Create outbox message for OpenSearch within transaction
	outbox := &entgen.Outbox{
		ID:            id.New(),
		AggregateType: "car",
		AggregateID:   car.ID,
		EventType:     "car_created",
		Payload: map[string]interface{}{
			"id":         car.ID,
			"tenant_id":  car.TenantID,
			"model":      car.Model,
			"created_at": car.CreatedAt,
			"updated_at": car.UpdatedAt,
		},
		CreatedAt: now,
		Status:    "pending",
	}

	if err := uc.outboxRepo.CreateInTx(ctx, tx, outbox); err != nil {
		if rollbackErr := uc.txManager.RollbackTx(ctx, tx); rollbackErr != nil {
			return nil, fmt.Errorf("failed to create outbox message: %w; also failed to rollback transaction: %v", err, rollbackErr)
		}
		return nil, fmt.Errorf("failed to create outbox message: %w", err)
	}

	// Commit the transaction
	if err := uc.txManager.CommitTx(ctx, tx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return car, nil
}

// GetByID retrieves a car by its ID
func (uc *carUsecase) GetByID(ctx context.Context, input input.GetCarByIDInput) (*model.Car, error) {
	// Validate input
	if err := Validate(input); err != nil {
		return nil, err
	}

	return uc.carRepo.GetByID(ctx, input.ID)
}

// GetByIDWithTenant retrieves a car by its ID along with its tenant information
func (uc *carUsecase) GetByIDWithTenant(ctx context.Context, input input.GetCarByIDInput) (*model.Car, error) {
	// Validate input
	if err := Validate(input); err != nil {
		return nil, err
	}

	return uc.carRepo.GetByIDWithTenant(ctx, input.ID)
}
