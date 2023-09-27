package helper

import (
	"database/sql"
	my_err "ecommerce-api/internal/errors"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ErrorResp(c *gin.Context, err error) {
	// 格式化 mysql error
	tmpErr := err
	if errors.Is(err, sql.ErrNoRows) {
		tmpErr = my_err.ErrNotFound
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"code":    -1,
		"message": tmpErr.Error(),
		"data":    nil,
	})
}
