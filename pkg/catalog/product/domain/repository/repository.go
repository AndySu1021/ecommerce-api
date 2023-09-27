package repository

import (
	"context"
	"ecommerce-api/pkg/catalog/product/domain/entity"
	"ecommerce-api/pkg/catalog/product/domain/vo"
)

type ProductRepository interface {
	ListProductCategory(ctx context.Context, merchantId uint64) ([]vo.ListProductCategoryRow, error)
	GetProductCategoryTopCount(ctx context.Context, merchantId uint64) (int64, error)
	CreateProductCategory(ctx context.Context, params vo.CreateProductCategoryParams) (uint64, error)
	GetProductCategory(ctx context.Context, categoryId, merchantId uint64) (entity.ProductCategory, error)
	UpdateProductCategory(ctx context.Context, params vo.UpdateProductCategoryParams) error
	GetProductCountByCategoryID(ctx context.Context, categoryId, merchantId uint64) (int64, error)
	DeleteProductCategory(ctx context.Context, categoryId, merchantId uint64) error
	ListProduct(ctx context.Context, params vo.ListProductParams) ([]vo.ListProductRow, int64, error)
	GetProductSpecTitlesByProductID(ctx context.Context, productId, merchantId uint64) ([]string, error)
	CreateProduct(ctx context.Context, params vo.CreateProductParams) (int64, error)
	GetProduct(ctx context.Context, productId, merchantId uint64) (entity.Product, error)
	UpdateProduct(ctx context.Context, params vo.UpdateProductParams) error
	DeleteProduct(ctx context.Context, productId, merchantId uint64) error
	SwitchProductStatus(ctx context.Context, productId, merchantId uint64) error
	ListProductByFilter(ctx context.Context, params vo.ListProductByFilterParams) ([]vo.ListProductRow, int64, error)
	GetProductSpec(ctx context.Context, productId, merchantId uint64) (entity.ProductSpec, error)
	UpsertProductSpec(ctx context.Context, params vo.UpsertProductSpecParams) error
	DeleteProductSpec(ctx context.Context, specId, merchantId uint64) error
	GetProductStock(ctx context.Context, productId, merchantId uint64) ([]entity.ProductStock, error)
	UpsertProductStock(ctx context.Context, params vo.UpsertProductStockParams) error
	DeleteProductStock(ctx context.Context, stockId, merchantId uint64) error
}
