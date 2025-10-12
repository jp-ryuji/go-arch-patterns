package repository

import (
	"context"
	"time"

	sharedrepository "github.com/jp-ryuji/go-arch-patterns/internal/domain/repository"
	"github.com/jp-ryuji/go-arch-patterns/internal/infrastructure/postgres/entgen"
	"github.com/jp-ryuji/go-arch-patterns/internal/infrastructure/postgres/entgen/outbox"
)

type outboxRepository struct {
	client *entgen.Client
}

// NewOutboxRepository creates a new outbox repository
func NewOutboxRepository(client *entgen.Client) sharedrepository.OutboxRepository {
	return &outboxRepository{
		client: client,
	}
}

// Create inserts a new outbox message
func (r *outboxRepository) Create(ctx context.Context, message *sharedrepository.OutboxMessage) error {
	_, err := r.client.Outbox.Create().
		SetID(message.ID).
		SetAggregateType(message.AggregateType).
		SetAggregateID(message.AggregateID).
		SetEventType(message.EventType).
		SetPayload(message.Payload).
		SetCreatedAt(message.CreatedAt).
		SetStatus(string(message.Status)).
		SetVersion(message.Version).
		Save(ctx)
	return err
}

// CreateInTx inserts a new outbox message within a transaction
func (r *outboxRepository) CreateInTx(ctx context.Context, tx sharedrepository.Transaction, message *sharedrepository.OutboxMessage) error {
	entTx, err := asEntTx(tx)
	if err != nil {
		return err
	}

	_, err = entTx.Outbox.Create().
		SetID(message.ID).
		SetAggregateType(message.AggregateType).
		SetAggregateID(message.AggregateID).
		SetEventType(message.EventType).
		SetPayload(message.Payload).
		SetCreatedAt(message.CreatedAt).
		SetStatus(string(message.Status)).
		SetVersion(message.Version).
		Save(ctx)
	return err
}

// GetPending retrieves pending outbox messages up to the specified limit
func (r *outboxRepository) GetPending(ctx context.Context, limit int) ([]*sharedrepository.OutboxMessage, error) {
	result, err := r.client.Outbox.Query().
		Where(outbox.Status("pending")).
		Limit(limit).
		Order(entgen.Asc(outbox.FieldCreatedAt)). // Process in FIFO order
		All(ctx)
	if err != nil {
		return nil, err
	}

	return toDomainMessages(result), nil
}

// GetPendingWithLock retrieves pending outbox messages with locking
func (r *outboxRepository) GetPendingWithLock(ctx context.Context, limit int, processorID string) ([]*sharedrepository.OutboxMessage, error) {
	pendingMessages, err := r.client.Outbox.Query().
		Where(outbox.Status("pending")).
		Limit(limit).
		Order(entgen.Asc(outbox.FieldCreatedAt)).
		ForUpdate().
		All(ctx)
	if err != nil {
		return nil, err
	}

	// Update the locked messages with processor information
	now := time.Now()
	lockedMessages := make([]*sharedrepository.OutboxMessage, 0, len(pendingMessages))

	for _, msg := range pendingMessages {
		updatedMsg, err := r.client.Outbox.UpdateOneID(msg.ID).
			SetLockedAt(now).
			SetLockedBy(processorID).
			Save(ctx)
		if err == nil {
			lockedMessages = append(lockedMessages, toDomainMessage(updatedMsg))
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
func (r *outboxRepository) MarkAsProcessedInTx(ctx context.Context, tx sharedrepository.Transaction, id string, processedAt time.Time) error {
	entTx, err := asEntTx(tx)
	if err != nil {
		return err
	}

	return entTx.Outbox.UpdateOneID(id).
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
func (r *outboxRepository) MarkAsFailedInTx(ctx context.Context, tx sharedrepository.Transaction, id string, errorMessage string) error {
	entTx, err := asEntTx(tx)
	if err != nil {
		return err
	}

	return entTx.Outbox.UpdateOneID(id).
		SetStatus("failed").
		SetErrorMessage(errorMessage).
		ClearLockedAt().
		ClearLockedBy().
		Exec(ctx)
}

// GetFailed retrieves failed outbox messages up to the specified limit
func (r *outboxRepository) GetFailed(ctx context.Context, limit int) ([]*sharedrepository.OutboxMessage, error) {
	result, err := r.client.Outbox.Query().
		Where(outbox.Status("failed")).
		Limit(limit).
		All(ctx)
	if err != nil {
		return nil, err
	}

	return toDomainMessages(result), nil
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

func toDomainMessages(outboxes []*entgen.Outbox) []*sharedrepository.OutboxMessage {
	msgs := make([]*sharedrepository.OutboxMessage, 0, len(outboxes))
	for _, o := range outboxes {
		msgs = append(msgs, toDomainMessage(o))
	}
	return msgs
}

func toDomainMessage(outbox *entgen.Outbox) *sharedrepository.OutboxMessage {
	if outbox == nil {
		return nil
	}

	var lockedAt *time.Time
	if outbox.LockedAt != nil {
		lockedAt = outbox.LockedAt
	}

	var lockedBy *string
	if outbox.LockedBy != nil {
		lockedBy = outbox.LockedBy
	}

	var processedAt *time.Time
	if outbox.ProcessedAt != nil {
		processedAt = outbox.ProcessedAt
	}

	var errorMessage *string
	if outbox.ErrorMessage != nil {
		errorMessage = outbox.ErrorMessage
	}

	return &sharedrepository.OutboxMessage{
		ID:            outbox.ID,
		AggregateType: outbox.AggregateType,
		AggregateID:   outbox.AggregateID,
		EventType:     outbox.EventType,
		Payload:       outbox.Payload,
		CreatedAt:     outbox.CreatedAt,
		Status:        sharedrepository.OutboxStatus(outbox.Status),
		Version:       outbox.Version,
		LockedAt:      lockedAt,
		LockedBy:      lockedBy,
		ProcessedAt:   processedAt,
		ErrorMessage:  errorMessage,
	}
}
