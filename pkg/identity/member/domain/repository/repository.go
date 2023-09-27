package repository

import (
	"context"
	"ecommerce-api/pkg/identity/member/domain/entity"
	"ecommerce-api/pkg/identity/member/domain/vo"
)

type MemberRepository interface {
	CreateMember(ctx context.Context, params vo.CreateMemberParams) error
	MemberLogin(ctx context.Context, params vo.MemberLoginParams) (entity.Member, error)
	UpdateMemberByMap(ctx context.Context, memberId uint64, values map[string]interface{}) error
	ListMember(ctx context.Context, params vo.ListMemberParams) ([]vo.ListMemberRow, int64, error)
	GetMemberInfo(ctx context.Context, memberId uint64) (vo.MemberInfo, error)
	UpdateMemberPassword(ctx context.Context, params vo.UpdateMemberPasswordRepoParams) error
	ResetPassword(ctx context.Context, params vo.ResetPasswordRepoParams) error
	CheckEmailExist(ctx context.Context, merchantId uint64, email string) (int32, error)
}
