package application

import (
	"context"
	"ecommerce-api/pkg/constant"
	"ecommerce-api/pkg/identity/merchant/domain/entity"
	"ecommerce-api/pkg/identity/merchant/domain/repository"
	"fmt"
	"github.com/go-redis/redis/v8"
)

type MerchantService interface {
	GetCurrentMerchant(ctx context.Context) (entity.Merchant, error)
	GetMerchantByHost(ctx context.Context, host string) (entity.Merchant, error)
}

type service struct {
	rdb  redis.UniversalClient
	repo repository.MerchantRepository
}

func NewMerchantService(rdb redis.UniversalClient, repo repository.MerchantRepository) MerchantService {
	return &service{
		rdb:  rdb,
		repo: repo,
	}
}

func (s *service) GetCurrentMerchant(ctx context.Context) (entity.Merchant, error) {
	host, ok := ctx.Value(constant.ContextKeyHost).(string)
	if host == "" || !ok {
		return entity.Merchant{}, fmt.Errorf("invalid host")
	}

	return s.GetMerchantByHost(ctx, host)
}

func (s *service) GetMerchantByHost(ctx context.Context, host string) (entity.Merchant, error) {
	return s.repo.GetMerchantByHost(ctx, host)
}
