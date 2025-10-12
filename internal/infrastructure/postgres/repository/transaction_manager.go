package repository

import (
	"context"
	"fmt"

	sharedrepository "github.com/jp-ryuji/go-arch-patterns/internal/domain/repository"
	"github.com/jp-ryuji/go-arch-patterns/internal/infrastructure/postgres/entgen"
)

type transactionManager struct {
	client *entgen.Client
}

// NewTransactionManager creates a new transaction manager
func NewTransactionManager(client *entgen.Client) sharedrepository.TransactionManager {
	return &transactionManager{
		client: client,
	}
}

// BeginTx starts a new transaction
func (tm *transactionManager) BeginTx(ctx context.Context) (sharedrepository.Transaction, error) {
	tx, err := tm.client.Tx(ctx)
	if err != nil {
		return nil, err
	}

	return &entTransaction{tx: tx}, nil
}

// CommitTx commits a transaction
func (tm *transactionManager) CommitTx(ctx context.Context, tx sharedrepository.Transaction) error {
	entTx, err := unwrapEntTx(tx)
	if err != nil {
		return err
	}

	return entTx.Commit()
}

// RollbackTx rolls back a transaction
func (tm *transactionManager) RollbackTx(ctx context.Context, tx sharedrepository.Transaction) error {
	entTx, err := unwrapEntTx(tx)
	if err != nil {
		return err
	}

	return entTx.Rollback()
}

type entTransaction struct {
	tx *entgen.Tx
}

func (t *entTransaction) Raw() any {
	return t.tx
}

func unwrapEntTx(tx sharedrepository.Transaction) (*entgen.Tx, error) {
	if tx == nil {
		return nil, fmt.Errorf("transaction is nil")
	}

	entTx, ok := tx.Raw().(*entgen.Tx)
	if !ok {
		return nil, fmt.Errorf("unsupported transaction type %T", tx.Raw())
	}

	return entTx, nil
}
