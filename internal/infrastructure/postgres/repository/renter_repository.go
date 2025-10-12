package repository

import (
	"context"
	"strings"
	"time"

	"github.com/jp-ryuji/go-arch-patterns/internal/infrastructure/postgres/entgen"
	rentertable "github.com/jp-ryuji/go-arch-patterns/internal/infrastructure/postgres/entgen/renter"
	renter "github.com/jp-ryuji/go-arch-patterns/internal/renter"
)

type renterRepository struct {
	client *entgen.Client
}

func NewRenterRepository(client *entgen.Client) renter.RenterRepository {
	return &renterRepository{client: client}
}

func (r *renterRepository) Create(ctx context.Context, renterModel *renter.Renter) error {
	_, err := r.client.Renter.
		Create().
		SetID(renterModel.ID).
		SetTenantID(renterModel.TenantID).
		SetType(string(renterModel.Type)).
		SetCreatedAt(renterModel.CreatedAt).
		SetUpdatedAt(renterModel.UpdatedAt).
		Save(ctx)
	return err
}

func (r *renterRepository) GetByID(ctx context.Context, id string) (*renter.Renter, error) {
	renterDB, err := r.client.Renter.
		Query().
		Where(rentertable.ID(id)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	var renterType renter.Type
	switch renterDB.Type {
	case string(renter.TypeCompany):
		renterType = renter.TypeCompany
	case string(renter.TypeIndividual):
		renterType = renter.TypeIndividual
	default:
		renterType = ""
	}

	return &renter.Renter{
		ID:        renterDB.ID,
		TenantID:  renterDB.TenantID,
		Type:      renterType,
		CreatedAt: renterDB.CreatedAt,
		UpdatedAt: renterDB.UpdatedAt,
	}, nil
}

func (r *renterRepository) Update(ctx context.Context, renterModel *renter.Renter) error {
	renterModel.UpdatedAt = time.Now()

	_, err := r.client.Renter.
		UpdateOneID(renterModel.ID).
		SetTenantID(renterModel.TenantID).
		SetType(string(renterModel.Type)).
		SetUpdatedAt(renterModel.UpdatedAt).
		Save(ctx)
	return err
}

func (r *renterRepository) Delete(ctx context.Context, id string) error {
	err := r.client.Renter.
		DeleteOneID(id).
		Exec(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "no rows in result set") {
			return nil
		}
		return err
	}

	return nil
}
