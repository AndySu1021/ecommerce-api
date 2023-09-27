package entity

import (
	"ecommerce-api/pkg/catalog/product/domain/vo"
	"ecommerce-api/pkg/constant"
	"time"
)

// Product Aggregate Root
type Product struct {
	ID                    uint64                   `json:"id"`
	Name                  string                   `json:"name"`
	CurrencyID            constant.Currency        `json:"currency_id"`
	Price                 uint64                   `json:"price"`
	SpecialPrice          uint64                   `json:"special_price"`
	SpecialPriceStart     time.Time                `json:"special_price_start"`
	SpecialPriceEnd       time.Time                `json:"special_price_end"`
	CategoryID            uint64                   `json:"category_id"`
	CategoryName          string                   `json:"category_name"`
	Description           string                   `json:"description"`
	Pictures              vo.PictureArray          `json:"pictures"`
	PictureUrls           vo.PictureArray          `json:"picture_urls"`
	SingleOrderLimit      int32                    `json:"single_order_limit"`
	IsSingleOrderOnly     constant.YesNo           `json:"is_single_order_only"`
	Temperature           vo.Temperature           `json:"temperature"`
	Length                int32                    `json:"length"`
	Width                 int32                    `json:"width"`
	Height                int32                    `json:"height"`
	Weight                int32                    `json:"weight"`
	SupportDeliveryMethod vo.SupportDeliveryMethod `json:"support_delivery_method"`
	IsAirContraband       constant.YesNo           `json:"is_air_contraband"`
	Extra                 vo.Extra                 `json:"extra"`
	IsEnabled             constant.YesNo           `json:"is_enabled"`
	TotalStockQuantity    int32                    `json:"total_stock_quantity"` // 總共剩餘的庫存數量
	ProductSpec           ProductSpec              `json:"product_spec,omitempty"`
	ProductStocks         []ProductStock           `json:"product_stocks,omitempty"`
}

type ProductCategory struct {
	ID   uint64
	Name string
}

type ProductSpec struct {
	Spec1Title   vo.ProductSpecItem   `json:"spec_1_title"`
	Spec1Options []vo.ProductSpecItem `json:"spec_1_options"`
	Spec2Title   vo.ProductSpecItem   `json:"spec_2_title"`
	Spec2Options []vo.ProductSpecItem `json:"spec_2_options"`
}

type ProductStock struct {
	ID        uint64 `json:"id"`
	Spec1ID   uint64 `json:"spec_1_id"`
	Spec1Name string `json:"spec_1_name"`
	Spec2ID   uint64 `json:"spec_2_id"`
	Spec2Name string `json:"spec_2_name"`
	Quantity  int32  `json:"quantity"`
	Code      string `json:"code"`
	Spec      string `json:"spec,omitempty"`
}
