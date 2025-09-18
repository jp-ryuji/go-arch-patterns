package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/jp-ryuji/go-arch-patterns/internal/domain/model"
	"github.com/jp-ryuji/go-arch-patterns/internal/domain/repository"
	"github.com/jp-ryuji/go-arch-patterns/internal/infrastructure/opensearch/client"
	"github.com/jp-ryuji/go-arch-patterns/internal/infrastructure/opensearch/processor/metrics"
	opensearchrepo "github.com/jp-ryuji/go-arch-patterns/internal/infrastructure/opensearch/repository"
	"github.com/jp-ryuji/go-arch-patterns/internal/infrastructure/postgres/entgen"
	"github.com/oklog/ulid/v2"
	opensearch "github.com/opensearch-project/opensearch-go/v4"
)

// OutboxProcessor processes outbox messages and sends them to OpenSearch
type OutboxProcessor struct {
	outboxRepo      repository.OutboxRepository
	opensearchRepo  *opensearchrepo.CarRepository
	batchSize       int
	maxRetries      int
	retryDelay      time.Duration
	processorID     string
	metrics         *metrics.OutboxMetrics
	messageSequence map[string]int64 // For deduplication
	sequenceMutex   sync.RWMutex
}

// NewOutboxProcessor creates a new outbox processor
func NewOutboxProcessor(
	outboxRepo repository.OutboxRepository,
	batchSize int,
	maxRetries int,
	retryDelay time.Duration,
) (*OutboxProcessor, error) {
	// Create config directly like in the client
	cfg := opensearch.Config{
		Addresses: []string{
			fmt.Sprintf("http://localhost:%s", "9200"),
		},
		Username: "admin",
		Password: "EX2R3L(,M-tR",
	}

	// Create OpenSearch repository
	opensearchRepo := opensearchrepo.NewCarRepository(
		cfg,
		client.GetDefaultIndex(),
	)

	return &OutboxProcessor{
		outboxRepo:      outboxRepo,
		opensearchRepo:  opensearchRepo,
		batchSize:       batchSize,
		maxRetries:      maxRetries,
		retryDelay:      retryDelay,
		processorID:     ulid.Make().String(),
		metrics:         metrics.NewOutboxMetrics(),
		messageSequence: make(map[string]int64),
	}, nil
}

// ProcessMessage processes a single outbox message with retry mechanisms, metrics, and deduplication
func (p *OutboxProcessor) ProcessMessage(ctx context.Context, message *entgen.Outbox) error {
	startTime := time.Now()

	// Check for deduplication
	if p.isDuplicateMessage(message) {
		log.Printf("Skipping duplicate message %s", message.ID)
		// Still mark as processed to avoid reprocessing
		if err := p.outboxRepo.MarkAsProcessed(ctx, message.ID, time.Now()); err != nil {
			log.Printf("Failed to mark duplicate message %s as processed: %v", message.ID, err)
			return err
		}
		return nil
	}

	processStartTime := time.Now()
	retries, err := p.processMessageWithRetry(ctx, message)
	processDuration := time.Since(processStartTime)

	if err == nil {
		p.metrics.RecordMessageProcessed(processDuration, retries)
		// Record message sequence for deduplication
		p.recordMessageSequence(message)
		log.Printf("Successfully processed message %s in %v", message.ID, processDuration)
	} else {
		p.metrics.RecordMessageFailed(p.getErrorType(err))
		log.Printf("Failed to process message %s after retries: %v", message.ID, err)
		return err
	}

	log.Printf("Processed message in %v. Metrics: %+v", time.Since(startTime), p.metrics.GetMetrics())
	return nil
}

// Process processes a batch of outbox messages (for backward compatibility)
func (p *OutboxProcessor) Process(ctx context.Context) error {
	startTime := time.Now()

	// Unlock orphaned messages first
	if _, err := p.outboxRepo.UnlockOrphanedMessages(ctx, 5*time.Minute); err != nil {
		log.Printf("Failed to unlock orphaned messages: %v", err)
	}

	// Get pending messages with locking for concurrency control
	messages, err := p.outboxRepo.GetPendingWithLock(ctx, p.batchSize, p.processorID)
	if err != nil {
		p.metrics.RecordMessageFailed("database")
		return fmt.Errorf("failed to get pending messages: %w", err)
	}

	// Update queue metrics
	failedCount, _ := p.getFailedMessageCount(ctx) // Ignore error for metrics
	p.metrics.UpdateQueueMetrics(int64(len(messages)), failedCount)
	p.metrics.RecordBatchProcessed(len(messages))

	for _, message := range messages {
		// Process each message individually
		if err := p.ProcessMessage(ctx, message); err != nil {
			log.Printf("Failed to process message %s: %v", message.ID, err)
			// Continue processing other messages
		}
	}

	// Cleanup old processed messages
	if _, err := p.outboxRepo.CleanupProcessedMessages(ctx, 24*time.Hour); err != nil {
		log.Printf("Failed to cleanup processed messages: %v", err)
	}

	log.Printf("Processed batch in %v. Metrics: %+v", time.Since(startTime), p.metrics.GetMetrics())
	return nil
}

// processMessageWithRetry processes a single outbox message with retry mechanisms
func (p *OutboxProcessor) processMessageWithRetry(ctx context.Context, message *entgen.Outbox) (int, error) {
	var lastErr error

	for attempt := 0; attempt <= p.maxRetries; attempt++ {
		if attempt > 0 {
			// Wait before retrying
			time.Sleep(p.retryDelay * time.Duration(attempt))
		}

		err := p.processMessage(ctx, message)
		if err == nil {
			// Success - mark message as processed
			if err := p.outboxRepo.MarkAsProcessed(ctx, message.ID, time.Now()); err != nil {
				log.Printf("Failed to mark message %s as processed: %v", message.ID, err)
			}
			return attempt, nil
		}

		lastErr = err
		log.Printf("Attempt %d failed for message %s: %v", attempt+1, message.ID, err)
	}

	// All retries failed - mark message as failed
	errorMessage := fmt.Sprintf("Failed after %d attempts: %v", p.maxRetries+1, lastErr)
	if err := p.outboxRepo.MarkAsFailed(ctx, message.ID, errorMessage); err != nil {
		log.Printf("Failed to mark message %s as failed: %v", message.ID, err)
	}

	return p.maxRetries, fmt.Errorf("message processing failed after %d attempts: %w", p.maxRetries+1, lastErr)
}

// processMessage processes a single outbox message
func (p *OutboxProcessor) processMessage(ctx context.Context, message *entgen.Outbox) error {
	switch message.EventType {
	case "car_created":
		return p.processCarCreated(ctx, message)
	case "car_updated":
		return p.processCarUpdated(ctx, message)
	case "car_deleted":
		return p.processCarDeleted(ctx, message)
	default:
		return fmt.Errorf("unknown event type: %s", message.EventType)
	}
}

// processCarCreated processes a car created event
func (p *OutboxProcessor) processCarCreated(ctx context.Context, message *entgen.Outbox) error {
	// Convert payload to car model
	car, err := p.payloadToCar(message.Payload)
	if err != nil {
		return fmt.Errorf("failed to convert payload to car: %w", err)
	}

	// Send to OpenSearch
	if err := p.opensearchRepo.Create(ctx, car); err != nil {
		return fmt.Errorf("failed to create car in OpenSearch: %w", err)
	}

	return nil
}

// processCarUpdated processes a car updated event
func (p *OutboxProcessor) processCarUpdated(ctx context.Context, message *entgen.Outbox) error {
	// Convert payload to car model
	car, err := p.payloadToCar(message.Payload)
	if err != nil {
		return fmt.Errorf("failed to convert payload to car: %w", err)
	}

	// Send to OpenSearch
	if err := p.opensearchRepo.Update(ctx, car); err != nil {
		return fmt.Errorf("failed to update car in OpenSearch: %w", err)
	}

	return nil
}

// processCarDeleted processes a car deleted event
func (p *OutboxProcessor) processCarDeleted(ctx context.Context, message *entgen.Outbox) error {
	// Get car ID from payload
	carID, ok := message.Payload["id"].(string)
	if !ok {
		return fmt.Errorf("failed to get car ID from payload")
	}

	// Send to OpenSearch
	if err := p.opensearchRepo.Delete(ctx, carID); err != nil {
		return fmt.Errorf("failed to delete car from OpenSearch: %w", err)
	}

	return nil
}

// payloadToCar converts a payload to a car model
func (p *OutboxProcessor) payloadToCar(payload map[string]interface{}) (*model.Car, error) {
	// Marshal payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Unmarshal JSON to car model
	var car model.Car
	if err := json.Unmarshal(jsonData, &car); err != nil {
		return nil, fmt.Errorf("failed to unmarshal to car: %w", err)
	}

	return &car, nil
}

// StartProcessing starts the outbox processor in a queue-based mode
// This method is intended to be used with a queue mechanism instead of polling
func (p *OutboxProcessor) StartProcessing(ctx context.Context) error {
	// Unlock orphaned messages first
	if _, err := p.outboxRepo.UnlockOrphanedMessages(ctx, 5*time.Minute); err != nil {
		log.Printf("Failed to unlock orphaned messages: %v", err)
		return err
	}

	// Note: In a queue-based implementation, messages would be received from a queue
	// and processed using ProcessMessage method
	// This is a placeholder for queue integration

	log.Printf("Outbox processor ready for queue-based processing")
	return nil
}

// GetMetrics returns current metrics
func (p *OutboxProcessor) GetMetrics() map[string]interface{} {
	return p.metrics.GetMetrics()
}

// ResetMetrics resets all metrics
func (p *OutboxProcessor) ResetMetrics() {
	p.metrics.Reset()
}

// isDuplicateMessage checks if a message is a duplicate
func (p *OutboxProcessor) isDuplicateMessage(message *entgen.Outbox) bool {
	p.sequenceMutex.RLock()
	defer p.sequenceMutex.RUnlock()

	// Check if we've already processed this message
	if seq, exists := p.messageSequence[message.ID]; exists {
		// Check if this is a retry of the same message
		currentSeq := message.Version
		if currentSeq <= seq {
			return true
		}
	}

	return false
}

// recordMessageSequence records a message sequence for deduplication
func (p *OutboxProcessor) recordMessageSequence(message *entgen.Outbox) {
	p.sequenceMutex.Lock()
	defer p.sequenceMutex.Unlock()

	p.messageSequence[message.ID] = message.Version
}

// getErrorType categorizes errors for metrics
func (p *OutboxProcessor) getErrorType(err error) string {
	errStr := err.Error()
	switch {
	case contains(errStr, "database", "sql", "connection"):
		return "database"
	case contains(errStr, "opensearch", "index", "search"):
		return "opensearch"
	case contains(errStr, "marshal", "unmarshal", "json"):
		return "serialization"
	default:
		return "unknown"
	}
}

// getFailedMessageCount gets the count of failed messages for metrics
func (p *OutboxProcessor) getFailedMessageCount(ctx context.Context) (int64, error) {
	failedMessages, err := p.outboxRepo.GetFailed(ctx, 1000) // Limit to avoid performance issues
	if err != nil {
		return 0, err
	}
	return int64(len(failedMessages)), nil
}

// contains checks if a string contains any of the given substrings
func contains(s string, substrs ...string) bool {
	for _, substr := range substrs {
		if containsIgnoreCase(s, substr) {
			return true
		}
	}
	return false
}

// containsIgnoreCase checks if a string contains a substring (case insensitive)
func containsIgnoreCase(s, substr string) bool {
	// Simple case-insensitive contains
	sLower := ""
	substrLower := ""
	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			sLower += string(r + 32)
		} else {
			sLower += string(r)
		}
	}
	for _, r := range substr {
		if r >= 'A' && r <= 'Z' {
			substrLower += string(r + 32)
		} else {
			substrLower += string(r)
		}
	}
	return len(sLower) >= len(substrLower) &&
		(len(sLower) == len(substrLower) && sLower == substrLower ||
			len(sLower) > len(substrLower) && (sLower[:len(substrLower)] == substrLower ||
				sLower[len(sLower)-len(substrLower):] == substrLower ||
				containsSubstring(sLower, substrLower)))
}

// containsSubstring checks if a substring exists within a string
func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
