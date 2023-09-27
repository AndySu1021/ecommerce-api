package application

import (
	"context"
	"crypto/md5"
	"database/sql"
	"ecommerce-api/internal/helper"
	"ecommerce-api/pkg/constant"
	"ecommerce-api/pkg/identity/admin/domain/repository"
	"ecommerce-api/pkg/identity/admin/domain/vo"
	"ecommerce-api/pkg/identity/merchant/application"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/rs/xid"
	"time"
)

//go:embed lua/login.lua
var loginScript string

//go:embed lua/logout.lua
var logoutScript string

type AdminService interface {
	GetCurrentAdmin(ctx context.Context) (vo.CurrentAdmin, error)
	Login(ctx context.Context, params vo.LoginParams) (vo.CurrentAdmin, error)
	Logout(ctx context.Context, params vo.LogoutParams) error
}

type service struct {
	loginScript  *redis.Script
	logoutScript *redis.Script
	rdb          redis.UniversalClient
	repo         repository.AdminRepository
	merchantSvc  application.MerchantService
}

func NewAdminService(rdb redis.UniversalClient, repo repository.AdminRepository, merchantSvc application.MerchantService) AdminService {
	return &service{
		loginScript:  redis.NewScript(loginScript),
		logoutScript: redis.NewScript(logoutScript),
		rdb:          rdb,
		repo:         repo,
		merchantSvc:  merchantSvc,
	}
}

func (s *service) GetCurrentAdmin(ctx context.Context) (vo.CurrentAdmin, error) {
	merchant, err := s.merchantSvc.GetCurrentMerchant(ctx)
	if err != nil {
		return vo.CurrentAdmin{}, err
	}

	token, ok := ctx.Value(constant.ContextKeyAdminToken).(string)
	if token == "" || !ok {
		return vo.CurrentAdmin{}, fmt.Errorf("invalid token")
	}

	redisKey := fmt.Sprintf("token:member:%s:%s", merchant.Host, token)
	result, err := s.rdb.Get(ctx, redisKey).Result()
	if err != nil {
		return vo.CurrentAdmin{}, err
	}

	var tmp vo.CurrentAdmin
	if err = json.Unmarshal([]byte(result), &tmp); err != nil {
		return vo.CurrentAdmin{}, err
	}

	return tmp, nil
}

func (s *service) Login(ctx context.Context, params vo.LoginParams) (vo.CurrentAdmin, error) {
	admin, err := s.repo.AdminLogin(ctx, vo.AdminLoginParams{
		MerchantID: params.MerchantID,
		Email:      params.Email,
		Password:   helper.SaltEncrypt(params.Password, params.MerchantEncryptSalt),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return vo.CurrentAdmin{}, fmt.Errorf("帳號或密碼錯誤")
		}
		return vo.CurrentAdmin{}, err
	}

	// 檢查管理員狀態
	if admin.IsEnabled == constant.No {
		return vo.CurrentAdmin{}, fmt.Errorf("禁止登入")
	}

	// 更新登入時間
	if err = s.repo.UpdateAdminByMap(ctx, admin.ID, map[string]interface{}{
		"last_login_time": time.Now().UTC(),
	}); err != nil {
		return vo.CurrentAdmin{}, err
	}

	token := genToken()
	currentAdmin := vo.CurrentAdmin{
		ID:    admin.ID,
		Email: admin.Email,
		Token: token,
		Merchant: vo.Merchant{
			ID:          params.MerchantID,
			Host:        params.MerchantHost,
			EncryptSalt: params.MerchantEncryptSalt,
		},
	}

	result, err := json.Marshal(currentAdmin)
	if err != nil {
		return vo.CurrentAdmin{}, err
	}

	expire := int64(24 * time.Hour / time.Second)
	if err = s.loginScript.Run(ctx, s.rdb, []string{params.MerchantHost}, admin.Email, token, result, expire).Err(); err != nil {
		return vo.CurrentAdmin{}, err
	}

	return currentAdmin, nil
}

func (s *service) Logout(ctx context.Context, params vo.LogoutParams) error {
	return s.logoutScript.Run(ctx, s.rdb, []string{params.MerchantHost}, params.Email).Err()
}

func genToken() string {
	str := time.Now().String() + xid.New().String()
	str = fmt.Sprintf("%x", md5.Sum([]byte(str)))
	return str
}
