package persistence

import (
	"context"
	"database/sql"
	"ecommerce-api/internal/db/model"
	"ecommerce-api/pkg/identity/merchant/domain/entity"
)

type MerchantRepository struct {
	db      *sql.DB
	tx      *sql.Tx
	queries *model.Queries
}

func NewMerchantRepository(db *sql.DB, queries *model.Queries) *MerchantRepository {
	return &MerchantRepository{
		db:      db,
		tx:      nil,
		queries: queries,
	}
}

func (r *MerchantRepository) GetMerchantByHost(ctx context.Context, host string) (entity.Merchant, error) {
	return entity.Merchant{
		ID:          1,
		Name:        "X電商",
		Host:        "x.ecommerce.com",
		EncryptSalt: "abcde12345",
	}, nil

	merchant, err := r.queries.GetMerchantByHost(ctx, host)
	if err != nil {
		return entity.Merchant{}, err
	}

	return entity.Merchant{
		ID:          merchant.ID,
		Name:        merchant.Name,
		Host:        merchant.Host,
		EncryptSalt: merchant.EncryptSalt,
	}, nil
}
