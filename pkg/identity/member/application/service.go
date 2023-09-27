package application

import (
	"context"
	"crypto/md5"
	"database/sql"
	"ecommerce-api/internal/helper"
	"ecommerce-api/pkg/constant"
	"ecommerce-api/pkg/identity/member/domain/event"
	"ecommerce-api/pkg/identity/member/domain/repository"
	"ecommerce-api/pkg/identity/member/domain/vo"
	"ecommerce-api/pkg/identity/merchant/application"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/go-sql-driver/mysql"
	"github.com/rs/xid"
	"math/rand"
	"strconv"
	"time"
)

//go:embed lua/login.lua
var loginScript string

//go:embed lua/logout.lua
var logoutScript string

type MemberService interface {
	GetCurrentMember(ctx context.Context) (vo.CurrentMember, error)
	Register(ctx context.Context, params vo.RegisterParams) (vo.CurrentMember, error)
	Login(ctx context.Context, params vo.LoginParams) (vo.CurrentMember, error)
	Logout(ctx context.Context, params vo.LogoutParams) error
	ForgetPassword(ctx context.Context, params vo.ForgetPasswordParams) error
	CheckForgetCode(ctx context.Context, params vo.CheckForgetCodeParams) error
	ResetPassword(ctx context.Context, params vo.ResetPasswordParams) error
	GetMemberInfo(ctx context.Context, memberId uint64) (vo.MemberInfo, error)
	UpdateMemberInfo(ctx context.Context, params vo.UpdateMemberInfoParams) error
	UpdateMemberPassword(ctx context.Context, params vo.UpdateMemberPasswordParams) error
	ListMember(ctx context.Context, params vo.ListMemberParams) ([]vo.ListMemberRow, int64, error)
	UpdateMemberStatus(ctx context.Context, params vo.UpdateMemberStatusParams) error
}

type service struct {
	loginScript  *redis.Script
	logoutScript *redis.Script
	rdb          redis.UniversalClient
	repo         repository.MemberRepository
	eventHandler *event.Handler
	merchantSvc  application.MerchantService
}

func NewMemberService(rdb redis.UniversalClient, repo repository.MemberRepository, handler *event.Handler, merchantSvc application.MerchantService) MemberService {
	return &service{
		loginScript:  redis.NewScript(loginScript),
		logoutScript: redis.NewScript(logoutScript),
		rdb:          rdb,
		repo:         repo,
		eventHandler: handler,
		merchantSvc:  merchantSvc,
	}
}

func (s *service) GetCurrentMember(ctx context.Context) (vo.CurrentMember, error) {
	merchant, err := s.merchantSvc.GetCurrentMerchant(ctx)
	if err != nil {
		return vo.CurrentMember{}, err
	}

	token, ok := ctx.Value(constant.ContextKeyToken).(string)
	if token == "" || !ok {
		return vo.CurrentMember{}, fmt.Errorf("invalid token")
	}

	redisKey := fmt.Sprintf("token:member:%s:%s", merchant.Host, token)
	result, err := s.rdb.Get(ctx, redisKey).Result()
	if err != nil {
		return vo.CurrentMember{}, err
	}

	var tmp vo.CurrentMember
	if err = json.Unmarshal([]byte(result), &tmp); err != nil {
		return vo.CurrentMember{}, err
	}

	return tmp, nil
}

func (s *service) Register(ctx context.Context, params vo.RegisterParams) (vo.CurrentMember, error) {
	if err := s.repo.CreateMember(ctx, vo.CreateMemberParams{
		MerchantID: params.MerchantID,
		Email:      params.Email,
		Password:   helper.SaltEncrypt(params.Password, params.MerchantEncryptSalt),
	}); err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return vo.CurrentMember{}, fmt.Errorf("信箱已存在")
		}
		return vo.CurrentMember{}, err
	}

	return s.Login(ctx, *(*vo.LoginParams)(&params))
}

func (s *service) Login(ctx context.Context, params vo.LoginParams) (vo.CurrentMember, error) {
	member, err := s.repo.MemberLogin(ctx, vo.MemberLoginParams{
		MerchantID: params.MerchantID,
		Email:      params.Email,
		Password:   helper.SaltEncrypt(params.Password, params.MerchantEncryptSalt),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return vo.CurrentMember{}, fmt.Errorf("帳號或密碼錯誤")
		}
		return vo.CurrentMember{}, err
	}

	// 檢查會員狀態
	if member.IsEnabled == constant.No {
		return vo.CurrentMember{}, fmt.Errorf("禁止登入")
	}

	// 更新登入時間
	if err = s.repo.UpdateMemberByMap(ctx, member.ID, map[string]interface{}{
		"last_login_time": time.Now().UTC(),
	}); err != nil {
		return vo.CurrentMember{}, err
	}

	token := genToken()
	currentMember := vo.CurrentMember{
		ID:    member.ID,
		Email: member.Email,
		Token: token,
		Merchant: vo.Merchant{
			ID:          params.MerchantID,
			EncryptSalt: params.MerchantEncryptSalt,
		},
	}

	result, err := json.Marshal(currentMember)
	if err != nil {
		return vo.CurrentMember{}, err
	}

	expire := int64(24 * time.Hour / time.Second)
	if err = s.loginScript.Run(ctx, s.rdb, []string{params.MerchantHost}, member.Email, token, result, expire).Err(); err != nil {
		return vo.CurrentMember{}, err
	}

	return currentMember, nil
}

func (s *service) Logout(ctx context.Context, params vo.LogoutParams) error {
	return s.logoutScript.Run(ctx, s.rdb, []string{params.MerchantHost}, params.Email).Err()
}

func (s *service) ForgetPassword(ctx context.Context, params vo.ForgetPasswordParams) error {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := strconv.FormatInt(100000+r.Int63n(900000), 10)

	if _, err := s.repo.CheckEmailExist(ctx, params.MerchantID, params.Email); err != nil {
		return err
	}

	result, err := s.rdb.SetNX(ctx, fmt.Sprintf("forget:password:%s:%s", params.MerchantHost, params.Email), code, 3*time.Minute).Result()
	if err != nil {
		return err
	}
	if !result {
		return fmt.Errorf("認證信已寄出，請稍後再試")
	}

	if err = s.rdb.Set(ctx, fmt.Sprintf("forget:password:%s:%s", params.MerchantHost, code), params.Email, 10*time.Minute).Err(); err != nil {
		return err
	}

	if err = s.eventHandler.Handle(event.MemberForgetPassword{
		LogoUrl:      "",
		MerchantName: params.MerchantName,
		RealName:     "",
		Email:        "",
		Code:         code,
	}); err != nil {
		return err
	}

	return nil
}

func (s *service) CheckForgetCode(ctx context.Context, params vo.CheckForgetCodeParams) error {
	redisKey := fmt.Sprintf("forget:password:%s:%s", params.MerchantHost, params.Code)
	email, err := s.rdb.Get(ctx, redisKey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}

	if errors.Is(err, redis.Nil) || email != params.Email {
		return fmt.Errorf("驗證碼錯誤")
	}

	return nil
}

func (s *service) ResetPassword(ctx context.Context, params vo.ResetPasswordParams) error {
	redisKey := fmt.Sprintf("forget:password:%s:%s", params.MerchantHost, params.Code)
	email, err := s.rdb.Get(ctx, redisKey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}

	if errors.Is(err, redis.Nil) {
		return fmt.Errorf("驗證碼錯誤")
	}

	if err = s.repo.ResetPassword(ctx, vo.ResetPasswordRepoParams{
		MerchantID: params.MerchantID,
		Email:      email,
		Password:   helper.SaltEncrypt(params.Password, params.MerchantEncryptSalt),
	}); err != nil {
		return err
	}

	return s.rdb.Del(ctx, redisKey).Err()
}

func (s *service) GetMemberInfo(ctx context.Context, memberId uint64) (vo.MemberInfo, error) {
	return s.repo.GetMemberInfo(ctx, memberId)
}

func (s *service) UpdateMemberInfo(ctx context.Context, params vo.UpdateMemberInfoParams) error {
	return s.repo.UpdateMemberByMap(ctx, params.MemberID, map[string]interface{}{
		"real_name":  params.RealName,
		"mobile":     params.Mobile,
		"birthday":   params.BirthdayTime,
		"sex":        params.Sex,
		"city":       params.City,
		"district":   params.District,
		"address":    params.Address,
		"zip_code":   params.ZipCode,
		"updated_at": time.Now().UTC(),
	})
}

func (s *service) UpdateMemberPassword(ctx context.Context, params vo.UpdateMemberPasswordParams) error {
	return s.repo.UpdateMemberPassword(ctx, vo.UpdateMemberPasswordRepoParams{
		OldPassword: helper.SaltEncrypt(params.OldPassword, params.MerchantEncryptSalt),
		Password:    helper.SaltEncrypt(params.Password, params.MerchantEncryptSalt),
		MemberID:    params.MemberID,
	})
}

func (s *service) ListMember(ctx context.Context, params vo.ListMemberParams) ([]vo.ListMemberRow, int64, error) {
	return s.repo.ListMember(ctx, params)
}

func (s *service) UpdateMemberStatus(ctx context.Context, params vo.UpdateMemberStatusParams) error {
	return s.repo.UpdateMemberByMap(ctx, params.MemberID, map[string]interface{}{
		"is_enabled": params.IsEnabled,
		"updated_at": time.Now().UTC(),
	})
}

func genToken() string {
	str := time.Now().String() + xid.New().String()
	str = fmt.Sprintf("%x", md5.Sum([]byte(str)))
	return str
}
