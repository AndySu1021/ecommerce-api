package vo

import (
	"database/sql"
	"ecommerce-api/pkg/constant"
)

type CurrentMember struct {
	ID       uint64   `json:"id"`
	Email    string   `json:"email"`
	Token    string   `json:"token"`
	Merchant Merchant `json:"merchant"`
}

type MemberInfo struct {
	RealName string       `json:"real_name"`
	Email    string       `json:"email"`
	Mobile   string       `json:"mobile"`
	Sex      Sex          `json:"sex"`
	Birthday sql.NullTime `json:"birthday"`
	City     string       `json:"city"`
	District string       `json:"district"`
	Address  string       `json:"address"`
	ZipCode  string       `json:"zip_code"`
}

type ListMemberRow struct {
	ID        uint64         `json:"id"`
	RealName  string         `json:"real_name"`
	Email     string         `json:"email"`
	Mobile    string         `json:"mobile"`
	Birthday  sql.NullTime   `json:"birthday"`
	IsEnabled constant.YesNo `json:"is_enabled"`
}
