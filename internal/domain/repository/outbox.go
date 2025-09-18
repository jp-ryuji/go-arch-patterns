package repository

import (
	"context"
	"time"

	"github.com/jp-ryuji/go-arch-patterns/internal/infrastructure/postgres/entgen"
)

//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock_repository
type OutboxRepository interface {
	Create(ctx context.Context, outbox *entgen.Outbox) error
	CreateInTx(ctx context.Context, tx *entgen.Tx, outbox *entgen.Outbox) error
	GetPending(ctx context.Context, limit int) ([]*entgen.Outbox, error)
	GetPendingWithLock(ctx context.Context, limit int, processorID string) ([]*entgen.Outbox, error)
	MarkAsProcessed(ctx context.Context, id string, processedAt time.Time) error
	MarkAsProcessedInTx(ctx context.Context, tx *entgen.Tx, id string, processedAt time.Time) error
	MarkAsFailed(ctx context.Context, id string, errorMessage string) error
	MarkAsFailedInTx(ctx context.Context, tx *entgen.Tx, id string, errorMessage string) error
	GetFailed(ctx context.Context, limit int) ([]*entgen.Outbox, error)
	UnlockOrphanedMessages(ctx context.Context, olderThan time.Duration) (int, error)
	CleanupProcessedMessages(ctx context.Context, olderThan time.Duration) (int, error)
}
