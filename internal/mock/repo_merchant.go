package mock

import (
	"ecommerce-api/pkg/identity/merchant/domain/entity"
	"ecommerce-api/pkg/identity/merchant/domain/repository"
	"github.com/golang/mock/gomock"
	"testing"
)

func NewMockMerchantRepo(t *testing.T) repository.MerchantRepository {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	merchant := NewMockMerchantRepository(ctl)
	merchant.EXPECT().GetMerchantByHost(gomock.Any(), gomock.Any()).AnyTimes().Return(entity.Merchant{
		ID:          1,
		Name:        "Test",
		Code:        "123456",
		Host:        "a.b.c",
		EncryptSalt: "abcd12345",
	}, nil)

	return merchant
}
