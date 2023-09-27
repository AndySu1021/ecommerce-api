package persistence

import (
	"context"
	"database/sql"
	"ecommerce-api/internal/db/model"
	"ecommerce-api/internal/helper"
	"ecommerce-api/pkg/constant"
	"ecommerce-api/pkg/identity/member/domain/entity"
	"ecommerce-api/pkg/identity/member/domain/vo"
	"github.com/Masterminds/squirrel"
	"time"
)

type MemberRepository struct {
	db      *sql.DB
	queries *model.Queries
}

func NewMemberRepository(db *sql.DB, queries *model.Queries) *MemberRepository {
	return &MemberRepository{
		db:      db,
		queries: queries,
	}
}

func (r *MemberRepository) ListMember(ctx context.Context, params vo.ListMemberParams) (members []vo.ListMemberRow, total int64, err error) {
	columns := []string{
		"id",
		"real_name",
		"email",
		"mobile",
		"birthday",
		"is_enabled",
	}

	dQuery := squirrel.Select(columns...).From("member").
		Where(squirrel.Gt{"member.id": 1})
	cQuery := squirrel.Select("count(*) AS count").From("member").
		Where(squirrel.Gt{"member.id": 1})

	if params.IsEnabled != constant.All {
		dQuery = dQuery.Where(squirrel.Eq{"is_enabled": params.IsEnabled})
		cQuery = cQuery.Where(squirrel.Eq{"is_enabled": params.IsEnabled})
	}

	dQuery = dQuery.RunWith(r.db)
	cQuery = cQuery.RunWith(r.db)

	if err = helper.PageQuery(ctx, dQuery, params.Pagination, &members); err != nil {
		return
	}
	if err = helper.TotalQuery(ctx, cQuery, &total); err != nil {
		return
	}

	return
}

func (r *MemberRepository) UpdateMemberByMap(ctx context.Context, memberId uint64, values map[string]interface{}) error {
	builder := squirrel.Update("member")

	for k, v := range values {
		builder = builder.Set(k, v)
	}

	builder = builder.Where(squirrel.Eq{"id": memberId}).RunWith(r.db)

	if _, err := builder.ExecContext(ctx); err != nil {
		return err
	}

	return nil
}

func (r *MemberRepository) CreateMember(ctx context.Context, params vo.CreateMemberParams) error {
	now := time.Now().UTC()
	return r.queries.CreateMember(ctx, model.CreateMemberParams{
		MerchantID: params.MerchantID,
		Email:      params.Email,
		Password:   params.Password,
		CreatedAt:  now,
		UpdatedAt:  now,
	})
}

func (r *MemberRepository) MemberLogin(ctx context.Context, params vo.MemberLoginParams) (entity.Member, error) {
	member, err := r.queries.MemberLogin(ctx, model.MemberLoginParams{
		MerchantID: params.MerchantID,
		Email:      params.Email,
		Password:   params.Password,
	})
	if err != nil {
		return entity.Member{}, err
	}

	return entity.Member{
		ID:            member.ID,
		MerchantID:    member.MerchantID,
		Email:         member.Email,
		RealName:      member.RealName,
		Mobile:        member.Mobile,
		Sex:           member.Sex,
		Birthday:      member.Birthday,
		City:          member.City,
		District:      member.District,
		Address:       member.Address,
		ZipCode:       member.ZipCode,
		LastLoginTime: member.LastLoginTime,
		IsEnabled:     member.IsEnabled,
		CreatedAt:     member.CreatedAt,
	}, nil
}

func (r *MemberRepository) GetMemberInfo(ctx context.Context, memberId uint64) (vo.MemberInfo, error) {
	info, err := r.queries.GetMemberInfo(ctx, memberId)
	if err != nil {
		return vo.MemberInfo{}, err
	}

	return vo.MemberInfo{
		RealName: info.RealName,
		Email:    info.Email,
		Mobile:   info.Mobile,
		Sex:      info.Sex,
		Birthday: info.Birthday,
		City:     info.City,
		District: info.District,
		Address:  info.Address,
		ZipCode:  info.ZipCode,
	}, nil
}

func (r *MemberRepository) UpdateMemberPassword(ctx context.Context, params vo.UpdateMemberPasswordRepoParams) error {
	return r.queries.UpdateMemberPassword(ctx, model.UpdateMemberPasswordParams{
		Password:    params.Password,
		UpdatedAt:   time.Now().UTC(),
		ID:          params.MemberID,
		OldPassword: params.OldPassword,
	})
}

func (r *MemberRepository) ResetPassword(ctx context.Context, params vo.ResetPasswordRepoParams) error {
	return r.queries.ResetPassword(ctx, model.ResetPasswordParams{
		Password:   params.Password,
		MerchantID: params.MerchantID,
		Email:      params.Email,
	})
}

func (r *MemberRepository) CheckEmailExist(ctx context.Context, merchantId uint64, email string) (int32, error) {
	return r.queries.CheckEmailExist(ctx, model.CheckEmailExistParams{
		MerchantID: merchantId,
		Email:      email,
	})
}
