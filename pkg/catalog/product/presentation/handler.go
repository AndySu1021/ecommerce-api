package presentation

import (
	"ecommerce-api/internal/errors"
	"ecommerce-api/internal/helper"
	"ecommerce-api/pkg/catalog/product/application"
	"ecommerce-api/pkg/catalog/product/domain/vo"
	app_admin "ecommerce-api/pkg/identity/admin/application"
	app_merchant "ecommerce-api/pkg/identity/merchant/application"
	ginTool "github.com/AndySu1021/go-util/gin"
	"github.com/gin-gonic/gin"
	"strconv"
)

type ProductHandler struct {
	adminSvc    app_admin.AdminService
	merchantSvc app_merchant.MerchantService
	productSvc  application.ProductService
}

func NewProductHandler(adminSvc app_admin.AdminService, merchantSvc app_merchant.MerchantService, productSvc application.ProductService) *ProductHandler {
	return &ProductHandler{
		adminSvc:    adminSvc,
		merchantSvc: merchantSvc,
		productSvc:  productSvc,
	}
}

/* Product Category */

func (h *ProductHandler) ListProductCategory(c *gin.Context) {
	var (
		err error
		ctx = c.Request.Context()
	)

	admin, err := h.adminSvc.GetCurrentAdmin(ctx)
	if err != nil {
		helper.ErrorResp(c, err)
		return
	}

	categories, err := h.productSvc.ListProductCategory(ctx, admin.Merchant.ID)
	if err != nil {
		helper.ErrorResp(c, err)
		return
	}

	ginTool.SuccessWithData(c, categories)
}

func (h *ProductHandler) CreateProductCategory(c *gin.Context) {
	var (
		err    error
		params vo.CreateProductCategoryParams
		ctx    = c.Request.Context()
	)

	if err = c.ShouldBindJSON(&params); err != nil {
		helper.ErrorResp(c, errors.NewErrValidation(err))
		return
	}

	admin, err := h.adminSvc.GetCurrentAdmin(ctx)
	if err != nil {
		helper.ErrorResp(c, err)
		return
	}

	params.MerchantID = admin.Merchant.ID

	categoryId, err := h.productSvc.CreateProductCategory(ctx, params)
	if err != nil {
		helper.ErrorResp(c, err)
		return
	}

	ginTool.SuccessWithData(c, categoryId)
}

func (h *ProductHandler) GetProductCategory(c *gin.Context) {
	var (
		err error
		ctx = c.Request.Context()
	)

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.ErrorResp(c, err)
		return
	}

	admin, err := h.adminSvc.GetCurrentAdmin(ctx)
	if err != nil {
		helper.ErrorResp(c, err)
		return
	}

	category, err := h.productSvc.GetProductCategory(ctx, id, admin.Merchant.ID)
	if err != nil {
		helper.ErrorResp(c, err)
		return
	}

	ginTool.SuccessWithData(c, category)
}

func (h *ProductHandler) UpdateProductCategory(c *gin.Context) {
	var (
		err    error
		params vo.UpdateProductCategoryParams
		ctx    = c.Request.Context()
	)

	if err = c.ShouldBindJSON(&params); err != nil {
		helper.ErrorResp(c, errors.NewErrValidation(err))
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.ErrorResp(c, err)
		return
	}

	admin, err := h.adminSvc.GetCurrentAdmin(ctx)
	if err != nil {
		helper.ErrorResp(c, err)
		return
	}

	params.CategoryID = id
	params.MerchantID = admin.Merchant.ID

	if err = h.productSvc.UpdateProductCategory(ctx, params); err != nil {
		helper.ErrorResp(c, err)
		return
	}

	ginTool.Success(c)
}

func (h *ProductHandler) DeleteProductCategory(c *gin.Context) {
	var (
		err error
		ctx = c.Request.Context()
	)

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.ErrorResp(c, err)
		return
	}

	admin, err := h.adminSvc.GetCurrentAdmin(ctx)
	if err != nil {
		helper.ErrorResp(c, err)
		return
	}

	if err = h.productSvc.DeleteProductCategory(ctx, id, admin.Merchant.ID); err != nil {
		helper.ErrorResp(c, err)
		return
	}

	ginTool.Success(c)
}

/* Product */

func (h *ProductHandler) ListProduct(c *gin.Context) {
	var (
		err    error
		params vo.ListProductParams
		ctx    = c.Request.Context()
	)

	if err = c.ShouldBindQuery(&params); err != nil {
		helper.ErrorResp(c, errors.NewErrValidation(err))
		return
	}

	admin, err := h.adminSvc.GetCurrentAdmin(ctx)
	if err != nil {
		helper.ErrorResp(c, err)
		return
	}

	params.MerchantID = admin.Merchant.ID

	products, total, err := h.productSvc.ListProduct(ctx, params)
	if err != nil {
		helper.ErrorResp(c, err)
		return
	}

	params.Pagination.Total = total

	ginTool.SuccessWithPagination(c, products, params.Pagination)
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var (
		err    error
		params vo.CreateProductParams
		ctx    = c.Request.Context()
	)

	if err = c.ShouldBindJSON(&params); err != nil {
		helper.ErrorResp(c, errors.NewErrValidation(err))
		return
	}

	if params.SpecialPrice != 0 {
		strArr := []string{params.SpecialPriceStartStr, params.SpecialPriceEndStr}
		timeArr, parseErr := helper.ParseUTC8Datetime(strArr)
		if parseErr != nil {
			helper.ErrorResp(c, errors.ErrWrongTimeFormat)
			return
		}
		params.SpecialPriceStart = timeArr[0]
		params.SpecialPriceEnd = timeArr[1]
	}

	admin, err := h.adminSvc.GetCurrentAdmin(ctx)
	if err != nil {
		helper.ErrorResp(c, err)
		return
	}

	params.MerchantID = admin.Merchant.ID

	id, err := h.productSvc.CreateProduct(ctx, params)
	if err != nil {
		helper.ErrorResp(c, err)
		return
	}

	ginTool.SuccessWithData(c, id)
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
	var (
		err error
		ctx = c.Request.Context()
	)

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.ErrorResp(c, err)
		return
	}

	admin, err := h.adminSvc.GetCurrentAdmin(ctx)
	if err != nil {
		helper.ErrorResp(c, err)
		return
	}

	product, err := h.productSvc.GetProduct(ctx, id, admin.Merchant.ID)
	if err != nil {
		helper.ErrorResp(c, err)
		return
	}

	ginTool.SuccessWithData(c, product)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	var (
		err    error
		params vo.UpdateProductParams
		ctx    = c.Request.Context()
	)

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.ErrorResp(c, err)
		return
	}

	if err = c.ShouldBindJSON(&params); err != nil {
		helper.ErrorResp(c, errors.NewErrValidation(err))
		return
	}

	if params.SpecialPrice != 0 {
		strArr := []string{params.SpecialPriceStartStr, params.SpecialPriceEndStr}
		timeArr, parseErr := helper.ParseUTC8Datetime(strArr)
		if parseErr != nil {
			helper.ErrorResp(c, errors.ErrWrongTimeFormat)
			return
		}
		params.SpecialPriceStart = timeArr[0]
		params.SpecialPriceEnd = timeArr[1]
	}

	params.ProductID = id

	admin, err := h.adminSvc.GetCurrentAdmin(ctx)
	if err != nil {
		helper.ErrorResp(c, err)
		return
	}

	params.MerchantID = admin.Merchant.ID

	if err = h.productSvc.UpdateProduct(ctx, params); err != nil {
		helper.ErrorResp(c, err)
		return
	}

	ginTool.Success(c)
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	var (
		err error
		ctx = c.Request.Context()
	)

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.ErrorResp(c, err)
		return
	}

	admin, err := h.adminSvc.GetCurrentAdmin(ctx)
	if err != nil {
		helper.ErrorResp(c, err)
		return
	}

	if err = h.productSvc.DeleteProduct(ctx, id, admin.Merchant.ID); err != nil {
		helper.ErrorResp(c, err)
		return
	}

	ginTool.Success(c)
}

func (h *ProductHandler) SwitchProductStatus(c *gin.Context) {
	var (
		err error
		ctx = c.Request.Context()
	)

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.ErrorResp(c, err)
		return
	}

	admin, err := h.adminSvc.GetCurrentAdmin(ctx)
	if err != nil {
		helper.ErrorResp(c, err)
		return
	}

	if err = h.productSvc.SwitchProductStatus(ctx, id, admin.Merchant.ID); err != nil {
		helper.ErrorResp(c, err)
		return
	}

	ginTool.Success(c)
}

func (h *ProductHandler) ListProductByFilter(c *gin.Context) {
	var (
		err    error
		params vo.ListProductByFilterParams
		ctx    = c.Request.Context()
	)

	if err = c.ShouldBindQuery(&params); err != nil {
		helper.ErrorResp(c, errors.NewErrValidation(err))
		return
	}

	merchant, err := h.merchantSvc.GetCurrentMerchant(ctx)
	if err != nil {
		helper.ErrorResp(c, err)
		return
	}

	params.MerchantID = merchant.ID

	products, total, err := h.productSvc.ListProductByFilter(ctx, params)
	if err != nil {
		helper.ErrorResp(c, err)
		return
	}

	params.Pagination.Total = total

	ginTool.SuccessWithPagination(c, products, params.Pagination)
}

func (h *ProductHandler) GetProductDetail(c *gin.Context) {
	var (
		err error
		ctx = c.Request.Context()
	)

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.ErrorResp(c, err)
		return
	}

	merchant, err := h.merchantSvc.GetCurrentMerchant(ctx)
	if err != nil {
		helper.ErrorResp(c, err)
		return
	}

	detail, err := h.productSvc.GetProductDetail(ctx, id, merchant.ID)
	if err != nil {
		helper.ErrorResp(c, err)
		return
	}

	ginTool.SuccessWithData(c, detail)
}

/* Product Spec */

func (h *ProductHandler) GetProductSpec(c *gin.Context) {
	var (
		err error
		ctx = c.Request.Context()
	)

	productId, err := strconv.ParseUint(c.Param("productId"), 10, 64)
	if err != nil {
		helper.ErrorResp(c, err)
		return
	}

	admin, err := h.adminSvc.GetCurrentAdmin(ctx)
	if err != nil {
		helper.ErrorResp(c, err)
		return
	}

	spec, err := h.productSvc.GetProductSpec(ctx, productId, admin.Merchant.ID)
	if err != nil {
		helper.ErrorResp(c, err)
		return
	}

	ginTool.SuccessWithData(c, spec)
}

func (h *ProductHandler) UpsertProductSpec(c *gin.Context) {
	var (
		err    error
		params vo.UpsertProductSpecParams
		ctx    = c.Request.Context()
	)

	if err = c.ShouldBindJSON(&params); err != nil {
		helper.ErrorResp(c, errors.NewErrValidation(err))
		return
	}

	productId, err := strconv.ParseUint(c.Param("productId"), 10, 64)
	if err != nil {
		helper.ErrorResp(c, err)
		return
	}

	admin, err := h.adminSvc.GetCurrentAdmin(ctx)
	if err != nil {
		helper.ErrorResp(c, err)
		return
	}

	params.ProductID = productId
	params.MerchantID = admin.Merchant.ID

	if err = h.productSvc.UpsertProductSpec(ctx, params); err != nil {
		helper.ErrorResp(c, err)
		return
	}

	ginTool.Success(c)
}

func (h *ProductHandler) DeleteProductSpec(c *gin.Context) {
	var (
		err error
		ctx = c.Request.Context()
	)

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.ErrorResp(c, err)
		return
	}

	admin, err := h.adminSvc.GetCurrentAdmin(ctx)
	if err != nil {
		helper.ErrorResp(c, err)
		return
	}

	if err = h.productSvc.DeleteProductSpec(ctx, id, admin.Merchant.ID); err != nil {
		helper.ErrorResp(c, err)
		return
	}

	ginTool.Success(c)
}

/* Product Stock */

func (h *ProductHandler) GetProductStock(c *gin.Context) {
	var (
		err error
		ctx = c.Request.Context()
	)

	productId, err := strconv.ParseUint(c.Param("productId"), 10, 64)
	if err != nil {
		helper.ErrorResp(c, err)
		return
	}

	admin, err := h.adminSvc.GetCurrentAdmin(ctx)
	if err != nil {
		helper.ErrorResp(c, err)
		return
	}

	stocks, err := h.productSvc.GetProductStock(ctx, productId, admin.Merchant.ID)
	if err != nil {
		helper.ErrorResp(c, err)
		return
	}

	ginTool.SuccessWithData(c, stocks)
}

func (h *ProductHandler) UpsertProductStock(c *gin.Context) {
	var (
		err    error
		params vo.UpsertProductStockParams
		ctx    = c.Request.Context()
	)

	if err = c.ShouldBindJSON(&params); err != nil {
		helper.ErrorResp(c, errors.NewErrValidation(err))
		return
	}

	productId, err := strconv.ParseUint(c.Param("productId"), 10, 64)
	if err != nil {
		helper.ErrorResp(c, err)
		return
	}

	admin, err := h.adminSvc.GetCurrentAdmin(ctx)
	if err != nil {
		helper.ErrorResp(c, err)
		return
	}

	params.ProductID = productId
	params.MerchantID = admin.Merchant.ID

	if err = h.productSvc.UpsertProductStock(ctx, params); err != nil {
		helper.ErrorResp(c, err)
		return
	}

	ginTool.Success(c)
}

func (h *ProductHandler) DeleteProductStock(c *gin.Context) {
	var (
		err error
		ctx = c.Request.Context()
	)

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.ErrorResp(c, err)
		return
	}

	admin, err := h.adminSvc.GetCurrentAdmin(ctx)
	if err != nil {
		helper.ErrorResp(c, err)
		return
	}

	if err = h.productSvc.DeleteProductStock(ctx, id, admin.Merchant.ID); err != nil {
		helper.ErrorResp(c, err)
		return
	}

	ginTool.Success(c)
}
