package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
	"time"
)

var Logger *zap.SugaredLogger

type Config struct {
	AppName     string        `mapstructure:"app_name"`
	Environment string        `mapstructure:"environment"` // local, development, staging, production
	Level       zapcore.Level `mapstructure:"level"`       // debug: -1, info: 0, ...
	Directory   string        `mapstructure:"directory"`   // 檔案儲存目錄
}

func InitZapLogger(cfg Config) error {
	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())

	syncSlice := []zapcore.WriteSyncer{zapcore.AddSync(os.Stdout)}
	if cfg.Directory != "" {
		f, err := resolveFile(cfg.Directory, cfg.AppName)
		if err != nil {
			return err
		}
		syncSlice = append(syncSlice, zapcore.AddSync(f))
	}

	syncer := zapcore.NewMultiWriteSyncer(syncSlice...)
	core := zapcore.NewCore(encoder, syncer, zap.NewAtomicLevelAt(cfg.Level))

	Logger = zap.New(core).
		With(zap.String("app", cfg.AppName)).
		With(zap.String("env", cfg.Environment)).
		Sugar()

	return nil
}

func resolveFile(dir, app string) (*os.File, error) {
	dir = strings.TrimRight(dir, "/")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err = os.MkdirAll(dir, 0700); err != nil {
			return nil, err
		}
	}
	filename := fmt.Sprintf("%s_%s.log", app, time.Now().Format("20060102"))
	return os.OpenFile(dir+"/"+filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
}
