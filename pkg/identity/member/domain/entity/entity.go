package entity

import (
	"database/sql"
	"ecommerce-api/pkg/constant"
	"ecommerce-api/pkg/identity/member/domain/vo"
	"time"
)

// Member Aggregate Root
type Member struct {
	ID            uint64         `json:"id"`
	MerchantID    uint64         `json:"merchant_id"`
	Email         string         `json:"email"`
	Password      string         `json:"password,omitempty"`
	RealName      string         `json:"real_name"`
	Mobile        string         `json:"mobile"`
	Sex           vo.Sex         `json:"sex"`
	Birthday      sql.NullTime   `json:"birthday"`
	City          string         `json:"city"`
	District      string         `json:"district"`
	Address       string         `json:"address"`
	ZipCode       string         `json:"zip_code"`
	LastLoginTime sql.NullTime   `json:"last_login_time"`
	IsEnabled     constant.YesNo `json:"is_enabled"`
	CreatedAt     time.Time      `json:"created_at"`
}
