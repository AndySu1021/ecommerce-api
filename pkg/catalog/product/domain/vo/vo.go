package vo

import (
	"database/sql"
	"ecommerce-api/pkg/constant"
)

type ListProductCategoryRow struct {
	ID       int64
	Name     string
	ParentID int64
	Level    int32
}

type ListAgentProductCategoryRow struct {
	ID              int64
	Name            string
	ParentID        int64
	Level           int32
	ProductQuantity int32
}

type DeleteProductCategoryParams struct {
	TopID     int64
	TreeLeft  int64
	TreeRight int64
}

type ProductSpecItem struct {
	ID   uint64 `json:"id" binding:"gte=0"`
	Name string `json:"name" binding:"required"`
}

type ListProductRow struct {
	ID                 uint64            `json:"id"`
	Name               string            `json:"name"`
	CurrencyID         constant.Currency `json:"currency_id"`
	Price              int64             `json:"price"`
	SpecialPrice       int64             `json:"special_price"`
	SpecialPriceStart  sql.NullTime      `json:"special_price_start"`
	SpecialPriceEnd    sql.NullTime      `json:"special_price_end"`
	Pictures           []string          `json:"pictures"`
	IsEnabled          constant.YesNo    `json:"is_enabled"`
	TotalStockQuantity int32             `json:"total_stock_quantity"`
}
