package application

import (
	"context"
	"ecommerce-api/internal/mock"
	"ecommerce-api/pkg/identity/merchant/domain/entity"
	"ecommerce-api/pkg/identity/merchant/domain/repository"
	"fmt"
	"github.com/go-redis/redis/v8"
	"reflect"
	"testing"
)

func TestMain(m *testing.M) {
	fmt.Println("begin")
	m.Run()
	fmt.Println("end")
}

func Test_service_GetMerchantByHost(t *testing.T) {
	type fields struct {
		rdb  redis.UniversalClient
		repo repository.MerchantRepository
	}
	type args struct {
		ctx  context.Context
		host string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    entity.Merchant
		wantErr bool
	}{
		{
			name: "aaaa",
			fields: fields{
				rdb:  nil,
				repo: mock.NewMockMerchantRepo(t),
			},
			args: args{
				ctx:  context.Background(),
				host: "a.b.c",
			},
			want: entity.Merchant{
				ID:          0,
				Name:        "Test",
				Code:        "123456",
				Host:        "a.b.c",
				EncryptSalt: "abcd12345",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				rdb:  tt.fields.rdb,
				repo: tt.fields.repo,
			}
			got, err := s.GetMerchantByHost(tt.args.ctx, tt.args.host)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMerchantByHost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMerchantByHost() got = %v, want %v", got, tt.want)
			}
		})
	}
}
