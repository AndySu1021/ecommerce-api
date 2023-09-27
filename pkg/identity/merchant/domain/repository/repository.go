package repository

import (
	"context"
	"ecommerce-api/pkg/identity/merchant/domain/entity"
)

type MerchantRepository interface {
	GetMerchantByHost(ctx context.Context, host string) (entity.Merchant, error)
}
