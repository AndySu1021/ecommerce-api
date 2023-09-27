package presentation

import (
	"ecommerce-api/internal/helper"
	"ecommerce-api/internal/storage/driver"
	"fmt"
	ginTool "github.com/AndySu1021/go-util/gin"
	"github.com/gin-gonic/gin"
)

type CommonHandler struct {
	baseUrl string
	storage driver.Client
}

func NewCommonHandler(baseUrl string, storage driver.Client) *CommonHandler {
	return &CommonHandler{
		baseUrl: baseUrl,
		storage: storage,
	}
}

func (h *CommonHandler) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		helper.ErrorResp(c, err)
		return
	}

	tmp, err := file.Open()
	if err != nil {
		helper.ErrorResp(c, err)
		return
	}
	defer tmp.Close()

	path, err := h.storage.Upload(c.Request.Context(), tmp, "upload", file.Filename)
	if err != nil {
		helper.ErrorResp(c, err)
		return
	}

	ginTool.SuccessWithData(c, gin.H{
		"url":                fmt.Sprintf("%s/%s", h.baseUrl, path),
		"url_without_domain": path,
	})
}
