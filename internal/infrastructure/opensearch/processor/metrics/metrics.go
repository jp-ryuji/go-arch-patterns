package metrics

import (
	"sync/atomic"
	"time"
)

// OutboxMetrics tracks metrics for the outbox processor
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
	databaseErrors      atomic.Int64
	opensearchErrors    atomic.Int64
	serializationErrors atomic.Int64

	// Queue metrics
	pendingMessages atomic.Int64
	failedMessages  atomic.Int64
}

// NewOutboxMetrics creates a new metrics tracker
func NewOutboxMetrics() *OutboxMetrics {
	m := &OutboxMetrics{}
	// Initialize min to max int64
	m.minProcessingTime.Store(int64(^uint64(0) >> 1))
	return m
}

// RecordMessageProcessed records a successfully processed message
func (m *OutboxMetrics) RecordMessageProcessed(duration time.Duration, retries int) {
	m.messagesProcessed.Add(1)
	m.messagesSucceeded.Add(1)

	// Update timing metrics
	durationNanos := int64(duration)
	m.totalProcessingTime.Add(durationNanos)

	// Update min/max timing
	for {
		currentMin := m.minProcessingTime.Load()
		if durationNanos < currentMin {
			if m.minProcessingTime.CompareAndSwap(currentMin, durationNanos) {
				break
			}
		} else {
			break
		}
	}

	for {
		currentMax := m.maxProcessingTime.Load()
		if durationNanos > currentMax {
			if m.maxProcessingTime.CompareAndSwap(currentMax, durationNanos) {
				break
			}
		} else {
			break
		}
	}

	// Update retry metrics
	m.totalRetries.Add(int64(retries))
	for {
		currentMaxRetries := m.maxRetries.Load()
		if int64(retries) > currentMaxRetries {
			if m.maxRetries.CompareAndSwap(currentMaxRetries, int64(retries)) {
				break
			}
		} else {
			break
		}
	}
}

// RecordMessageFailed records a failed message
func (m *OutboxMetrics) RecordMessageFailed(errorType string) {
	m.messagesProcessed.Add(1)
	m.messagesFailed.Add(1)

	switch errorType {
	case "database":
		m.databaseErrors.Add(1)
	case "opensearch":
		m.opensearchErrors.Add(1)
	case "serialization":
		m.serializationErrors.Add(1)
	}
}

// RecordBatchProcessed records a processed batch
func (m *OutboxMetrics) RecordBatchProcessed(batchSize int) {
	m.batchesProcessed.Add(1)
	// Simple average calculation - in a real implementation, you might want a more sophisticated approach
	m.avgBatchSize.Store(int64(batchSize))
}

// UpdateQueueMetrics updates queue-related metrics
func (m *OutboxMetrics) UpdateQueueMetrics(pending, failed int64) {
	m.pendingMessages.Store(pending)
	m.failedMessages.Store(failed)
}

// GetMetrics returns a snapshot of current metrics
func (m *OutboxMetrics) GetMetrics() map[string]interface{} {
	return map[string]interface{}{
		"messages_processed_total":   m.messagesProcessed.Load(),
		"messages_succeeded_total":   m.messagesSucceeded.Load(),
		"messages_failed_total":      m.messagesFailed.Load(),
		"processing_time_avg_ms":     float64(m.totalProcessingTime.Load()) / float64(m.messagesProcessed.Load()+1) / float64(time.Millisecond),
		"processing_time_min_ms":     float64(m.minProcessingTime.Load()) / float64(time.Millisecond),
		"processing_time_max_ms":     float64(m.maxProcessingTime.Load()) / float64(time.Millisecond),
		"batches_processed_total":    m.batchesProcessed.Load(),
		"avg_batch_size":             m.avgBatchSize.Load(),
		"total_retries":              m.totalRetries.Load(),
		"max_retries":                m.maxRetries.Load(),
		"database_errors_total":      m.databaseErrors.Load(),
		"opensearch_errors_total":    m.opensearchErrors.Load(),
		"serialization_errors_total": m.serializationErrors.Load(),
		"pending_messages":           m.pendingMessages.Load(),
		"failed_messages":            m.failedMessages.Load(),
	}
}

// Reset resets all metrics
func (m *OutboxMetrics) Reset() {
	m.messagesProcessed.Store(0)
	m.messagesFailed.Store(0)
	m.messagesSucceeded.Store(0)
	m.totalProcessingTime.Store(0)
	m.minProcessingTime.Store(int64(^uint64(0) >> 1))
	m.maxProcessingTime.Store(0)
	m.batchesProcessed.Store(0)
	m.avgBatchSize.Store(0)
	m.totalRetries.Store(0)
	m.maxRetries.Store(0)
	m.databaseErrors.Store(0)
	m.opensearchErrors.Store(0)
	m.serializationErrors.Store(0)
	m.pendingMessages.Store(0)
	m.failedMessages.Store(0)
}
