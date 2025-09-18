# Outbox Pattern Implementation

This document explains the implementation of the Outbox pattern in this project, which ensures data consistency between PostgreSQL and external systems.

## Overview

The Outbox pattern is a technique for reliably publishing events from a service to external systems. It addresses the challenge of ensuring that data changes in the primary database are reliably propagated to external systems usually via message queues.

In our implementation, we coordinate between:

1. PostgreSQL (primary database)
2. External systems (such as search engines, messaging systems, etc.)

## Implementation Details

### Key Files

1. **Domain Layer**:
   - `internal/domain/repository/outbox.go` - Outbox repository interface
   - `internal/domain/repository/transaction.go` - Transaction manager interface

2. **Infrastructure Layer**:
   - `internal/infrastructure/postgres/ent/schema/outbox.go` - Outbox table schema
   - `internal/infrastructure/postgres/repository/outbox_repository.go` - Outbox repository implementation
   - `internal/infrastructure/postgres/repository/transaction_manager.go` - Transaction manager implementation

3. **Application Layer**:
   - `internal/usecase/car_impl.go` - Car usecase implementation with outbox pattern
   - `internal/usecase/car.go` - Car usecase interface

4. **Tests**:
   - `internal/usecase/car_impl_test.go` - Unit tests for car usecase with transactional outbox

### Outbox Flow

The outbox pattern is implemented in the car usecase:

1. When a car is registered, a database transaction is started
2. The car data is saved to PostgreSQL within the transaction
3. An outbox message is created for external systems within the same transaction
4. The transaction is committed (ensuring atomicity)

This approach ensures that either both the car data and the outbox message are saved, or neither is, maintaining consistency between the database and the outbox table.

## Transactional Guarantees

The implementation ensures atomicity between the main entity creation and outbox message creation by using database transactions. Both operations happen within the same transaction, so either both succeed or both fail.

## Concurrency Control

The implementation uses `SELECT ... FOR UPDATE SKIP LOCKED` for efficient concurrency control when processing outbox messages. This approach:

1. **Prevents Race Conditions**: Uses database-level locking to ensure only one processor can claim a message
2. **Improves Performance**: Eliminates the need for separate SELECT and UPDATE operations
3. **Handles Failures Gracefully**: Automatically skips locked messages that are being processed by other instances
4. **Ensures FIFO Processing**: Messages are processed in first-in-first-out order

The implementation uses Ent's query builder with the `ForUpdate` method and `SkipLocked` option, which requires enabling the `sql/lock` feature flag during Ent code generation:

```go
pendingMessages, err := r.client.Outbox.Query().
    Where(outbox.Status("pending")).
    Limit(limit).
    Order(entgen.Asc(outbox.FieldCreatedAt)).
    ForUpdate(
        sql.WithLockAction(sql.SkipLocked),
    ).
    All(ctx)
```

This generates SQL equivalent to:

```sql
SELECT * FROM outbox 
WHERE status = 'pending' 
ORDER BY created_at ASC 
LIMIT $1 
FOR UPDATE SKIP LOCKED
```

To enable this functionality, the Ent code generation includes the `sql/lock` feature flag:

```go
//go:generate go run -mod=mod entgo.io/ent/cmd/ent generate --target ../entgen --feature sql/lock ./schema
```

The `sql/lock` feature flag is sufficient for this implementation. The `sql/execquery` feature flag is not required but may be included for future enhancements that require direct SQL execution.

## Schema Evolution

The outbox schema is designed to support schema evolution:

1. **Flexible Payload**: JSON payload field can accommodate different message formats
2. **Version Field**: Message versioning for tracking schema changes
3. **Extensible Fields**: Additional fields can be added without breaking existing code

## Future Improvements

This repository focuses on demonstrating the core concept of the outbox pattern: saving records to the outbox table within a database transaction to ensure atomicity. The actual processing of those messages is suggested as a future improvement.

In a production environment, this would typically be handled by external services. For example:

1. **AWS EventBridge Scheduler** could trigger a Lambda function at regular intervals to process records in the outbox table
2. The Lambda function would read pending messages from the outbox table and send them to **AWS SQS FIFO** queues
   - **AWS SQS (FIFO)** - Amazon's managed message queue service that ensures message ordering and exactly-once processing
3. Separate **Lambda functions triggered by SQS FIFO** queues would then process these messages and update external systems

This approach provides a scalable, serverless solution for processing outbox messages with guaranteed ordering and exactly-once delivery semantics.

The outbox table records would be processed by these external services, which would:

1. Poll or be triggered by the outbox table for new messages
2. Send the messages to external systems
3. Mark the messages as processed in the outbox table
4. Handle errors and retries as appropriate for the specific service

Additional improvements could include:

1. **Enhanced Monitoring**: Integration with Prometheus/Grafana for advanced metrics
2. **Message Priority**: Support for priority-based message processing
3. **Partitioning**: Message partitioning for better scalability
4. **Circuit Breaker**: Circuit breaker pattern for external service failures
5. **Dead Letter Queue (DLQ) Processing**: Implementation of a dead letter queue pattern for handling messages that fail repeatedly processing. Messages that exceed a configurable retry threshold could be moved to a separate "dead_letter" status in the outbox table, allowing for:
   - Isolation of problematic messages to prevent blocking processing of other messages
   - Manual inspection and resolution of failed messages
   - Separate monitoring and alerting for dead letter messages
   - Potential replay mechanisms for messages after issues are resolved

## References

- [Outbox Pattern - Microsoft Docs](https://docs.microsoft.com/en-us/azure/architecture/patterns/outbox)
- [Reliable Event Processing - Confluent](https://www.confluent.io/blog/reliable-event-processing/)
- [Transactional Outbox Pattern - Chris Richardson](https://microservices.io/patterns/data/transactional-outbox.html)
- [Distributed Systems Monitoring](https://sre.google/sre-book/monitoring-distributed-systems/)
