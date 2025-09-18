package repository

import (
	"context"

	"github.com/jp-ryuji/go-arch-patterns/internal/infrastructure/postgres/entgen"
)

//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock_repository
type TransactionManager interface {
	BeginTx(ctx context.Context) (*entgen.Tx, error)
	CommitTx(ctx context.Context, tx *entgen.Tx) error
	RollbackTx(ctx context.Context, tx *entgen.Tx) error
}
