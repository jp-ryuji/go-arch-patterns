package repository

import (
	"context"

	"github.com/jp-ryuji/go-arch-patterns/internal/domain/repository"
	"github.com/jp-ryuji/go-arch-patterns/internal/infrastructure/postgres/entgen"
)

type transactionManager struct {
	client *entgen.Client
}

// NewTransactionManager creates a new transaction manager
func NewTransactionManager(client *entgen.Client) repository.TransactionManager {
	return &transactionManager{
		client: client,
	}
}

// BeginTx starts a new transaction
func (tm *transactionManager) BeginTx(ctx context.Context) (*entgen.Tx, error) {
	return tm.client.Tx(ctx)
}

// CommitTx commits a transaction
func (tm *transactionManager) CommitTx(ctx context.Context, tx *entgen.Tx) error {
	return tx.Commit()
}

// RollbackTx rolls back a transaction
func (tm *transactionManager) RollbackTx(ctx context.Context, tx *entgen.Tx) error {
	return tx.Rollback()
}
