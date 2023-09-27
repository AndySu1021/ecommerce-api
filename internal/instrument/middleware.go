package instrument

import (
	"github.com/gin-gonic/gin"
	"time"
)

func RequestMetric(op string) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func(begin time.Time) {
			instance.Request.RequestCount.WithLabelValues(op).Add(1)
			instance.Request.RequestLatency.WithLabelValues(op).Observe(time.Since(begin).Seconds())
		}(time.Now())
	}
}
