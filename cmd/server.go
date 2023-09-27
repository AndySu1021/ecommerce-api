package cmd

import (
	"context"
	"ecommerce-api/internal/config"
	"ecommerce-api/internal/db"
	"ecommerce-api/internal/instrument"
	"ecommerce-api/internal/logger"
	"ecommerce-api/internal/redis"
	"ecommerce-api/internal/storage"
	"ecommerce-api/pkg"
	"ecommerce-api/pkg/common/infrastructure/email"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

var configPath string

func NewServerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "server",
		Run: runServer,
	}

	cmd.Flags().StringVarP(&configPath, "config", "c", "", "config file location")

	return cmd
}

func runServer(cobraCmd *cobra.Command, args []string) {
	// Recover
	defer handleRecover()

	// Global content
	ctx := context.Background()

	// Init config
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Init Logger
	if err = logger.InitZapLogger(cfg.Log); err != nil {
		fmt.Println(err)
		return
	}

	// Init database
	dbClient, err := db.NewDatabase(&cfg.Database)
	if err != nil {
		logger.Logger.Errorf("Init database error: %s", err)
		return
	}

	// Init redis
	redisClient, err := redis.NewRedisClient(cfg.Redis)
	if err != nil {
		logger.Logger.Errorf("Init redis error: %s", err)
		return
	}

	// Set server id
	serverId, err := redisClient.Incr(ctx, "server:id").Result()
	if err != nil {
		logger.Logger.Errorf("Get server id error: %s", err)
		return
	}
	cfg.App.ServerID = serverId % 1000

	// Init prometheus instrument
	instrument.InitInstrument(cfg.App)

	// Init http handler
	engine := newGin(cfg.Http.Mode)
	disk, err := storage.NewStorage(cfg.Storage)

	// Init Mailgun
	mailClient := email.NewMailgun(cfg.Mailgun)

	// Init module
	pkg.InitModules(pkg.ModuleParams{
		Engine:  engine,
		Redis:   redisClient,
		DB:      dbClient,
		Storage: disk,
		Config:  cfg,
		Mailgun: mailClient,
	})

	// Init http server
	exitCode := 0
	server := &http.Server{
		Addr:         ":" + cfg.Http.Port,
		Handler:      engine,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 120 * time.Second,
	}
	//server.SetKeepAlivesEnabled(false)

	go func(port string) {
		logger.Logger.Infof("Starting gin server, listen on %s", port)
		if err = server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Logger.Errorf("Start gin server error: %s", err)
		}
	}(cfg.Http.Port)

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	<-stopChan
	logger.Logger.Info("Shutting down server...")

	stopCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	if err = server.Shutdown(stopCtx); err != nil {
		logger.Logger.Errorf("Stop server error: %s", err)
		return
	}

	os.Exit(exitCode)
}

func newGin(mode string) *gin.Engine {
	gin.SetMode(mode)
	engine := gin.New()
	return engine
}

func handleRecover() {
	if r := recover(); r != nil {
		var msg string
		for i := 0; ; i++ {
			_, file, line, ok := runtime.Caller(i)
			if !ok {
				break
			}
			msg += fmt.Sprintf("%s:%d\n", file, line)
		}
		logger.Logger.Errorf("%s\n↧↧↧↧↧↧ PANIC ↧↧↧↧↧↧\n%s↥↥↥↥↥↥ PANIC ↥↥↥↥↥↥", r, msg)
	}
}
