package vo

import (
	"ecommerce-api/pkg/constant"
	"time"
)

/* Product Category */

type CreateProductCategoryParams struct {
	ParentID   uint64 `json:"parent_id" binding:"gte=0"`
	Name       string `json:"name" binding:"required"`
	MerchantID uint64
}

type UpdateProductCategoryParams struct {
	Name       string `json:"name" binding:"required"`
	CategoryID uint64
	MerchantID uint64
}

/* Product */

type ListProductParams struct {
	IsEnabled constant.YesNo `form:"is_enabled" binding:"min=-1,max=1"`
	constant.Pagination
	MerchantID uint64
}

type CreateProductParams struct {
	Name                  string                `json:"name" binding:"required"`
	Price                 uint64                `json:"price" binding:"required"`
	SpecialPrice          uint64                `json:"special_price" binding:"gte=0"`
	SpecialPriceStartStr  string                `json:"special_price_start" binding:""`
	SpecialPriceEndStr    string                `json:"special_price_end" binding:""`
	CategoryID            uint64                `json:"category_id" binding:"required"`
	Description           string                `json:"description" binding:"required"`
	Pictures              PictureArray          `json:"pictures" binding:"required"`
	SingleOrderLimit      int32                 `json:"single_order_limit" binding:"gte=0"`
	IsSingleOrderOnly     constant.YesNo        `json:"is_single_order_only" binding:"oneof=0 1"`
	Temperature           Temperature           `json:"temperature" binding:"min=1,max=3"`
	Length                int32                 `json:"length" binding:"gte=1"`
	Width                 int32                 `json:"width" binding:"gte=1"`
	Height                int32                 `json:"height" binding:"gte=1"`
	Weight                int32                 `json:"weight" binding:"required"`
	SupportDeliveryMethod SupportDeliveryMethod `json:"support_delivery_method" binding:"required"`
	IsAirContraband       constant.YesNo        `json:"is_air_contraband" binding:"oneof=0 1"`
	SpecialPriceStart     time.Time
	SpecialPriceEnd       time.Time
	MerchantID            uint64
}

type UpdateProductParams struct {
	Name                  string                `json:"name" binding:"required"`
	Price                 uint64                `json:"price" binding:"required"`
	SpecialPrice          uint64                `json:"special_price" binding:"gte=0"`
	SpecialPriceStartStr  string                `json:"special_price_start" binding:""`
	SpecialPriceEndStr    string                `json:"special_price_end" binding:""`
	CategoryID            uint64                `json:"category_id" binding:"required"`
	Description           string                `json:"description" binding:"required"`
	Pictures              PictureArray          `json:"pictures" binding:"required"`
	SingleOrderLimit      int32                 `json:"single_order_limit" binding:"gte=0"`
	IsSingleOrderOnly     constant.YesNo        `json:"is_single_order_only" binding:"oneof=0 1"`
	Temperature           Temperature           `json:"temperature" binding:"min=1,max=3"`
	Length                int32                 `json:"length" binding:"gte=1"`
	Width                 int32                 `json:"width" binding:"gte=1"`
	Height                int32                 `json:"height" binding:"gte=1"`
	Weight                int32                 `json:"weight" binding:"required"`
	SupportDeliveryMethod SupportDeliveryMethod `json:"support_delivery_method" binding:"required"`
	IsAirContraband       constant.YesNo        `json:"is_air_contraband" binding:"oneof=0 1"`
	SpecialPriceStart     time.Time
	SpecialPriceEnd       time.Time
	ProductID             uint64
	MerchantID            uint64
}

/* Product Spec */

type OptionalProductSpecItem struct {
	ID   uint64 `json:"id" binding:"gte=0"`
	Name string `json:"name" binding:""`
}

type UpsertProductSpecParams struct {
	Spec1Title   ProductSpecItem           `json:"spec_1_title" binding:"required"`
	Spec1Options []ProductSpecItem         `json:"spec_1_options" binding:"required"`
	Spec2Title   OptionalProductSpecItem   `json:"spec_2_title" binding:""`
	Spec2Options []OptionalProductSpecItem `json:"spec_2_options" binding:""`
	ProductID    uint64
	MerchantID   uint64
}

/* Product Stock */

type CreateProductStockParams struct {
	ProductID int64 `json:"product_id" binding:"required"`
	Spec1ID   int64 `json:"spec_1_id" binding:"required"`
	Spec2ID   int64 `json:"spec_2_id" binding:"gte=0"`
	Quantity  int32 `json:"quantity" binding:"gte=0"`
}

type ProductStockItem struct {
	Spec1ID  int64  `json:"spec_1_id" binding:"required"`
	Spec2ID  int64  `json:"spec_2_id" binding:"gte=0"`
	Quantity int32  `json:"quantity" binding:"gte=0"`
	Code     string `json:"code" binding:"required"`
}

type UpsertProductStockParams struct {
	Stocks     []ProductStockItem `json:"stocks" binding:"required"`
	ProductID  uint64
	MerchantID uint64
}

type ListProductByFilterParams struct {
	Keyword    string `form:"keyword" binding:""`
	CategoryID uint64 `form:"category_id" binding:""`
	PriceRange string `form:"price_range"  binding:""`
	Sort       SortBy `form:"sort"  binding:""`
	Direction  string `form:"direction"  binding:"oneof=asc desc"`
	constant.Pagination
	MerchantID uint64
}
