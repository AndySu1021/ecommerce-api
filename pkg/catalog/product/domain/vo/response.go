package vo

import "ecommerce-api/pkg/constant"

type ProductCategoryItem struct {
	ID              int64                  `json:"id"`
	Label           string                 `json:"label"`
	ProductQuantity int32                  `json:"product_quantity"`
	Children        []*ProductCategoryItem `json:"children,omitempty"`
}

type ListProductRecord struct {
	ID                 uint64            `json:"id"`
	Name               string            `json:"name"`
	CurrencyID         constant.Currency `json:"currency_id"`
	Price              int64             `json:"price"`
	SpecialPrice       int64             `json:"special_price"`
	Pictures           []string          `json:"pictures"`
	PictureUrls        []string          `json:"picture_urls"`
	IsEnabled          constant.YesNo    `json:"is_enabled"`
	TotalStockQuantity int32             `json:"total_stock_quantity"`
	SpecItems          []string          `json:"spec_items"`
}
