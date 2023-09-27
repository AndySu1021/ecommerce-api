package pkg

import (
	"database/sql"
	"ecommerce-api/internal/config"
	"ecommerce-api/internal/db/model"
	"ecommerce-api/internal/logger"
	"ecommerce-api/internal/middleware"
	"ecommerce-api/internal/storage/driver"
	app_product "ecommerce-api/pkg/catalog/product/application"
	persis_product "ecommerce-api/pkg/catalog/product/infrastructure/persistence"
	present_product "ecommerce-api/pkg/catalog/product/presentation"
	"ecommerce-api/pkg/common/infrastructure/email"
	present_common "ecommerce-api/pkg/common/presentation"
	app_admin "ecommerce-api/pkg/identity/admin/application"
	persis_admin "ecommerce-api/pkg/identity/admin/infrastructure/persistence"
	present_admin "ecommerce-api/pkg/identity/admin/presentation"
	app_member "ecommerce-api/pkg/identity/member/application"
	"ecommerce-api/pkg/identity/member/domain/event"
	"ecommerce-api/pkg/identity/member/infrastructure/persistence"
	present_member "ecommerce-api/pkg/identity/member/presentation"
	app_merchant "ecommerce-api/pkg/identity/merchant/application"
	persis_merchant "ecommerce-api/pkg/identity/merchant/infrastructure/persistence"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
)

type ModuleParams struct {
	Engine  *gin.Engine
	Redis   redis.UniversalClient
	DB      *sql.DB
	Storage driver.Client
	Config  *config.Config
	Mailgun *email.Mailgun
}

func InitModules(params ModuleParams) {
	// New queries
	queries := model.New(params.DB)

	// Exec migration
	if err := execMigration(params.DB, params.Config.App.MigrationPath); err != nil {
		logger.Logger.Errorf("Exec migration error: %s", err)
		return
	}

	// Exec seed
	if err := execSeed(queries); err != nil {
		logger.Logger.Errorf("Exec seed error: %s", err)
		return
	}

	s := NewService(params, queries)
	h := NewHandler(s, params.Config.Storage.BaseUrl, params.Storage)
	InitTransport(params.Engine, h)
}

type S struct {
	merchantSvc app_merchant.MerchantService
	memberSvc   app_member.MemberService
	adminSvc    app_admin.AdminService
	productSvc  app_product.ProductService
}

type H struct {
	memberHandler  *present_member.MemberHandler
	commonHandler  *present_common.CommonHandler
	adminHandler   *present_admin.AdminHandler
	productHandler *present_product.ProductHandler
}

func NewService(params ModuleParams, queries *model.Queries) *S {
	s := &S{}
	//conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	//if err != nil {
	//	log.Fatalf("did not connect: %v", err)
	//}
	s.merchantSvc = app_merchant.NewMerchantService(params.Redis, persis_merchant.NewMerchantRepository(params.DB, queries))
	s.memberSvc = app_member.NewMemberService(params.Redis, persistence.NewMemberRepository(params.DB, queries), event.NewHandler(params.Mailgun), s.merchantSvc)
	s.adminSvc = app_admin.NewAdminService(params.Redis, persis_admin.NewAdminRepository(params.DB, queries), s.merchantSvc)
	s.productSvc = app_product.NewProductService(params.Config.Storage.BaseUrl, persis_product.NewProductRepository(params.DB, queries))
	return s
}

func NewHandler(svc *S, baseUrl string, storage driver.Client) *H {
	h := &H{}
	h.memberHandler = present_member.NewMemberHandler(svc.merchantSvc, svc.memberSvc, svc.adminSvc)
	h.commonHandler = present_common.NewCommonHandler(baseUrl, storage)
	h.adminHandler = present_admin.NewAdminHandler(svc.adminSvc, svc.merchantSvc)
	h.productHandler = present_product.NewProductHandler(svc.adminSvc, svc.merchantSvc, svc.productSvc)
	return h
}

func InitTransport(e *gin.Engine, h *H) {
	e.Use(middleware.CORS())

	present_member.InitTransport(e, h.memberHandler)
	present_common.InitTransport(e, h.commonHandler)
	present_admin.InitTransport(e, h.adminHandler)
	present_product.InitTransport(e, h.productHandler)
}

func execMigration(db *sql.DB, path string) error {
	d, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return err
	}
	instance, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s", path), "mysql", d)
	if err != nil {
		return err
	}
	err = instance.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	}
	return err
}

func execSeed(queries *model.Queries) (err error) {
	return nil
}
