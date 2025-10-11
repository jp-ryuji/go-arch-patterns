package repository

import (
	"context"
	"time"
)

type OutboxStatus string

const (
	OutboxStatusPending   OutboxStatus = "pending"
	OutboxStatusProcessed OutboxStatus = "processed"
	OutboxStatusFailed    OutboxStatus = "failed"
)

type OutboxMessage struct {
	ID            string
	AggregateType string
	AggregateID   string
	EventType     string
	Payload       map[string]any
	CreatedAt     time.Time
	Status        OutboxStatus
	Version       int64
	LockedAt      *time.Time
	LockedBy      *string
	ProcessedAt   *time.Time
	ErrorMessage  *string
}

//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock_repository
type OutboxRepository interface {
	Create(ctx context.Context, outbox *OutboxMessage) error
	CreateInTx(ctx context.Context, tx Transaction, outbox *OutboxMessage) error
	GetPending(ctx context.Context, limit int) ([]*OutboxMessage, error)
	GetPendingWithLock(ctx context.Context, limit int, processorID string) ([]*OutboxMessage, error)
	MarkAsProcessed(ctx context.Context, id string, processedAt time.Time) error
	MarkAsProcessedInTx(ctx context.Context, tx Transaction, id string, processedAt time.Time) error
	MarkAsFailed(ctx context.Context, id string, errorMessage string) error
	MarkAsFailedInTx(ctx context.Context, tx Transaction, id string, errorMessage string) error
	GetFailed(ctx context.Context, limit int) ([]*OutboxMessage, error)
	UnlockOrphanedMessages(ctx context.Context, olderThan time.Duration) (int, error)
	CleanupProcessedMessages(ctx context.Context, olderThan time.Duration) (int, error)
}
