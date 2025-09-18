# Outbox Pattern Implementation

This document explains the implementation of the Outbox pattern in this project, which ensures data consistency across multiple services (PostgreSQL and OpenSearch).

## Overview

The Outbox pattern is a technique for reliably publishing events from a service to external systems. It addresses the challenge of ensuring that data changes in the primary database are reliably propagated to external systems like search engines usually via message queues.

In our implementation, we coordinate between:

1. PostgreSQL (primary database)
2. OpenSearch (search engine)

## Implementation Details

### Key Files

1. **Domain Layer**:
   - `internal/domain/repository/outbox.go` - Outbox repository interface
   - `internal/domain/repository/transaction.go` - Transaction manager interface

2. **Infrastructure Layer**:
   - `internal/infrastructure/postgres/ent/schema/outbox.go` - Outbox table schema
   - `internal/infrastructure/postgres/repository/outbox_repository.go` - Outbox repository implementation
   - `internal/infrastructure/postgres/repository/transaction_manager.go` - Transaction manager implementation
   - `internal/infrastructure/opensearch/repository/car_repository.go` - OpenSearch repository implementation
   - `internal/infrastructure/opensearch/client/client.go` - OpenSearch client configuration
   - `internal/infrastructure/opensearch/processor/processor.go` - Outbox processor that sends messages to OpenSearch
   - `internal/infrastructure/opensearch/processor/metrics/metrics.go` - Metrics tracking for the processor

3. **Application Layer**:
   - `internal/infrastructure/usecase/car_impl.go` - Car usecase implementation with outbox pattern
   - `internal/infrastructure/usecase/car.go` - Car usecase interface

4. **Tests**:
   - `internal/infrastructure/opensearch/processor/processor_test.go` - Unit tests for outbox processor
   - `internal/infrastructure/usecase/car_impl_test.go` - Unit tests for car usecase with transactional outbox

### Outbox Flow

The outbox pattern is implemented in the car usecase:

1. When a car is registered, a database transaction is started
2. The car data is saved to PostgreSQL within the transaction
3. An outbox message is created for OpenSearch within the same transaction
4. The transaction is committed (ensuring atomicity)
5. The outbox processor receives messages from a queue and sends them to OpenSearch

This approach separates the concerns of saving data and processing messages, with the latter being handled by a different process using a queue mechanism instead of polling.

## Transactional Guarantees

The implementation ensures atomicity between the main entity creation and outbox message creation by using database transactions. Both operations happen within the same transaction, so either both succeed or both fail.

## Message Deduplication and Sequencing

The implementation includes robust deduplication mechanisms:

1. **Message Versioning**: Each message has a version field to track duplicates
2. **In-Memory Sequence Tracking**: The processor maintains an in-memory map of processed message sequences
3. **Idempotent Operations**: OpenSearch operations are idempotent by design

## Concurrency Control

The implementation supports multiple processor instances with:

1. **Message Locking**: Messages are locked when retrieved to prevent race conditions
2. **Orphaned Message Cleanup**: Periodic cleanup of messages locked by crashed processors
3. **FIFO Processing**: Messages are processed in first-in-first-out order

## Metrics and Monitoring

The implementation includes comprehensive metrics tracking:

```go
type OutboxMetrics struct {
    // Message processing metrics
    messagesProcessed atomic.Int64
    messagesFailed    atomic.Int64
    messagesSucceeded atomic.Int64

    // Timing metrics
    totalProcessingTime atomic.Int64
    minProcessingTime   atomic.Int64
    maxProcessingTime   atomic.Int64

    // Batch metrics
    batchesProcessed atomic.Int64
    avgBatchSize     atomic.Int64

    // Retry metrics
    totalRetries atomic.Int64
    maxRetries   atomic.Int64

    // Error metrics
    databaseErrors     atomic.Int64
    opensearchErrors   atomic.Int64
    serializationErrors atomic.Int64

    // Queue metrics
    pendingMessages atomic.Int64
    failedMessages  atomic.Int64
}
```

## Error Handling and Retry Mechanisms

The implementation includes:

1. **Retry Logic**: Failed messages are retried up to a configurable number of times
2. **Exponential Backoff**: Delay between retries increases with each attempt
3. **Dead Letter Queue**: Messages that fail after all retries are marked as "failed" for manual inspection
4. **Graceful Error Handling**: Processing continues even if individual messages fail
5. **Error Categorization**: Different error types are tracked separately for better observability

## Message Retention and Cleanup

The implementation includes automatic cleanup mechanisms:

1. **Processed Message Cleanup**: Processed messages are automatically removed after 24 hours
2. **Failed Message Retention**: Failed messages are retained for manual inspection
3. **Orphaned Message Cleanup**: Messages locked by crashed processors are automatically unlocked

## Schema Evolution

The outbox schema is designed to support schema evolution:

1. **Flexible Payload**: JSON payload field can accommodate different message formats
2. **Version Field**: Message versioning for tracking schema changes
3. **Extensible Fields**: Additional fields can be added without breaking existing code

## Performance Optimizations

The implementation includes several performance optimizations:

1. **Database Indexes**: Proper indexing on status, created_at, and locked fields
2. **Batch Processing**: Messages are processed in batches for efficiency
3. **Connection Pooling**: Efficient use of database and OpenSearch connections
4. **Asynchronous Processing**: Non-blocking message processing

## Outbox Processor

The outbox processor is responsible for sending messages to OpenSearch with robust error handling. The processor has been refactored to support both batch processing and queue-based processing:

```go
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
```

The processor can now be integrated with message queues like Redis, RabbitMQ, or Apache Kafka for more scalable and efficient processing.

## Testing

The implementation includes comprehensive tests:

### Unit Tests

Unit tests verify:

- Successful transactional outbox pattern execution
- Proper error handling when database operations fail
- Input validation
- Transaction rollback on errors
- Metrics tracking
- Concurrency control

Tests can be found in `internal/infrastructure/usecase/car_impl_test.go` and `internal/infrastructure/opensearch/processor/processor_test.go`.

To run tests:

```bash
go test ./internal/infrastructure/usecase/...
go test ./internal/infrastructure/opensearch/processor/...
```

## Future Improvements

1. **Queue Integration**: Integrate with message queues like Redis, RabbitMQ, or Apache Kafka for more scalable and efficient processing
2. **Enhanced Monitoring**: Integration with Prometheus/Grafana for advanced metrics
3. **Message Priority**: Support for priority-based message processing
4. **Partitioning**: Message partitioning for better scalability
5. **Circuit Breaker**: Circuit breaker pattern for external service failures
6. **Advanced Retries**: More sophisticated retry strategies (jitter, etc.)

## References

- [Outbox Pattern - Microsoft Docs](https://docs.microsoft.com/en-us/azure/architecture/patterns/outbox)
- [Reliable Event Processing - Confluent](https://www.confluent.io/blog/reliable-event-processing/)
- [Transactional Outbox Pattern - Chris Richardson](https://microservices.io/patterns/data/transactional-outbox.html)
- [Distributed Systems Monitoring](https://sre.google/sre-book/monitoring-distributed-systems/)
