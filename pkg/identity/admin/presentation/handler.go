package presentation

import (
	"ecommerce-api/internal/errors"
	"ecommerce-api/internal/helper"
	"ecommerce-api/pkg/identity/admin/application"
	"ecommerce-api/pkg/identity/admin/domain/vo"
	app_merchant "ecommerce-api/pkg/identity/merchant/application"
	ginTool "github.com/AndySu1021/go-util/gin"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	adminSvc    application.AdminService
	merchantSvc app_merchant.MerchantService
}

func NewAdminHandler(adminSvc application.AdminService, merchantSvc app_merchant.MerchantService) *AdminHandler {
	return &AdminHandler{
		adminSvc:    adminSvc,
		merchantSvc: merchantSvc,
	}
}

func (h *AdminHandler) Login(c *gin.Context) {
	var (
		err    error
		params vo.LoginParams
		ctx    = c.Request.Context()
	)

	if err = c.ShouldBindJSON(&params); err != nil {
		helper.ErrorResp(c, errors.NewErrValidation(err))
		return
	}

	merchant, err := h.merchantSvc.GetCurrentMerchant(ctx)
	if err != nil {
		helper.ErrorResp(c, err)
		return
	}

	params.MerchantID = merchant.ID
	params.MerchantHost = merchant.Host
	params.MerchantEncryptSalt = merchant.EncryptSalt

	admin, err := h.adminSvc.Login(ctx, params)
	if err != nil {
		helper.ErrorResp(c, err)
		return
	}

	ginTool.SuccessWithData(c, admin)
}

func (h *AdminHandler) Logout(c *gin.Context) {
	var (
		err    error
		params vo.LogoutParams
		ctx    = c.Request.Context()
	)

	admin, err := h.adminSvc.GetCurrentAdmin(ctx)
	if err != nil {
		helper.ErrorResp(c, err)
		return
	}

	params.Email = admin.Email
	params.MerchantHost = admin.Merchant.Host

	if err = h.adminSvc.Logout(ctx, params); err != nil {
		helper.ErrorResp(c, err)
		return
	}

	ginTool.Success(c)
}

func (h *AdminHandler) CheckToken(c *gin.Context) {
	var (
		err error
		ctx = c.Request.Context()
	)

	admin, err := h.adminSvc.GetCurrentAdmin(ctx)
	if err != nil {
		helper.ErrorResp(c, err)
		return
	}

	ginTool.SuccessWithData(c, admin)
}
