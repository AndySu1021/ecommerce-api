package config

import (
	"ecommerce-api/internal/db"
	"ecommerce-api/internal/logger"
	"ecommerce-api/internal/redis"
	"ecommerce-api/internal/storage"
	"github.com/AndySu1021/go-util/gin"
	"github.com/spf13/viper"
)

type AppConfig struct {
	Name          string `mapstructure:"name"`
	Env           string `mapstructure:"env"`
	MigrationPath string `mapstructure:"migration_path"`
	TemplatePath  string `mapstructure:"template_path"`
	ServerID      int64
}

type TapPayConfig struct {
	IPList []string `mapstructure:"ip_list"`
}

type MailgunConfig struct {
	Sender string `mapstructure:"sender"`
	Domain string `mapstructure:"domain"`
	ApiKey string `mapstructure:"api_key"`
}

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Storage  storage.Config `mapstructure:"storage"`
	Log      logger.Config  `mapstructure:"log"`
	Http     gin.Config     `mapstructure:"http"`
	Database db.Config      `mapstructure:"database"`
	Redis    redis.Config   `mapstructure:"redis"`
	TapPay   TapPayConfig   `mapstructure:"tappay"`
	Mailgun  MailgunConfig  `mapstructure:"mailgun"`
}

func NewConfig(path string) (*Config, error) {
	var config *Config

	if path == "" {
		path = "./"
	}

	viper.AutomaticEnv()
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return config, nil
}
