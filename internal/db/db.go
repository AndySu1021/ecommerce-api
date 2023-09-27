package db

import (
	"database/sql"
	"ecommerce-api/internal/logger"
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
	"time"
)

type DatabaseType string

func (d DatabaseType) String() string {
	return string(d)
}

const (
	MySQL  DatabaseType = "mysql"
	SQLite DatabaseType = "sqlite"
)

type Config struct {
	Debug              bool         `mapstructure:"debug"`
	Type               DatabaseType `mapstructure:"type"`
	Host               string       `mapstructure:"host"`
	Port               int          `mapstructure:"port"`
	Username           string       `mapstructure:"username"`
	Password           string       `mapstructure:"password"`
	DBName             string       `mapstructure:"db_name"`
	MaxIdleConnections int          `mapstructure:"max_idle_connections"`
	MaxOpenConnections int          `mapstructure:"max_open_connections"`
	MaxLifetimeSec     int          `mapstructure:"max_lifetime_sec"`
}

func GetConnectionStr(cfg *Config) (connectionString string, err error) {
	switch cfg.Type {
	case MySQL:
		connectionString = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&multiStatements=true&parseTime=true", cfg.Username, cfg.Password, cfg.Host+":"+strconv.Itoa(cfg.Port), cfg.DBName)
	case SQLite:
		if cfg.Host == "" {
			connectionString = path.Join(os.Getenv("PROJ_DIR"), "test/.data", "sqlite.db?cache=shared")
		} else {
			connectionString = cfg.Host
		}
	default:
		return "", errors.New("not support driver")
	}

	return
}

func NewDatabase(cfg *Config) (db *sql.DB, err error) {
	dsn, err := GetConnectionStr(cfg)
	if err != nil {
		return nil, err
	}

	db, err = sql.Open(cfg.Type.String(), dsn)
	if err != nil {
		logger.Logger.Errorf("fail to open connection, err: %+v", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	logger.Logger.Infof("ping %s success", cfg.Type)

	if cfg.MaxIdleConnections != 0 {
		db.SetMaxIdleConns(cfg.MaxIdleConnections)
	} else {
		db.SetMaxIdleConns(25)
	}

	if cfg.MaxOpenConnections != 0 {
		db.SetMaxOpenConns(cfg.MaxOpenConnections)
	} else {
		db.SetMaxOpenConns(25)
	}

	if cfg.MaxLifetimeSec != 0 {
		db.SetConnMaxLifetime(time.Duration(cfg.MaxLifetimeSec) * time.Second)
	} else {
		db.SetConnMaxLifetime(5 * time.Minute)
	}

	return db, err
}
