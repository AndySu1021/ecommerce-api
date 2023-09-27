package persistence

import (
	"context"
	"database/sql"
	"ecommerce-api/internal/db/model"
	"ecommerce-api/internal/helper"
	"ecommerce-api/pkg/catalog/product/domain/entity"
	"ecommerce-api/pkg/catalog/product/domain/vo"
	"ecommerce-api/pkg/constant"
	"fmt"
	"github.com/Masterminds/squirrel"
	"strings"
	"time"
)

type ProductRepository struct {
	db      *sql.DB
	queries *model.Queries
}

func NewProductRepository(db *sql.DB, queries *model.Queries) *ProductRepository {
	return &ProductRepository{
		db:      db,
		queries: queries,
	}
}

func (r *ProductRepository) ListProductCategory(ctx context.Context, merchantId uint64) (categories []vo.ListProductCategoryRow, err error) {
	columns := []string{"id", "name", "parent_id"}

	builder := squirrel.Select(columns...).From("product_category").
		Where(squirrel.Eq{"merchant_id": merchantId}).
		OrderBy("parent_id").
		RunWith(r.db)

	if err = helper.StructQuery(ctx, builder, &categories); err != nil {
		return nil, err
	}

	return
}

func (r *ProductRepository) GetProductCategoryTopCount(ctx context.Context, merchantId uint64) (int64, error) {
	return r.queries.GetProductCategoryTopCount(ctx, merchantId)
}

func (r *ProductRepository) CreateProductCategory(ctx context.Context, params vo.CreateProductCategoryParams) (uint64, error) {
	var (
		categoryId uint64
		err        error
		now        = time.Now().UTC()
	)

	modelParams := model.CreateProductCategoryParams{
		MerchantID: params.MerchantID,
		Name:       params.Name,
		TopID:      0,
		ParentID:   params.ParentID,
		TreeLeft:   1,
		TreeRight:  2,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err = helper.Transaction(ctx, r.db, func(ctx context.Context, tx *sql.Tx) error {
		if modelParams.ParentID != 0 {
			category, errRepo := r.queries.WithTx(tx).GetProductCategory(ctx, model.GetProductCategoryParams{
				ID:         params.ParentID,
				MerchantID: params.MerchantID,
			})
			if errRepo != nil {
				return errRepo
			}

			modelParams.TreeLeft = category.TreeRight
			modelParams.TreeRight = category.TreeRight + 1
			modelParams.TopID = category.TopID

			if errRepo = r.queries.WithTx(tx).UpdateProductCategoryLeftTree(ctx, model.UpdateProductCategoryLeftTreeParams{
				UpdatedAt: now,
				TreeLeft:  category.TreeRight,
				TreeRight: category.TreeRight,
				TopID:     category.TopID,
			}); errRepo != nil {
				return errRepo
			}

			if errRepo = r.queries.WithTx(tx).UpdateProductCategoryRightTree(ctx, model.UpdateProductCategoryRightTreeParams{
				UpdatedAt: now,
				TreeLeft:  category.TreeRight,
				TopID:     category.TopID,
			}); errRepo != nil {
				return errRepo
			}
		}

		result, errRepo2 := r.queries.WithTx(tx).CreateProductCategory(ctx, modelParams)
		if errRepo2 != nil {
			return errRepo2
		}

		tmpId, err := result.LastInsertId()
		if err != nil {
			return err
		}

		categoryId = uint64(tmpId)

		// 有上級就可以先回傳了
		if modelParams.ParentID != 0 {
			return nil
		}

		return r.queries.WithTx(tx).UpdateProductCategoryTopID(ctx, model.UpdateProductCategoryTopIDParams{
			TopID:      categoryId,
			ID:         categoryId,
			MerchantID: params.MerchantID,
		})
	}); err != nil {
		return 0, err
	}

	return categoryId, err
}

func (r *ProductRepository) GetProductCategory(ctx context.Context, categoryId, merchantId uint64) (entity.ProductCategory, error) {
	category, err := r.queries.GetProductCategory(ctx, model.GetProductCategoryParams{
		ID:         categoryId,
		MerchantID: merchantId,
	})
	if err != nil {
		return entity.ProductCategory{}, err
	}

	return entity.ProductCategory{
		ID:   category.ID,
		Name: category.Name,
	}, nil
}

func (r *ProductRepository) UpdateProductCategory(ctx context.Context, params vo.UpdateProductCategoryParams) error {
	return r.queries.UpdateProductCategory(ctx, model.UpdateProductCategoryParams{
		Name:       params.Name,
		UpdatedAt:  time.Now().UTC(),
		ID:         params.CategoryID,
		MerchantID: params.MerchantID,
	})
}

func (r *ProductRepository) GetProductCountByCategoryID(ctx context.Context, categoryId, merchantId uint64) (int64, error) {
	return r.queries.GetProductCountByCategoryID(ctx, model.GetProductCountByCategoryIDParams{
		ID:         categoryId,
		MerchantID: merchantId,
	})
}

func (r *ProductRepository) DeleteProductCategory(ctx context.Context, categoryId, merchantId uint64) error {
	category, err := r.queries.GetProductCategory(ctx, model.GetProductCategoryParams{
		ID:         categoryId,
		MerchantID: merchantId,
	})
	if err != nil {
		return err
	}

	return r.queries.DeleteProductCategory(ctx, model.DeleteProductCategoryParams{
		TopID:      category.TopID,
		TreeLeft:   category.TreeLeft,
		TreeRight:  category.TreeRight,
		MerchantID: category.MerchantID,
	})
}

func (r *ProductRepository) ListProduct(ctx context.Context, params vo.ListProductParams) (products []vo.ListProductRow, total int64, err error) {
	columns := []string{
		"p.id",
		"p.name",
		"p.currency_id",
		"p.price",
		"p.special_price",
		"p.special_price_start",
		"p.special_price_end",
		"p.pictures",
		"p.is_enabled",
		"IFNULL(tmp.quantity, 0) AS total_stock_quantity",
	}

	dQuery := squirrel.Select(columns...).
		From("product p").
		Where(squirrel.Eq{"p.merchant_id": params.MerchantID}).
		LeftJoin("(select product_id, SUM(quantity) as quantity from product_stock where merchant_id = ? group by product_id) tmp ON tmp.product_id = p.id", params.MerchantID)
	cQuery := squirrel.Select("count(*) AS count").
		From("product p").
		Where(squirrel.Eq{"p.merchant_id": params.MerchantID})

	if params.IsEnabled != constant.All {
		dQuery = dQuery.Where(squirrel.Eq{"p.is_enabled": params.IsEnabled})
		cQuery = cQuery.Where(squirrel.Eq{"p.is_enabled": params.IsEnabled})
	}

	dQuery = dQuery.RunWith(r.db)
	cQuery = cQuery.RunWith(r.db)

	if err = helper.PageQuery(ctx, dQuery, params.Pagination, &products); err != nil {
		return
	}
	if err = helper.TotalQuery(ctx, cQuery, &total); err != nil {
		return
	}

	return
}

func (r *ProductRepository) GetProductSpecTitlesByProductID(ctx context.Context, productId, merchantId uint64) ([]string, error) {
	return r.queries.GetProductSpecTitlesByProductID(ctx, model.GetProductSpecTitlesByProductIDParams{
		ProductID:  productId,
		MerchantID: merchantId,
	})
}

func (r *ProductRepository) CreateProduct(ctx context.Context, params vo.CreateProductParams) (int64, error) {
	now := time.Now().UTC()
	modelParams := model.CreateProductParams{
		Name:                  params.Name,
		CategoryID:            params.CategoryID,
		CurrencyID:            constant.CurrencyTWD,
		Price:                 params.Price,
		SpecialPrice:          params.SpecialPrice,
		SpecialPriceStart:     sql.NullTime{},
		SpecialPriceEnd:       sql.NullTime{},
		SingleOrderLimit:      params.SingleOrderLimit,
		IsSingleOrderOnly:     params.IsSingleOrderOnly,
		Temperature:           params.Temperature,
		Length:                params.Length,
		Width:                 params.Width,
		Height:                params.Height,
		Weight:                params.Weight,
		SupportDeliveryMethod: params.SupportDeliveryMethod,
		IsAirContraband:       params.IsAirContraband,
		Description:           params.Description,
		Pictures:              params.Pictures,
		Extra:                 vo.Extra{},
		CreatedAt:             now,
		UpdatedAt:             now,
	}

	if params.SpecialPrice != 0 {
		modelParams.SpecialPriceStart = sql.NullTime{
			Time:  params.SpecialPriceStart,
			Valid: true,
		}
		modelParams.SpecialPriceEnd = sql.NullTime{
			Time:  params.SpecialPriceEnd,
			Valid: true,
		}
	}

	result, err := r.queries.CreateProduct(ctx, modelParams)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *ProductRepository) GetProduct(ctx context.Context, productId, merchantId uint64) (entity.Product, error) {
	product, err := r.queries.GetProduct(ctx, model.GetProductParams{
		ID:         productId,
		MerchantID: merchantId,
	})
	if err != nil {
		return entity.Product{}, err
	}

	specs, err := r.queries.GetProductSpec(ctx, model.GetProductSpecParams{
		ProductID:  productId,
		MerchantID: merchantId,
	})
	if err != nil {
		return entity.Product{}, err
	}

	stocks, err := r.queries.GetProductStock(ctx, model.GetProductStockParams{
		ProductID:  productId,
		MerchantID: merchantId,
	})
	if err != nil {
		return entity.Product{}, err
	}

	result := entity.Product{
		ID:                    productId,
		Name:                  product.Name,
		CurrencyID:            product.CurrencyID,
		Price:                 product.Price,
		SpecialPrice:          product.SpecialPrice,
		SpecialPriceStart:     product.SpecialPriceStart.Time,
		SpecialPriceEnd:       product.SpecialPriceEnd.Time,
		CategoryID:            product.CategoryID,
		CategoryName:          product.CategoryName,
		Description:           product.Description,
		Pictures:              product.Pictures,
		PictureUrls:           make(vo.PictureArray, 0),
		SingleOrderLimit:      product.SingleOrderLimit,
		IsSingleOrderOnly:     product.IsSingleOrderOnly,
		Temperature:           product.Temperature,
		Length:                product.Length,
		Width:                 product.Width,
		Height:                product.Height,
		Weight:                product.Weight,
		SupportDeliveryMethod: product.SupportDeliveryMethod,
		IsAirContraband:       product.IsAirContraband,
		Extra:                 product.Extra,
		IsEnabled:             product.IsEnabled,
		TotalStockQuantity:    0,
		ProductSpec:           r.formatSpec(specs),
		ProductStocks:         make([]entity.ProductStock, 0),
	}

	for _, stock := range stocks {
		tmp := entity.ProductStock{
			ID:        stock.ID,
			Spec1ID:   stock.Spec1ID,
			Spec1Name: stock.Spec1Name,
			Spec2ID:   stock.Spec2ID,
			Spec2Name: string(stock.Spec2Name.([]uint8)),
			Quantity:  stock.Quantity,
			Code:      stock.Code,
			Spec:      stock.Spec1Name,
		}

		if tmp.Spec2Name != "" {
			tmp.Spec = fmt.Sprintf("%s / %s", stock.Spec1Name, stock.Spec2Name)
		}

		result.ProductStocks = append(result.ProductStocks, tmp)
	}

	return result, nil
}

func (r *ProductRepository) UpdateProduct(ctx context.Context, params vo.UpdateProductParams) error {
	modelParams := model.UpdateProductParams{
		Name:                  params.Name,
		CategoryID:            params.CategoryID,
		Description:           params.Description,
		Price:                 params.Price,
		SpecialPrice:          params.SpecialPrice,
		SpecialPriceStart:     sql.NullTime{},
		SpecialPriceEnd:       sql.NullTime{},
		SingleOrderLimit:      params.SingleOrderLimit,
		IsSingleOrderOnly:     params.IsSingleOrderOnly,
		Temperature:           params.Temperature,
		Length:                params.Length,
		Width:                 params.Width,
		Height:                params.Height,
		Weight:                params.Weight,
		SupportDeliveryMethod: params.SupportDeliveryMethod,
		IsAirContraband:       params.IsAirContraband,
		Pictures:              params.Pictures,
		UpdatedAt:             time.Now().UTC(),
		ID:                    params.ProductID,
		MerchantID:            params.MerchantID,
	}

	if params.SpecialPrice != 0 {
		modelParams.SpecialPriceStart = sql.NullTime{
			Time:  params.SpecialPriceStart,
			Valid: true,
		}
		modelParams.SpecialPriceEnd = sql.NullTime{
			Time:  params.SpecialPriceEnd,
			Valid: true,
		}
	}

	return r.queries.UpdateProduct(ctx, modelParams)
}

func (r *ProductRepository) DeleteProduct(ctx context.Context, productId, merchantId uint64) error {
	return helper.Transaction(ctx, r.db, func(ctx context.Context, tx *sql.Tx) error {
		if err := r.queries.WithTx(tx).DeleteProduct(ctx, model.DeleteProductParams{
			ID:         productId,
			MerchantID: merchantId,
		}); err != nil {
			return err
		}
		if err := r.queries.WithTx(tx).DeleteProductSpecByProductID(ctx, model.DeleteProductSpecByProductIDParams{
			ProductID:  productId,
			MerchantID: merchantId,
		}); err != nil {
			return err
		}
		if err := r.queries.WithTx(tx).DeleteProductStockByProductID(ctx, model.DeleteProductStockByProductIDParams{
			ProductID:  productId,
			MerchantID: merchantId,
		}); err != nil {
			return err
		}
		return nil
	})
}

func (r *ProductRepository) SwitchProductStatus(ctx context.Context, productId, merchantId uint64) error {
	return r.queries.SwitchProductStatus(ctx, model.SwitchProductStatusParams{
		ID:         productId,
		MerchantID: merchantId,
	})
}

func (r *ProductRepository) ListProductByFilter(ctx context.Context, params vo.ListProductByFilterParams) (products []vo.ListProductRow, total int64, err error) {
	categoryIds, err := r.queries.GetProductCategoryChildrenIDs(ctx, model.GetProductCategoryChildrenIDsParams{
		ID:         params.CategoryID,
		MerchantID: params.MerchantID,
	})
	if err != nil {
		return nil, 0, err
	}

	columns := []string{
		"p.id",
		"p.name",
		"p.currency_id",
		"p.price",
		"p.special_price",
		"p.special_price_start",
		"p.special_price_end",
		"p.pictures",
		"p.is_enabled",
		"IFNULL(tmp.quantity, 0) AS total_stock_quantity",
	}

	dQuery := squirrel.Select(columns...).From("product p").
		Where(squirrel.Eq{"p.is_enabled": constant.Yes}).
		Where(squirrel.Eq{"p.merchant_id": params.MerchantID})
	cQuery := squirrel.Select("count(*) AS count").From("product p").
		Where(squirrel.Eq{"p.is_enabled": constant.Yes}).
		Where(squirrel.Eq{"p.merchant_id": params.MerchantID})

	dQuery = dQuery.LeftJoin("(select product_id, SUM(quantity) as quantity from product_stock where merchant_id = ? group by product_id) tmp ON tmp.product_id = p.id", params.MerchantID)

	if params.Keyword != "" {
		dQuery = dQuery.Where(squirrel.Or{squirrel.Like{"p.name": "%" + params.Keyword + "%"}, squirrel.Like{"p.description": "%" + params.Keyword + "%"}})
		cQuery = cQuery.Where(squirrel.Or{squirrel.Like{"p.name": "%" + params.Keyword + "%"}, squirrel.Like{"p.description": "%" + params.Keyword + "%"}})
	}

	if len(categoryIds) > 0 {
		dQuery = dQuery.Where(squirrel.Eq{"p.category_id": categoryIds})
		cQuery = cQuery.Where(squirrel.Eq{"p.category_id": categoryIds})
	}

	if params.PriceRange != "" {
		now := time.Now().UTC()
		tmp := strings.Split(params.PriceRange, "-")
		dQuery = dQuery.Where(squirrel.Expr("((special_price >= ? AND special_price <= ? AND special_price_start < ? AND special_price_end > ?) OR (price >= ? AND price <= ?))", tmp[0], tmp[1], now, now, tmp[0], tmp[1]))
		cQuery = cQuery.Where(squirrel.Expr("((special_price >= ? AND special_price <= ? AND special_price_start < ? AND special_price_end > ?) OR (price >= ? AND price <= ?))", tmp[0], tmp[1], now, now, tmp[0], tmp[1]))
	}

	switch params.Sort {
	case vo.SortByDefault:
		dQuery = dQuery.OrderBy("p.updated_at desc")
	case vo.SortBySales:
		dQuery = dQuery.OrderBy("p.sales " + params.Direction)
	case vo.SortByStock:
		dQuery = dQuery.Where(squirrel.Gt{"p.tmp.quantity": 0})
	case vo.SortByPrice:
		dQuery = dQuery.OrderBy("p.price " + params.Direction)
	case vo.SortByTime:
		dQuery = dQuery.OrderBy("p.created_at " + params.Direction)
	}

	dQuery = dQuery.RunWith(r.db)
	cQuery = cQuery.RunWith(r.db)

	if err = helper.PageQuery(ctx, dQuery, params.Pagination, &products); err != nil {
		return
	}
	if err = helper.TotalQuery(ctx, cQuery, &total); err != nil {
		return
	}

	return
}

func (r *ProductRepository) GetProductSpec(ctx context.Context, productId, merchantId uint64) (entity.ProductSpec, error) {
	specs, err := r.queries.GetProductSpec(ctx, model.GetProductSpecParams{
		ProductID:  productId,
		MerchantID: merchantId,
	})
	if err != nil {
		return entity.ProductSpec{}, err
	}

	return r.formatSpec(specs), nil
}

func (r *ProductRepository) UpsertProductSpec(ctx context.Context, params vo.UpsertProductSpecParams) error {
	now := time.Now().UTC()
	return helper.Transaction(ctx, r.db, func(ctx context.Context, tx *sql.Tx) error {
		// exec spec 1
		if params.Spec1Title.ID != 0 {
			if err := r.queries.WithTx(tx).UpdateProductSpec(ctx, model.UpdateProductSpecParams{
				Name:       params.Spec1Title.Name,
				UpdatedAt:  now,
				ID:         params.Spec1Title.ID,
				MerchantID: params.MerchantID,
			}); err != nil {
				return err
			}
		} else {
			if err := r.queries.WithTx(tx).CreateProductSpec(ctx, model.CreateProductSpecParams{
				MerchantID: params.MerchantID,
				ProductID:  params.ProductID,
				Level:      1,
				Type:       1,
				Name:       params.Spec1Title.Name,
				CreatedAt:  now,
				UpdatedAt:  now,
			}); err != nil {
				return err
			}
		}

		for _, option := range params.Spec1Options {
			if option.ID != 0 {
				if err := r.queries.WithTx(tx).UpdateProductSpec(ctx, model.UpdateProductSpecParams{
					Name:       option.Name,
					UpdatedAt:  now,
					ID:         option.ID,
					MerchantID: params.MerchantID,
				}); err != nil {
					return err
				}
			} else {
				if err := r.queries.WithTx(tx).CreateProductSpec(ctx, model.CreateProductSpecParams{
					MerchantID: params.MerchantID,
					ProductID:  params.ProductID,
					Level:      1,
					Type:       2,
					Name:       option.Name,
					CreatedAt:  now,
					UpdatedAt:  now,
				}); err != nil {
					return err
				}
			}
		}

		// exec spec 2
		if params.Spec2Title.Name == "" {
			return nil
		}

		if params.Spec2Title.ID != 0 {
			if err := r.queries.WithTx(tx).UpdateProductSpec(ctx, model.UpdateProductSpecParams{
				Name:       params.Spec2Title.Name,
				UpdatedAt:  now,
				ID:         params.Spec2Title.ID,
				MerchantID: params.MerchantID,
			}); err != nil {
				return err
			}
		} else {
			if err := r.queries.WithTx(tx).CreateProductSpec(ctx, model.CreateProductSpecParams{
				MerchantID: params.MerchantID,
				ProductID:  params.ProductID,
				Level:      2,
				Type:       1,
				Name:       params.Spec2Title.Name,
				CreatedAt:  now,
				UpdatedAt:  now,
			}); err != nil {
				return err
			}
		}

		for _, option := range params.Spec2Options {
			if option.ID != 0 {
				if err := r.queries.WithTx(tx).UpdateProductSpec(ctx, model.UpdateProductSpecParams{
					Name:       option.Name,
					UpdatedAt:  now,
					ID:         option.ID,
					MerchantID: params.MerchantID,
				}); err != nil {
					return err
				}
			} else {
				if err := r.queries.WithTx(tx).CreateProductSpec(ctx, model.CreateProductSpecParams{
					MerchantID: params.MerchantID,
					ProductID:  params.ProductID,
					Level:      2,
					Type:       2,
					Name:       option.Name,
					CreatedAt:  now,
					UpdatedAt:  now,
				}); err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func (r *ProductRepository) DeleteProductSpec(ctx context.Context, specId, merchantId uint64) error {
	return helper.Transaction(ctx, r.db, func(ctx context.Context, tx *sql.Tx) error {
		spec, err := r.queries.WithTx(tx).GetProductSpecByID(ctx, model.GetProductSpecByIDParams{
			ID:         specId,
			MerchantID: merchantId,
		})
		if err != nil {
			return err
		}

		if spec.Level == 1 && spec.Type == 1 {
			return fmt.Errorf("刪除失敗，請重新操作")
		}

		specIds := []uint64{specId}
		if spec.Level == 2 && spec.Type == 1 {
			ids, tmpErr := r.queries.WithTx(tx).ListProductSecondLevelSpecIDs(ctx, model.ListProductSecondLevelSpecIDsParams{
				ProductID:  spec.ProductID,
				MerchantID: merchantId,
			})
			if tmpErr != nil {
				return tmpErr
			}

			specIds = ids

			if err = r.queries.WithTx(tx).DeleteProductSecondLevelSpec(ctx, model.DeleteProductSecondLevelSpecParams{
				ProductID:  spec.ProductID,
				MerchantID: merchantId,
			}); err != nil {
				return err
			}
		}

		stockIds := make([]int64, 0)
		rows, err := squirrel.Select("id").From("product_stock").
			Where(squirrel.Or{squirrel.Eq{"spec_1_id": specIds}, squirrel.Eq{"spec_2_id": specIds}}).
			RunWith(tx).
			QueryContext(ctx)
		if err != nil {
			return err
		}

		for rows.Next() {
			var i int64
			if err = rows.Scan(&i); err != nil {
				return err
			}
			stockIds = append(stockIds, i)
		}
		if err = rows.Close(); err != nil {
			return err
		}
		if err = rows.Err(); err != nil {
			return err
		}

		if _, err = squirrel.Delete("shopping_cart").
			Where(squirrel.Eq{"stock_id": stockIds}).
			RunWith(tx).
			ExecContext(ctx); err != nil {
			return err
		}

		if _, err = squirrel.Delete("product_stock").
			Where(squirrel.Or{squirrel.Eq{"spec_1_id": specIds}, squirrel.Eq{"spec_2_id": specIds}}).
			RunWith(tx).
			ExecContext(ctx); err != nil {
			return err
		}

		return r.queries.WithTx(tx).DeleteProductSpec(ctx, model.DeleteProductSpecParams{
			ID:         specId,
			MerchantID: merchantId,
		})
	})
}

func (r *ProductRepository) formatSpec(specs []model.GetProductSpecRow) entity.ProductSpec {
	productSpec := entity.ProductSpec{
		Spec1Title:   vo.ProductSpecItem{},
		Spec1Options: make([]vo.ProductSpecItem, 0),
		Spec2Title:   vo.ProductSpecItem{},
		Spec2Options: make([]vo.ProductSpecItem, 0),
	}
	for _, spec := range specs {
		item := vo.ProductSpecItem{
			ID:   spec.ID,
			Name: spec.Name,
		}
		if spec.Level == 1 {
			if spec.Type == 1 {
				productSpec.Spec1Title = item
			} else {
				productSpec.Spec1Options = append(productSpec.Spec1Options, item)
			}
		} else if spec.Level == 2 {
			if spec.Type == 1 {
				productSpec.Spec2Title = item
			} else {
				productSpec.Spec2Options = append(productSpec.Spec2Options, item)
			}
		}
	}
	return productSpec
}

func (r *ProductRepository) GetProductStock(ctx context.Context, productId, merchantId uint64) ([]entity.ProductStock, error) {
	stocks, err := r.queries.GetProductStock(ctx, model.GetProductStockParams{
		ProductID:  productId,
		MerchantID: merchantId,
	})
	if err != nil {
		return nil, err
	}

	result := make([]entity.ProductStock, len(stocks))

	for i := 0; i < len(stocks); i++ {
		result[i] = entity.ProductStock{
			ID:        stocks[i].ID,
			Spec1ID:   stocks[i].Spec1ID,
			Spec1Name: stocks[i].Spec1Name,
			Spec2ID:   stocks[i].Spec2ID,
			Spec2Name: string(stocks[i].Spec2Name.([]uint8)),
			Quantity:  stocks[i].Quantity,
			Code:      stocks[i].Code,
			Spec:      "",
		}
	}

	return result, nil
}

func (r *ProductRepository) UpsertProductStock(ctx context.Context, params vo.UpsertProductStockParams) error {
	if len(params.Stocks) == 0 {
		return nil
	}

	builder := squirrel.Insert("product_stock").
		Columns("product_id", "spec_1_id", "spec_2_id", "quantity", "code", "created_at", "updated_at")

	now := time.Now().UTC()
	for _, row := range params.Stocks {
		builder = builder.Values(
			params.ProductID,
			row.Spec1ID,
			row.Spec2ID,
			row.Quantity,
			row.Code,
			now,
			now,
		)
	}

	tmp := []string{
		"quantity = VALUES (quantity)",
		"code = VALUES (code)",
		"updated_at = VALUES (updated_at)",
	}
	if _, err := builder.SuffixExpr(squirrel.Expr("ON DUPLICATE KEY UPDATE " + strings.Join(tmp, ","))).
		RunWith(r.db).
		ExecContext(ctx); err != nil {
		return err
	}

	return nil
}

func (r *ProductRepository) DeleteProductStock(ctx context.Context, stockId, merchantId uint64) error {
	return r.queries.DeleteProductStock(ctx, model.DeleteProductStockParams{
		ID:         stockId,
		MerchantID: merchantId,
	})
}
