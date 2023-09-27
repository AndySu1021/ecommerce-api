package storage

import (
	"ecommerce-api/internal/storage/driver"
	"errors"
)

type Config struct {
	Driver  Driver `mapstructure:"driver"`
	BaseUrl string `mapstructure:"base_url"`
	GCS     GCS    `mapstructure:"gcs"`
	S3      S3     `mapstructure:"S3"`
}

type GCS struct {
	Bucket         string `mapstructure:"bucket"`
	CredentialPath string `mapstructure:"credential_path"`
}

type S3 struct {
	Key      string `mapstructure:"key"`
	Secret   string `mapstructure:"secret"`
	Region   string `mapstructure:"region"`
	Bucket   string `mapstructure:"bucket"`
	Endpoint string `mapstructure:"endpoint"`
}

func NewStorage(cfg Config) (driver.Client, error) {
	switch cfg.Driver {
	case DriverLocal:
		return driver.NewLocal(), nil
	case DriverS3:
		return driver.NewS3(cfg.S3.Key, cfg.S3.Secret, cfg.S3.Region, cfg.S3.Bucket, cfg.S3.Endpoint), nil
	case DriverGCS:
		return driver.NewGCS(cfg.GCS.Bucket, cfg.GCS.CredentialPath), nil
	}

	return nil, errors.New("driver not support")
}
