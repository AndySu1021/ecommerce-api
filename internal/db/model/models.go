// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0

package model

import (
	"database/sql"
	"time"

	product_vo "ecommerce-api/pkg/catalog/product/domain/vo"
	"ecommerce-api/pkg/constant"
	admin_vo "ecommerce-api/pkg/identity/admin/domain/vo"
	member_vo "ecommerce-api/pkg/identity/member/domain/vo"
)

// 管理員資料表
type Admin struct {
	ID uint64
	// 商戶ID
	MerchantID uint64
	// 信箱
	Email string
	// 密碼
	Password string
	// 真實姓名
	RealName string
	// 手機號
	Mobile string
	// 性別 1男 2女
	Sex admin_vo.Sex
	// 上次登入時間
	LastLoginTime sql.NullTime
	// 狀態 1開啟 2關閉
	IsEnabled constant.YesNo
	// 創建時間
	CreatedAt time.Time
	// 更新時間
	UpdatedAt time.Time
}

// 會員資料表
type Member struct {
	ID uint64
	// 商戶ID
	MerchantID uint64
	// 信箱
	Email string
	// 密碼
	Password string
	// 真實姓名
	RealName string
	// 手機號
	Mobile string
	// 性別 1男 2女
	Sex member_vo.Sex
	// 生日
	Birthday sql.NullTime
	// 城市
	City string
	// 區域
	District string
	// 剩餘地址
	Address string
	// 郵遞區號
	ZipCode string
	// 上次登入時間
	LastLoginTime sql.NullTime
	// 狀態 1開啟 2關閉
	IsEnabled constant.YesNo
	// 創建時間
	CreatedAt time.Time
	// 更新時間
	UpdatedAt time.Time
}

// 商戶資料表
type Merchant struct {
	ID uint64
	// 名稱
	Name string
	// 編號
	Code string
	// 域名
	Host string
	// 加密鹽
	EncryptSalt string
	// 是否啟用 0否 1是
	IsEnabled int32
	// 創建時間
	CreatedAt time.Time
	// 更新時間
	UpdatedAt time.Time
}

// 商品資料表
type Product struct {
	ID uint64
	// 商戶ID
	MerchantID uint64
	// 名稱
	Name string
	// 商品分類ID
	CategoryID uint64
	// 幣別 1新台幣
	CurrencyID constant.Currency
	// 價格
	Price uint64
	// 優惠價格
	SpecialPrice uint64
	// 優惠價格開始時間
	SpecialPriceStart sql.NullTime
	// 優惠價格結束時間
	SpecialPriceEnd sql.NullTime
	// 單筆訂單限購數量
	SingleOrderLimit int32
	// 是否一個商品成立一筆訂單 0否 1是
	IsSingleOrderOnly constant.YesNo
	// 溫層 1常溫 2冷藏 3冷凍
	Temperature product_vo.Temperature
	// 長度
	Length int32
	// 寬度
	Width int32
	// 高度
	Height int32
	// 重量(公克)
	Weight int32
	// 支援配送方式 1宅配到府 2超商取貨
	SupportDeliveryMethod product_vo.SupportDeliveryMethod
	// 是否為航空禁運品 0否 1是
	IsAirContraband constant.YesNo
	// 商品描述
	Description string
	// 商品圖片
	Pictures product_vo.PictureArray
	// 其他
	Extra product_vo.Extra
	// 是否啟用 0下架 1上架
	IsEnabled constant.YesNo
	// 總銷售量(每日更新)
	Sales uint64
	// 創建時間
	CreatedAt time.Time
	// 更新時間
	UpdatedAt time.Time
}

// 商品分類資料表
type ProductCategory struct {
	ID uint64
	// 商戶ID
	MerchantID uint64
	// 名稱
	Name string
	// 頂層ID
	TopID uint64
	// 父級ID
	ParentID uint64
	// 代理數節點左編號
	TreeLeft uint64
	// 代理數節點右編號
	TreeRight uint64
	// 創建時間
	CreatedAt time.Time
	// 更新時間
	UpdatedAt time.Time
}

// 商品規格資料表
type ProductSpec struct {
	ID uint64
	// 商戶ID
	MerchantID uint64
	// 商品ID
	ProductID uint64
	// 規格層級
	Level int32
	// 類型 1規格標題 2規格選項
	Type int32
	// 名稱
	Name string
	// 創建時間
	CreatedAt time.Time
	// 更新時間
	UpdatedAt time.Time
}

// 商品庫存資料表
type ProductStock struct {
	ID uint64
	// 商戶ID
	MerchantID uint64
	// 商品ID
	ProductID uint64
	// 第一層規格ID
	Spec1ID uint64
	// 第二層規格ID
	Spec2ID uint64
	// 庫存數量
	Quantity int32
	// SKU 貨號
	Code string
	// 創建時間
	CreatedAt time.Time
	// 更新時間
	UpdatedAt time.Time
}
