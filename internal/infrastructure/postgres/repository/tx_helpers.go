package repository

import (
	"fmt"

	sharedrepository "github.com/jp-ryuji/go-arch-patterns/internal/domain/repository"
	"github.com/jp-ryuji/go-arch-patterns/internal/infrastructure/postgres/entgen"
)

func asEntTx(tx sharedrepository.Transaction) (*entgen.Tx, error) {
	if tx == nil {
		return nil, fmt.Errorf("transaction is nil")
	}

	raw := tx.Raw()
	entTx, ok := raw.(*entgen.Tx)
	if !ok {
		return nil, fmt.Errorf("unsupported transaction type %T", raw)
	}

	return entTx, nil
}
