package presentation

import (
	"ecommerce-api/internal/instrument"
	"ecommerce-api/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func InitTransport(e *gin.Engine, h *CommonHandler) {
	//pprof.Register(e)

	e.GET("/metrics", gin.WrapH(promhttp.Handler()))
	e.GET("/health", func(c *gin.Context) { c.Status(http.StatusOK) })

	e.POST("/api/upload",
		middleware.SetAdminToken(),
		instrument.RequestMetric("UploadFile"),
		h.UploadFile,
	)
}
