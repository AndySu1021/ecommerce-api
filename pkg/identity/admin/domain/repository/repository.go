package repository

import (
	"context"
	"ecommerce-api/pkg/identity/admin/domain/entity"
	"ecommerce-api/pkg/identity/admin/domain/vo"
)

type AdminRepository interface {
	AdminLogin(ctx context.Context, params vo.AdminLoginParams) (entity.Admin, error)
	UpdateAdminByMap(ctx context.Context, adminId uint64, values map[string]interface{}) error
}
