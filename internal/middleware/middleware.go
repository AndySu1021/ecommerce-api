package middleware

import (
	"bytes"
	"context"
	"ecommerce-api/pkg/constant"
	"encoding/json"
	"fmt"
	ginTool "github.com/AndySu1021/go-util/gin"
	"github.com/AndySu1021/go-util/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"strings"
	"time"
)

func SetCurrentHost() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		host := c.GetHeader(constant.HeaderKeyHost.String())
		if host == "" {
			ginTool.ErrorAuth(c)
			c.Abort()
			return
		}

		valueCtx := context.WithValue(ctx, constant.ContextKeyHost, host)
		c.Request = c.Request.WithContext(valueCtx)
		c.Next()
	}
}

func SetMemberToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		token := c.GetHeader(constant.HeaderKeyMemberToken.String())
		if token == "" {
			ginTool.ErrorAuth(c)
			c.Abort()
			return
		}

		valueCtx := context.WithValue(ctx, constant.ContextKeyToken, token)
		c.Request = c.Request.WithContext(valueCtx)
		c.Next()
	}
}

func SetAdminToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		token := c.GetHeader(constant.HeaderKeyAdminToken.String())
		if token == "" {
			ginTool.ErrorAuth(c)
			c.Abort()
			return
		}

		valueCtx := context.WithValue(ctx, constant.ContextKeyAdminToken, token)
		c.Request = c.Request.WithContext(valueCtx)
		c.Next()
	}
}

func LogRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.Request.Header.Get(constant.HeaderKeyTraceID.String())
		if traceID == "" {
			traceID = uuid.New().String()
		}

		buf, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(buf))

		var requestPayload map[string]interface{}
		_ = json.Unmarshal(buf, &requestPayload)
		requestStr, err := json.Marshal(requestPayload)
		if err != nil {
			return
		}

		defer func(t time.Time) {
			if !strings.Contains(c.Request.URL.Path, "health") {
				latency := time.Since(t)
				logger.Logger.Infow("access log",
					"url", c.Request.URL.String(),
					"method", c.Request.Method,
					"latency", fmt.Sprintf("%d ms", latency.Milliseconds()),
					"request", string(requestStr),
					"status", c.Writer.Status(),
					"trace_id", traceID,
				)
			}
		}(time.Now())

		ctx := context.WithValue(c.Request.Context(), constant.ContextKeyTraceID, traceID)

		c.Request = c.Request.WithContext(ctx)
		c.Request.Header.Set(constant.HeaderKeyTraceID.String(), traceID)
		c.Writer.Header().Set(constant.HeaderKeyTraceID.String(), traceID)

		c.Next()
	}
}

func CORS() gin.HandlerFunc {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = []string{"*"}
	return cors.New(corsConfig)
}
