package vo

import (
	"database/sql"
	"ecommerce-api/pkg/constant"
)

type RegisterParams struct {
	Email               string `json:"email" binding:"required,email"`
	Password            string `json:"password" binding:"required"`
	MerchantID          uint64
	MerchantHost        string
	MerchantEncryptSalt string
}

type LoginParams struct {
	Email               string `json:"email" binding:"required"`
	Password            string `json:"password" binding:"required"`
	MerchantID          uint64
	MerchantHost        string
	MerchantEncryptSalt string
}

type LogoutParams struct {
	Email        string
	MerchantHost string
}

type ForgetPasswordParams struct {
	Email        string `json:"email" binding:"required"`
	MerchantID   uint64
	MerchantName string
	MerchantHost string
}

type CheckForgetCodeParams struct {
	Email        string `json:"email" binding:"required"`
	Code         string `json:"code" binding:"required"`
	MerchantHost string
}

type ResetPasswordParams struct {
	Code                string `json:"code" binding:"required"`
	Password            string `json:"password" binding:"required"`
	MerchantID          uint64
	MerchantHost        string
	MerchantEncryptSalt string
}

type UpdateMemberInfoParams struct {
	RealName     string `json:"real_name" binding:"required"`
	Mobile       string `json:"mobile" binding:""`
	Birthday     string `json:"birthday" binding:"required"`
	Sex          Sex    `json:"sex" binding:"required,oneof=1 2"`
	City         string `json:"city" binding:""`
	District     string `json:"district" binding:""`
	Address      string `json:"address" binding:""`
	ZipCode      string `json:"zip_code" binding:""`
	BirthdayTime sql.NullTime
	MemberID     uint64
}

type UpdateMemberPasswordParams struct {
	OldPassword         string `json:"old_password" binding:"required"`
	Password            string `json:"password" binding:"required"`
	MemberID            uint64
	MerchantEncryptSalt string
}

type ListMemberParams struct {
	IsEnabled constant.YesNo `form:"is_enabled" binding:"oneof=0 1"`
	constant.Pagination
}

type UpdateMemberStatusParams struct {
	IsEnabled constant.YesNo `json:"is_enabled" binding:"oneof=0 1"`
	MemberID  uint64
}
