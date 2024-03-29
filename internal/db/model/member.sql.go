// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: member.sql

package model

import (
	"context"
	"database/sql"
	"time"

	member_vo "ecommerce-api/pkg/identity/member/domain/vo"
)

const checkEmailExist = `-- name: CheckEmailExist :one
SELECT 1
FROM member
WHERE merchant_id = ?
  AND email = ?
`

type CheckEmailExistParams struct {
	MerchantID uint64
	Email      string
}

func (q *Queries) CheckEmailExist(ctx context.Context, arg CheckEmailExistParams) (int32, error) {
	row := q.queryRow(ctx, q.checkEmailExistStmt, checkEmailExist, arg.MerchantID, arg.Email)
	var column_1 int32
	err := row.Scan(&column_1)
	return column_1, err
}

const createMember = `-- name: CreateMember :exec
INSERT INTO member (merchant_id, email, password, created_at, updated_at)
VALUES (?, ?, ?, ?, ?)
`

type CreateMemberParams struct {
	MerchantID uint64
	Email      string
	Password   string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (q *Queries) CreateMember(ctx context.Context, arg CreateMemberParams) error {
	_, err := q.exec(ctx, q.createMemberStmt, createMember,
		arg.MerchantID,
		arg.Email,
		arg.Password,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	return err
}

const getMemberInfo = `-- name: GetMemberInfo :one
SELECT real_name,
       email,
       mobile,
       sex,
       birthday,
       city,
       district,
       address,
       zip_code
FROM member
WHERE id = ?
`

type GetMemberInfoRow struct {
	RealName string
	Email    string
	Mobile   string
	Sex      member_vo.Sex
	Birthday sql.NullTime
	City     string
	District string
	Address  string
	ZipCode  string
}

func (q *Queries) GetMemberInfo(ctx context.Context, id uint64) (GetMemberInfoRow, error) {
	row := q.queryRow(ctx, q.getMemberInfoStmt, getMemberInfo, id)
	var i GetMemberInfoRow
	err := row.Scan(
		&i.RealName,
		&i.Email,
		&i.Mobile,
		&i.Sex,
		&i.Birthday,
		&i.City,
		&i.District,
		&i.Address,
		&i.ZipCode,
	)
	return i, err
}

const memberLogin = `-- name: MemberLogin :one
SELECT id, merchant_id, email, password, real_name, mobile, sex, birthday, city, district, address, zip_code, last_login_time, is_enabled, created_at, updated_at
FROM member
WHERE merchant_id = ?
  AND email = ?
  AND password = ?
`

type MemberLoginParams struct {
	MerchantID uint64
	Email      string
	Password   string
}

func (q *Queries) MemberLogin(ctx context.Context, arg MemberLoginParams) (Member, error) {
	row := q.queryRow(ctx, q.memberLoginStmt, memberLogin, arg.MerchantID, arg.Email, arg.Password)
	var i Member
	err := row.Scan(
		&i.ID,
		&i.MerchantID,
		&i.Email,
		&i.Password,
		&i.RealName,
		&i.Mobile,
		&i.Sex,
		&i.Birthday,
		&i.City,
		&i.District,
		&i.Address,
		&i.ZipCode,
		&i.LastLoginTime,
		&i.IsEnabled,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const resetPassword = `-- name: ResetPassword :exec
UPDATE member
SET password = ?
WHERE merchant_id = ?
  AND email = ?
`

type ResetPasswordParams struct {
	Password   string
	MerchantID uint64
	Email      string
}

func (q *Queries) ResetPassword(ctx context.Context, arg ResetPasswordParams) error {
	_, err := q.exec(ctx, q.resetPasswordStmt, resetPassword, arg.Password, arg.MerchantID, arg.Email)
	return err
}

const updateMemberPassword = `-- name: UpdateMemberPassword :exec
UPDATE member
SET password   = ?,
    updated_at = ?
WHERE id = ?
  AND password = ?
`

type UpdateMemberPasswordParams struct {
	Password    string
	UpdatedAt   time.Time
	ID          uint64
	OldPassword string
}

func (q *Queries) UpdateMemberPassword(ctx context.Context, arg UpdateMemberPasswordParams) error {
	_, err := q.exec(ctx, q.updateMemberPasswordStmt, updateMemberPassword,
		arg.Password,
		arg.UpdatedAt,
		arg.ID,
		arg.OldPassword,
	)
	return err
}
