package cmd

import (
	"ecommerce-api/internal/config"
	"ecommerce-api/internal/db"
	"ecommerce-api/internal/logger"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/spf13/cobra"
	"os"
)

func NewMigrateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "migrate",
		Run: runMigrate,
	}

	cmd.Flags().StringVarP(&configPath, "config", "c", "", "config file location")

	return cmd
}

func runMigrate(cobraCmd *cobra.Command, args []string) {
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

	d, err := mysql.WithInstance(dbClient, &mysql.Config{})
	if err != nil {
		logger.Logger.Errorf("Init database error: %s", err)
		os.Exit(1)
	}

	instance, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s", cfg.App.MigrationPath), "mysql", d)
	if err != nil {
		logger.Logger.Errorf("Init database instance error: %s", err)
		os.Exit(1)
	}

	err = instance.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		logger.Logger.Infof("Migrate database success")
		os.Exit(0)
	}

	if err != nil {
		logger.Logger.Errorf("Migrate database error: %s", err)
		os.Exit(1)
	}
}
