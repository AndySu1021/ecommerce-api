package entity

import (
	"database/sql"
	"ecommerce-api/pkg/constant"
	"ecommerce-api/pkg/identity/admin/domain/vo"
	"time"
)

// Admin Aggregate Root
type Admin struct {
	ID            uint64         `json:"id"`
	MerchantID    uint64         `json:"merchant_id"`
	Email         string         `json:"email"`
	Password      string         `json:"password,omitempty"`
	RealName      string         `json:"real_name"`
	Mobile        string         `json:"mobile"`
	Sex           vo.Sex         `json:"sex"`
	LastLoginTime sql.NullTime   `json:"last_login_time"`
	IsEnabled     constant.YesNo `json:"is_enabled"`
	CreatedAt     time.Time      `json:"created_at"`
}
