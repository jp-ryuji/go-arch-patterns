package repository

import (
	"context"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/jp-ryuji/go-arch-patterns/internal/domain/repository"
	"github.com/jp-ryuji/go-arch-patterns/internal/infrastructure/postgres/entgen"
	"github.com/jp-ryuji/go-arch-patterns/internal/infrastructure/postgres/entgen/outbox"
)

type outboxRepository struct {
	client *entgen.Client
}

// NewOutboxRepository creates a new outbox repository
func NewOutboxRepository(client *entgen.Client) repository.OutboxRepository {
	return &outboxRepository{
		client: client,
	}
}

// Create inserts a new outbox message
func (r *outboxRepository) Create(ctx context.Context, outbox *entgen.Outbox) error {
	_, err := r.client.Outbox.Create().
		SetID(outbox.ID).
		SetAggregateType(outbox.AggregateType).
		SetAggregateID(outbox.AggregateID).
		SetEventType(outbox.EventType).
		SetPayload(outbox.Payload).
		SetCreatedAt(outbox.CreatedAt).
		SetStatus(outbox.Status).
		SetVersion(outbox.Version).
		Save(ctx)
	return err
}

// CreateInTx inserts a new outbox message within a transaction
func (r *outboxRepository) CreateInTx(ctx context.Context, tx *entgen.Tx, outbox *entgen.Outbox) error {
	_, err := tx.Outbox.Create().
		SetID(outbox.ID).
		SetAggregateType(outbox.AggregateType).
		SetAggregateID(outbox.AggregateID).
		SetEventType(outbox.EventType).
		SetPayload(outbox.Payload).
		SetCreatedAt(outbox.CreatedAt).
		SetStatus(outbox.Status).
		SetVersion(outbox.Version).
		Save(ctx)
	return err
}

// GetPending retrieves pending outbox messages up to the specified limit
func (r *outboxRepository) GetPending(ctx context.Context, limit int) ([]*entgen.Outbox, error) {
	return r.client.Outbox.Query().
		Where(outbox.Status("pending")).
		Limit(limit).
		Order(entgen.Asc(outbox.FieldCreatedAt)). // Process in FIFO order
		All(ctx)
}

// GetPendingWithLock retrieves pending outbox messages with locking for concurrency control
// Uses SELECT ... FOR UPDATE SKIP LOCKED for efficient concurrency control
func (r *outboxRepository) GetPendingWithLock(ctx context.Context, limit int, processorID string) ([]*entgen.Outbox, error) {
	// Use SELECT ... FOR UPDATE SKIP LOCKED for efficient concurrency control
	// This ensures that only one processor can claim a message at a time
	pendingMessages, err := r.client.Outbox.Query().
		Where(outbox.Status("pending")).
		Limit(limit).
		Order(entgen.Asc(outbox.FieldCreatedAt)).
		ForUpdate(
			sql.WithLockAction(sql.SkipLocked),
		).
		All(ctx)
	if err != nil {
		return nil, err
	}

	// Update the locked messages with processor information
	now := time.Now()
	lockedMessages := make([]*entgen.Outbox, 0, len(pendingMessages))

	for _, msg := range pendingMessages {
		updatedMsg, err := r.client.Outbox.UpdateOneID(msg.ID).
			SetLockedAt(now).
			SetLockedBy(processorID).
			Save(ctx)
		if err == nil {
			lockedMessages = append(lockedMessages, updatedMsg)
		}
		// If locking fails, another processor got it first - skip it
	}

	return lockedMessages, nil
}

// MarkAsProcessed marks an outbox message as processed
func (r *outboxRepository) MarkAsProcessed(ctx context.Context, id string, processedAt time.Time) error {
	return r.client.Outbox.UpdateOneID(id).
		SetProcessedAt(processedAt).
		SetStatus("processed").
		ClearLockedAt().
		ClearLockedBy().
		Exec(ctx)
}

// MarkAsProcessedInTx marks an outbox message as processed within a transaction
func (r *outboxRepository) MarkAsProcessedInTx(ctx context.Context, tx *entgen.Tx, id string, processedAt time.Time) error {
	return tx.Outbox.UpdateOneID(id).
		SetProcessedAt(processedAt).
		SetStatus("processed").
		ClearLockedAt().
		ClearLockedBy().
		Exec(ctx)
}

// MarkAsFailed marks an outbox message as failed
func (r *outboxRepository) MarkAsFailed(ctx context.Context, id string, errorMessage string) error {
	return r.client.Outbox.UpdateOneID(id).
		SetStatus("failed").
		SetErrorMessage(errorMessage).
		ClearLockedAt().
		ClearLockedBy().
		Exec(ctx)
}

// MarkAsFailedInTx marks an outbox message as failed within a transaction
func (r *outboxRepository) MarkAsFailedInTx(ctx context.Context, tx *entgen.Tx, id string, errorMessage string) error {
	return tx.Outbox.UpdateOneID(id).
		SetStatus("failed").
		SetErrorMessage(errorMessage).
		ClearLockedAt().
		ClearLockedBy().
		Exec(ctx)
}

// GetFailed retrieves failed outbox messages up to the specified limit
func (r *outboxRepository) GetFailed(ctx context.Context, limit int) ([]*entgen.Outbox, error) {
	return r.client.Outbox.Query().
		Where(outbox.Status("failed")).
		Limit(limit).
		All(ctx)
}

// UnlockOrphanedMessages unlocks messages that have been locked for too long
func (r *outboxRepository) UnlockOrphanedMessages(ctx context.Context, olderThan time.Duration) (int, error) {
	cutoffTime := time.Now().Add(-olderThan)

	affected, err := r.client.Outbox.Update().
		Where(
			outbox.LockedAtNotNil(),
			outbox.LockedAtLT(cutoffTime),
		).
		ClearLockedAt().
		ClearLockedBy().
		Save(ctx)

	return affected, err
}

// CleanupProcessedMessages removes processed messages older than the specified duration
func (r *outboxRepository) CleanupProcessedMessages(ctx context.Context, olderThan time.Duration) (int, error) {
	cutoffTime := time.Now().Add(-olderThan)

	affected, err := r.client.Outbox.Delete().
		Where(
			outbox.Status("processed"),
			outbox.ProcessedAtNotNil(),
			outbox.ProcessedAtLT(cutoffTime),
		).
		Exec(ctx)

	return affected, err
}
