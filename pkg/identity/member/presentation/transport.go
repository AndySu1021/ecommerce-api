package presentation

import (
	"ecommerce-api/internal/instrument"
	"ecommerce-api/internal/middleware"
	"github.com/gin-gonic/gin"
)

func InitTransport(e *gin.Engine, h *MemberHandler) {
	// 前台 api
	routes := e.Group("/api", middleware.SetCurrentHost())
	{
		routes.POST("/member/register",
			instrument.RequestMetric("Register"),
			h.Register,
		)

		routes.POST("/member/login",
			instrument.RequestMetric("Login"),
			h.Login,
		)

		routes.POST("/member/logout",
			middleware.SetMemberToken(),
			instrument.RequestMetric("Logout"),
			h.Logout,
		)

		routes.GET("/member/check-token",
			middleware.SetMemberToken(),
			instrument.RequestMetric("CheckToken"),
			h.CheckToken,
		)

		routes.POST("/member/forget-password",
			instrument.RequestMetric("ForgetPassword"),
			h.ForgetPassword,
		)

		routes.POST("/member/check-forget-code",
			instrument.RequestMetric("CheckForgetCode"),
			h.CheckForgetCode,
		)

		routes.POST("/member/reset-password",
			instrument.RequestMetric("ResetPassword"),
			h.ResetPassword,
		)

		routes.GET("/member/basic-info",
			middleware.SetMemberToken(),
			instrument.RequestMetric("GetMemberInfo"),
			h.GetMemberInfo,
		)

		routes.PATCH("/member/basic-info",
			middleware.SetMemberToken(),
			instrument.RequestMetric("UpdateMemberInfo"),
			h.UpdateMemberInfo,
		)

		routes.PATCH("/member/password",
			middleware.SetMemberToken(),
			instrument.RequestMetric("UpdateMemberPassword"),
			h.UpdateMemberPassword,
		)
	}

	// 後台 api
	adminRoutes := e.Group("/admin/api", middleware.SetMemberToken())
	{
		adminRoutes.GET("/members",
			instrument.RequestMetric("ListMember"),
			h.ListMember,
		)

		adminRoutes.PATCH("/member/:id/status",
			instrument.RequestMetric("UpdateMemberStatus"),
			h.UpdateMemberStatus,
		)
	}
}
