package rpc

import (
	"context"
)

type Merchant struct {
	UnimplementedMerchantServer
}

func (m *Merchant) GetMerchantByHost(ctx context.Context, in *GetMerchantByHostRequest) (*GetMerchantByHostResponse, error) {
	return &GetMerchantByHostResponse{
		Host:        "a.b.c",
		Name:        "AAA",
		Code:        "123456",
		EncryptSalt: "abcde12345",
	}, nil
}
