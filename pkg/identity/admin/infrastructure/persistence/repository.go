package persistence

import (
	"context"
	"database/sql"
	"ecommerce-api/internal/db/model"
	"ecommerce-api/pkg/identity/admin/domain/entity"
	"ecommerce-api/pkg/identity/admin/domain/vo"
	"github.com/Masterminds/squirrel"
)

type AdminRepository struct {
	db      *sql.DB
	queries *model.Queries
}

func NewAdminRepository(db *sql.DB, queries *model.Queries) *AdminRepository {
	return &AdminRepository{
		db:      db,
		queries: queries,
	}
}

func (r *AdminRepository) UpdateAdminByMap(ctx context.Context, adminId uint64, values map[string]interface{}) error {
	builder := squirrel.Update("admin")

	for k, v := range values {
		builder = builder.Set(k, v)
	}

	builder = builder.Where(squirrel.Eq{"id": adminId}).RunWith(r.db)

	if _, err := builder.ExecContext(ctx); err != nil {
		return err
	}

	return nil
}

func (r *AdminRepository) AdminLogin(ctx context.Context, params vo.AdminLoginParams) (entity.Admin, error) {
	admin, err := r.queries.AdminLogin(ctx, model.AdminLoginParams{
		MerchantID: params.MerchantID,
		Email:      params.Email,
		Password:   params.Password,
	})
	if err != nil {
		return entity.Admin{}, err
	}

	return entity.Admin{
		ID:            admin.ID,
		MerchantID:    admin.MerchantID,
		Email:         admin.Email,
		RealName:      admin.RealName,
		Mobile:        admin.Mobile,
		Sex:           admin.Sex,
		LastLoginTime: admin.LastLoginTime,
		IsEnabled:     admin.IsEnabled,
		CreatedAt:     admin.CreatedAt,
	}, nil
}
