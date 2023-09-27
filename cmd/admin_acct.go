package cmd

import (
	"context"
	"ecommerce-api/internal/config"
	"ecommerce-api/internal/db"
	"ecommerce-api/internal/db/model"
	"ecommerce-api/internal/helper"
	"ecommerce-api/internal/logger"
	"ecommerce-api/pkg/constant"
	admin_vo "ecommerce-api/pkg/identity/admin/domain/vo"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/spf13/cobra"
	"os"
	"time"
)

var (
	merchantId    uint64
	adminEmail    string
	adminPassword string
)

func NewAdminCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "admin",
		Run: runAdmin,
	}

	cmd.Flags().Uint64VarP(&merchantId, "merchant", "m", 0, "merchant ID")
	cmd.Flags().StringVarP(&adminEmail, "email", "e", "", "admin account email")
	cmd.Flags().StringVarP(&adminPassword, "password", "p", "", "admin account password")
	cmd.Flags().StringVarP(&configPath, "config", "c", "", "config file location")

	return cmd
}

func runAdmin(cobraCmd *cobra.Command, args []string) {
	ctx := context.Background()

	// Init config
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Init Logger
	if err = logger.InitZapLogger(cfg.Log); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Init database
	dbClient, err := db.NewDatabase(&cfg.Database)
	if err != nil {
		logger.Logger.Errorf("Init database error: %s", err)
		os.Exit(1)
	}

	queries := model.New(dbClient)

	merchant, err := queries.GetMerchant(ctx, merchantId)
	if err != nil {
		logger.Logger.Errorf("Get merchant error: %s", err)
		os.Exit(1)
	}

	if err = queries.CreateAdmin(ctx, model.CreateAdminParams{
		MerchantID: merchantId,
		Email:      adminEmail,
		Password:   helper.SaltEncrypt(adminPassword, merchant.EncryptSalt),
		RealName:   "超級管理員",
		Mobile:     "",
		Sex:        admin_vo.SexMale,
		IsEnabled:  constant.Yes,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}); err != nil {
		logger.Logger.Errorf("Create merchant admin error: %s", err)
		os.Exit(1)
	}

	logger.Logger.Info("Create merchant admin success")
}
