package application

import (
	"context"
	"ecommerce-api/internal/helper"
	"ecommerce-api/pkg/catalog/product/domain/entity"
	"ecommerce-api/pkg/catalog/product/domain/repository"
	"ecommerce-api/pkg/catalog/product/domain/vo"
	"fmt"
	"time"
)

type ProductService interface {
	ListProductCategory(ctx context.Context, merchantId uint64) ([]*vo.ProductCategoryItem, error)
	CreateProductCategory(ctx context.Context, params vo.CreateProductCategoryParams) (uint64, error)
	GetProductCategory(ctx context.Context, categoryId, merchantId uint64) (entity.ProductCategory, error)
	UpdateProductCategory(ctx context.Context, params vo.UpdateProductCategoryParams) error
	DeleteProductCategory(ctx context.Context, categoryId, merchantId uint64) error
	ListProduct(ctx context.Context, params vo.ListProductParams) ([]vo.ListProductRecord, int64, error)
	GetProduct(ctx context.Context, productId, merchantId uint64) (entity.Product, error)
	CreateProduct(ctx context.Context, params vo.CreateProductParams) (int64, error)
	UpdateProduct(ctx context.Context, params vo.UpdateProductParams) error
	DeleteProduct(ctx context.Context, productId, merchantId uint64) error
	SwitchProductStatus(ctx context.Context, productId, merchantId uint64) error
	ListProductByFilter(ctx context.Context, params vo.ListProductByFilterParams) ([]vo.ListProductRecord, int64, error)
	GetProductDetail(ctx context.Context, productId, merchantId uint64) (entity.Product, error)
	GetProductSpec(ctx context.Context, productId, merchantId uint64) (entity.ProductSpec, error)
	UpsertProductSpec(ctx context.Context, params vo.UpsertProductSpecParams) error
	DeleteProductSpec(ctx context.Context, specId, merchantId uint64) error
	GetProductStock(ctx context.Context, productId, merchantId uint64) ([]entity.ProductStock, error)
	UpsertProductStock(ctx context.Context, params vo.UpsertProductStockParams) error
	DeleteProductStock(ctx context.Context, stockId, merchantId uint64) error
}

type service struct {
	baseUrl string
	repo    repository.ProductRepository
}

func NewProductService(baseUrl string, repo repository.ProductRepository) ProductService {
	return &service{
		baseUrl: baseUrl,
		repo:    repo,
	}
}

/* Product Category */

func (s *service) ListProductCategory(ctx context.Context, merchantId uint64) ([]*vo.ProductCategoryItem, error) {
	categories, err := s.repo.ListProductCategory(ctx, merchantId)
	if err != nil {
		return nil, err
	}

	result := make([]*vo.ProductCategoryItem, 0)

	m := make(map[int64]*vo.ProductCategoryItem)
	for _, category := range categories {
		node := &vo.ProductCategoryItem{
			ID:       category.ID,
			Label:    category.Name,
			Children: make([]*vo.ProductCategoryItem, 0),
		}
		if _, ok := m[category.ID]; !ok {
			m[category.ID] = node
		}
		if category.ParentID == 0 {
			result = append(result, node)
		} else {
			tmp := m[category.ParentID]
			tmp.Children = append(tmp.Children, node)
		}
	}

	return result, nil
}

func (s *service) CreateProductCategory(ctx context.Context, params vo.CreateProductCategoryParams) (uint64, error) {
	if params.ParentID == 0 {
		// 檢查當前第 0 層上限 5 個
		count, err := s.repo.GetProductCategoryTopCount(ctx, params.MerchantID)
		if err != nil {
			return 0, err
		}

		if count == 5 {
			return 0, fmt.Errorf("商品分類上限為 5 個")
		}
	}

	return s.repo.CreateProductCategory(ctx, params)
}

func (s *service) GetProductCategory(ctx context.Context, categoryId, merchantId uint64) (entity.ProductCategory, error) {
	return s.repo.GetProductCategory(ctx, categoryId, merchantId)
}

func (s *service) UpdateProductCategory(ctx context.Context, params vo.UpdateProductCategoryParams) error {
	return s.repo.UpdateProductCategory(ctx, params)
}

func (s *service) DeleteProductCategory(ctx context.Context, categoryId, merchantId uint64) error {
	count, err := s.repo.GetProductCountByCategoryID(ctx, categoryId, merchantId)
	if err != nil {
		return err
	}
	if count != 0 {
		return fmt.Errorf("該分類底下尚有商品，請先轉移商品至其他分類")
	}

	return s.repo.DeleteProductCategory(ctx, categoryId, merchantId)
}

/* Product */

func (s *service) ListProduct(ctx context.Context, params vo.ListProductParams) ([]vo.ListProductRecord, int64, error) {
	products, total, err := s.repo.ListProduct(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	result := make([]vo.ListProductRecord, len(products))
	now := time.Now().UTC()

	for i := 0; i < len(products); i++ {
		specs, repoErr := s.repo.GetProductSpecTitlesByProductID(ctx, products[i].ID, params.MerchantID)
		if repoErr != nil {
			return nil, 0, err
		}

		result[i] = vo.ListProductRecord{
			ID:                 products[i].ID,
			Name:               products[i].Name,
			CurrencyID:         products[i].CurrencyID,
			Price:              products[i].Price,
			SpecialPrice:       0,
			Pictures:           products[i].Pictures,
			PictureUrls:        make([]string, 0),
			IsEnabled:          products[i].IsEnabled,
			TotalStockQuantity: products[i].TotalStockQuantity,
			SpecItems:          specs,
		}

		if products[i].SpecialPriceStart.Valid && products[i].SpecialPriceEnd.Valid {
			if now.After(products[i].SpecialPriceStart.Time) && now.Before(products[i].SpecialPriceEnd.Time) {
				result[i].SpecialPrice = products[i].SpecialPrice
			}
		}

		for _, picture := range products[i].Pictures {
			result[i].PictureUrls = append(result[i].PictureUrls, helper.GetUrl(s.baseUrl, picture))
		}
	}

	return result, total, nil
}

func (s *service) CreateProduct(ctx context.Context, params vo.CreateProductParams) (int64, error) {
	return s.repo.CreateProduct(ctx, params)
}

func (s *service) GetProduct(ctx context.Context, productId, merchantId uint64) (entity.Product, error) {
	product, err := s.repo.GetProduct(ctx, productId, merchantId)
	if err != nil {
		return entity.Product{}, err
	}

	for i := 0; i < len(product.Pictures); i++ {
		product.PictureUrls[i] = helper.GetUrl(s.baseUrl, product.Pictures[i])
	}

	return product, nil
}

func (s *service) UpdateProduct(ctx context.Context, params vo.UpdateProductParams) error {
	return s.repo.UpdateProduct(ctx, params)
}

func (s *service) DeleteProduct(ctx context.Context, productId, merchantId uint64) error {
	return s.repo.DeleteProduct(ctx, productId, merchantId)
}

func (s *service) SwitchProductStatus(ctx context.Context, productId, merchantId uint64) error {
	return s.repo.SwitchProductStatus(ctx, productId, merchantId)
}

func (s *service) ListProductByFilter(ctx context.Context, params vo.ListProductByFilterParams) ([]vo.ListProductRecord, int64, error) {
	products, total, err := s.repo.ListProductByFilter(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	result := make([]vo.ListProductRecord, len(products))
	now := time.Now().UTC()
	for i := 0; i < len(products); i++ {
		result[i] = vo.ListProductRecord{
			ID:                 products[i].ID,
			Name:               products[i].Name,
			CurrencyID:         products[i].CurrencyID,
			Price:              products[i].Price,
			SpecialPrice:       0,
			Pictures:           products[i].Pictures,
			PictureUrls:        make(vo.PictureArray, 0),
			IsEnabled:          products[i].IsEnabled,
			TotalStockQuantity: products[i].TotalStockQuantity,
		}

		if products[i].SpecialPriceStart.Valid && products[i].SpecialPriceEnd.Valid {
			if now.After(products[i].SpecialPriceStart.Time) && now.Before(products[i].SpecialPriceEnd.Time) {
				result[i].SpecialPrice = products[i].SpecialPrice
			}
		}

		for _, picture := range products[i].Pictures {
			result[i].PictureUrls = append(result[i].PictureUrls, helper.GetUrl(s.baseUrl, picture))
		}
	}

	return result, total, nil
}

func (s *service) GetProductDetail(ctx context.Context, productId, merchantId uint64) (entity.Product, error) {
	product, err := s.repo.GetProduct(ctx, productId, merchantId)
	if err != nil {
		return entity.Product{}, err
	}

	// 檢查優惠時間
	now := time.Now().UTC()
	if now.Before(product.SpecialPriceStart) || now.After(product.SpecialPriceEnd) {
		product.SpecialPrice = 0
	}

	for i := 0; i < len(product.Pictures); i++ {
		product.PictureUrls[i] = helper.GetUrl(s.baseUrl, product.Pictures[i])
	}

	return product, nil
}

/* Product Spec */

func (s *service) GetProductSpec(ctx context.Context, productId, merchantId uint64) (entity.ProductSpec, error) {
	return s.repo.GetProductSpec(ctx, productId, merchantId)
}

func (s *service) UpsertProductSpec(ctx context.Context, params vo.UpsertProductSpecParams) (err error) {
	return s.repo.UpsertProductSpec(ctx, params)
}

func (s *service) DeleteProductSpec(ctx context.Context, specId, merchantId uint64) error {
	return s.repo.DeleteProductSpec(ctx, specId, merchantId)
}

/* Product Stock */

func (s *service) GetProductStock(ctx context.Context, productId, merchantId uint64) ([]entity.ProductStock, error) {
	return s.repo.GetProductStock(ctx, productId, merchantId)
}

func (s *service) UpsertProductStock(ctx context.Context, params vo.UpsertProductStockParams) error {
	return s.repo.UpsertProductStock(ctx, params)
}

func (s *service) DeleteProductStock(ctx context.Context, stockId, merchantId uint64) error {
	return s.repo.DeleteProductStock(ctx, stockId, merchantId)
}
