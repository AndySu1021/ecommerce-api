package redis

import (
	"context"
	"ecommerce-api/internal/logger"
	"fmt"
	"github.com/cenkalti/backoff/v4"
	"github.com/go-redis/redis/v8"
	"runtime"
	"time"
)

type Config struct {
	MasterName string   `mapstructure:"master_name"`
	Addresses  []string `mapstructure:"addresses"`
	Password   string   `mapstructure:"password"`
	DB         int      `mapstructure:"db"`
}

func NewRedisClient(cfg Config) (redis.UniversalClient, error) {
	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = time.Duration(180) * time.Second

	if len(cfg.Addresses) == 0 {
		return nil, fmt.Errorf("redis config address is empty")
	}

	var client redis.UniversalClient

	if err := backoff.Retry(func() error {
		client = redis.NewUniversalClient(&redis.UniversalOptions{
			Addrs:              cfg.Addresses,
			DB:                 cfg.DB,
			Dialer:             nil,
			OnConnect:          nil,
			Username:           "",
			Password:           cfg.Password,
			SentinelUsername:   "",
			SentinelPassword:   "",
			MaxRetries:         5,
			MinRetryBackoff:    8 * time.Millisecond,
			MaxRetryBackoff:    512 * time.Millisecond,
			DialTimeout:        5 * time.Second,
			ReadTimeout:        3 * time.Second,
			WriteTimeout:       3 * time.Second,
			PoolFIFO:           false,
			PoolSize:           4 * runtime.NumCPU(),
			MinIdleConns:       10,
			MaxConnAge:         0,
			PoolTimeout:        4 * time.Second,
			IdleTimeout:        5 * time.Minute,
			IdleCheckFrequency: 0,
			MaxRedirects:       0,
			ReadOnly:           true,
			RouteByLatency:     true,
			RouteRandomly:      true,
			MasterName:         cfg.MasterName,
		})

		if err := client.Ping(context.Background()).Err(); err != nil {
			logger.Logger.Errorf("ping occurs error after connecting to redis: %s", err)
			return fmt.Errorf("ping occurs error after connecting to redis: %s", err)
		}
		logger.Logger.Info("ping redis success")
		return nil
	}, bo); err != nil {
		return nil, err
	}

	return client, nil
}
