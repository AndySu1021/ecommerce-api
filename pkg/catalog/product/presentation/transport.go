package presentation

import (
	"ecommerce-api/internal/instrument"
	"ecommerce-api/internal/middleware"
	"github.com/gin-gonic/gin"
)

func InitTransport(e *gin.Engine, h *ProductHandler) {
	// 後台 api
	adminRoutes := e.Group("/admin/api", middleware.SetCurrentHost())
	{
		/* Product Category */

		adminRoutes.GET("/product-categories",
			instrument.RequestMetric("ListProductCategory"),
			h.ListProductCategory,
		)

		adminRoutes.POST("/product-category",
			instrument.RequestMetric("CreateProductCategory"),
			h.CreateProductCategory,
		)

		adminRoutes.GET("/product-category/:id",
			instrument.RequestMetric("GetProductCategory"),
			h.GetProductCategory,
		)

		adminRoutes.PUT("/product-category/:id",
			instrument.RequestMetric("UpdateProductCategory"),
			h.UpdateProductCategory,
		)

		adminRoutes.DELETE("/product-category/:id",
			instrument.RequestMetric("DeleteProductCategory"),
			h.DeleteProductCategory,
		)

		/* Product */

		adminRoutes.GET("/products",
			instrument.RequestMetric("ListProduct"),
			h.ListProduct,
		)

		adminRoutes.POST("/product",
			instrument.RequestMetric("CreateProduct"),
			h.CreateProduct,
		)

		adminRoutes.GET("/product/:id",
			instrument.RequestMetric("GetProduct"),
			h.GetProduct,
		)

		adminRoutes.PUT("/product/:id",
			instrument.RequestMetric("UpdateProduct"),
			h.UpdateProduct,
		)

		adminRoutes.DELETE("/product/:id",
			instrument.RequestMetric("DeleteProduct"),
			h.DeleteProduct,
		)

		adminRoutes.PATCH("/product/:id/status",
			instrument.RequestMetric("SwitchProductStatus"),
			h.SwitchProductStatus,
		)

		/* Product Spec */

		adminRoutes.GET("/product-spec/:productId",
			instrument.RequestMetric("GetProductSpec"),
			h.GetProductSpec,
		)

		adminRoutes.PUT("/product-spec/:productId",
			instrument.RequestMetric("UpsertProductSpec"),
			h.UpsertProductSpec,
		)

		adminRoutes.DELETE("/product-spec/:id",
			instrument.RequestMetric("DeleteProductSpec"),
			h.DeleteProductSpec,
		)

		/* Product Stock */

		adminRoutes.GET("/product-stock/:productId",
			instrument.RequestMetric("GetProductStock"),
			h.GetProductStock,
		)

		adminRoutes.PUT("/product-stock/:productId",
			instrument.RequestMetric("UpsertProductStock"),
			h.UpsertProductStock,
		)

		adminRoutes.DELETE("/product-stock/:id",
			instrument.RequestMetric("DeleteProductStock"),
			h.DeleteProductStock,
		)
	}

	// 前台 api
	routes := e.Group("/api", middleware.SetCurrentHost())
	{
		/* Product */

		routes.GET("/products",
			instrument.RequestMetric("ListProductByFilter"),
			h.ListProductByFilter,
		)

		routes.GET("/product/:id",
			instrument.RequestMetric("GetProductDetail"),
			h.GetProductDetail,
		)
	}
}
