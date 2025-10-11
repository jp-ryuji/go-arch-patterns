package repository

import "context"

// Transaction represents an abstract unit of work boundary. Infrastructure
// implementations can expose driver specific handles through Raw(), but domain
// services should treat it as an opaque dependency.
type Transaction interface {
	Raw() any
}

//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock_repository
type TransactionManager interface {
	BeginTx(ctx context.Context) (Transaction, error)
	CommitTx(ctx context.Context, tx Transaction) error
	RollbackTx(ctx context.Context, tx Transaction) error
}
