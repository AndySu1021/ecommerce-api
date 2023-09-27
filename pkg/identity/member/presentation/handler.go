package presentation

import (
	"database/sql"
	"ecommerce-api/internal/errors"
	"ecommerce-api/internal/helper"
	app_admin "ecommerce-api/pkg/identity/admin/application"
	"ecommerce-api/pkg/identity/member/application"
	"ecommerce-api/pkg/identity/member/domain/vo"
	app_merchant "ecommerce-api/pkg/identity/merchant/application"
	ginTool "github.com/AndySu1021/go-util/gin"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

type MemberHandler struct {
	merchantSvc app_merchant.MerchantService
	memberSvc   application.MemberService
	adminSvc    app_admin.AdminService
}

func NewMemberHandler(merchantSvc app_merchant.MerchantService, memberSvc application.MemberService, adminSvc app_admin.AdminService) *MemberHandler {
	return &MemberHandler{
		merchantSvc: merchantSvc,
		memberSvc:   memberSvc,
		adminSvc:    adminSvc,
	}
}

func (h *MemberHandler) Register(c *gin.Context) {
	var (
		err    error
		params vo.RegisterParams
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

	member, err := h.memberSvc.Register(ctx, params)
	if err != nil {
		helper.ErrorResp(c, err)
		return
	}

	ginTool.SuccessWithData(c, member)
}

func (h *MemberHandler) Login(c *gin.Context) {
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

	member, err := h.memberSvc.Login(ctx, params)
	if err != nil {
		helper.ErrorResp(c, err)
		return
	}

	ginTool.SuccessWithData(c, member)
}

func (h *MemberHandler) Logout(c *gin.Context) {
	var (
		err    error
		params vo.LogoutParams
		ctx    = c.Request.Context()
	)

	member, err := h.memberSvc.GetCurrentMember(ctx)
	if err != nil {
		helper.ErrorResp(c, err)
		return
	}

	params.Email = member.Email
	params.MerchantHost = member.Merchant.Host

	if err = h.memberSvc.Logout(ctx, params); err != nil {
		helper.ErrorResp(c, err)
		return
	}

	ginTool.Success(c)
}

func (h *MemberHandler) CheckToken(c *gin.Context) {
	var (
		err error
		ctx = c.Request.Context()
	)

	member, err := h.memberSvc.GetCurrentMember(ctx)
	if err != nil {
		helper.ErrorResp(c, err)
		return
	}

	ginTool.SuccessWithData(c, member)
}

func (h *MemberHandler) ForgetPassword(c *gin.Context) {
	var (
		err    error
		params vo.ForgetPasswordParams
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
	params.MerchantName = merchant.Name
	params.MerchantHost = merchant.Host

	if err = h.memberSvc.ForgetPassword(ctx, params); err != nil {
		helper.ErrorResp(c, err)
		return
	}

	ginTool.Success(c)
}

func (h *MemberHandler) CheckForgetCode(c *gin.Context) {
	var (
		err    error
		params vo.CheckForgetCodeParams
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

	params.MerchantHost = merchant.Host

	if err = h.memberSvc.CheckForgetCode(ctx, params); err != nil {
		helper.ErrorResp(c, err)
		return
	}

	ginTool.Success(c)
}

func (h *MemberHandler) ResetPassword(c *gin.Context) {
	var (
		err    error
		params vo.ResetPasswordParams
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

	if err = h.memberSvc.ResetPassword(ctx, params); err != nil {
		helper.ErrorResp(c, err)
		return
	}

	ginTool.Success(c)
}

func (h *MemberHandler) GetMemberInfo(c *gin.Context) {
	var (
		err error
		ctx = c.Request.Context()
	)

	member, err := h.memberSvc.GetCurrentMember(ctx)
	if err != nil {
		helper.ErrorResp(c, err)
		return
	}

	info, err := h.memberSvc.GetMemberInfo(ctx, member.ID)
	if err != nil {
		helper.ErrorResp(c, err)
		return
	}

	ginTool.SuccessWithData(c, info)
}

func (h *MemberHandler) UpdateMemberInfo(c *gin.Context) {
	var (
		err    error
		params vo.UpdateMemberInfoParams
		ctx    = c.Request.Context()
	)

	if err = c.ShouldBindJSON(&params); err != nil {
		helper.ErrorResp(c, errors.NewErrValidation(err))
		return
	}

	birthday, err := time.Parse("2006-01-02", params.Birthday)
	if err != nil {
		helper.ErrorResp(c, errors.ErrWrongTimeFormat)
		return
	}

	params.BirthdayTime = sql.NullTime{
		Time:  birthday,
		Valid: true,
	}

	member, err := h.memberSvc.GetCurrentMember(ctx)
	if err != nil {
		helper.ErrorResp(c, err)
		return
	}

	params.MemberID = member.ID

	if err = h.memberSvc.UpdateMemberInfo(ctx, params); err != nil {
		helper.ErrorResp(c, err)
		return
	}

	ginTool.Success(c)
}

func (h *MemberHandler) UpdateMemberPassword(c *gin.Context) {
	var (
		err    error
		params vo.UpdateMemberPasswordParams
		ctx    = c.Request.Context()
	)

	if err = c.ShouldBindJSON(&params); err != nil {
		helper.ErrorResp(c, errors.NewErrValidation(err))
		return
	}

	member, err := h.memberSvc.GetCurrentMember(ctx)
	if err != nil {
		helper.ErrorResp(c, err)
		return
	}

	params.MemberID = member.ID
	params.MerchantEncryptSalt = member.Merchant.EncryptSalt

	if err = h.memberSvc.UpdateMemberPassword(ctx, params); err != nil {
		helper.ErrorResp(c, err)
		return
	}

	ginTool.Success(c)
}

func (h *MemberHandler) ListMember(c *gin.Context) {
	var (
		err    error
		params vo.ListMemberParams
		ctx    = c.Request.Context()
	)

	if err = c.ShouldBindQuery(&params); err != nil {
		helper.ErrorResp(c, errors.NewErrValidation(err))
		return
	}

	if _, err = h.adminSvc.GetCurrentAdmin(ctx); err != nil {
		helper.ErrorResp(c, err)
		return
	}

	members, total, err := h.memberSvc.ListMember(ctx, params)
	if err != nil {
		helper.ErrorResp(c, err)
		return
	}

	params.Pagination.Total = total

	ginTool.SuccessWithPagination(c, members, params.Pagination)
}

func (h *MemberHandler) UpdateMemberStatus(c *gin.Context) {
	var (
		err    error
		params vo.UpdateMemberStatusParams
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

	params.MemberID = id

	if _, err = h.adminSvc.GetCurrentAdmin(ctx); err != nil {
		helper.ErrorResp(c, err)
		return
	}

	if err = h.memberSvc.UpdateMemberStatus(ctx, params); err != nil {
		helper.ErrorResp(c, err)
		return
	}

	ginTool.Success(c)
}
