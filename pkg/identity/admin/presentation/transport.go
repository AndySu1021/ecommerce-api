package presentation

import (
	"ecommerce-api/internal/instrument"
	"ecommerce-api/internal/middleware"
	"github.com/gin-gonic/gin"
)

func InitTransport(e *gin.Engine, h *AdminHandler) {
	// 後台 api
	adminRoutes := e.Group("/admin/api")
	{
		adminRoutes.POST("/admin/login",
			instrument.RequestMetric("AdminLogin"),
			h.Login,
		)

		adminRoutes.POST("/admin/logout",
			middleware.SetAdminToken(),
			instrument.RequestMetric("AdminLogout"),
			h.Logout,
		)

		adminRoutes.GET("/admin/check-token",
			middleware.SetAdminToken(),
			instrument.RequestMetric("AdminCheckToken"),
			h.CheckToken,
		)
	}
}
