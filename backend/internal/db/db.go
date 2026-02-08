package db

import (
	"fmt"

	"github.com/khorzhenwin/go-cafe/backend/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func NewAWSClient(cfg *config.DBConfig) (*gorm.DB, error) {
	postgresCfg := postgres.Config{
		DSN:                  cfg.GetFormattedDSN(),
		PreferSimpleProtocol: true,
	}

	gormCfg := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{TablePrefix: "gocafe_"},
	}

	db, err := gorm.Open(postgres.New(postgresCfg), gormCfg)
	if err != nil {
		return nil, fmt.Errorf("DB connection failed: %w", err)
	}
	return db, nil
}

func NewLocalDbClient(cfg *config.DBConfig) (*gorm.DB, error) {
	postgresCfg := postgres.Config{
		DSN:                  cfg.GetLocalDSN(),
		PreferSimpleProtocol: true,
	}

	gormCfg := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{TablePrefix: "gocafe_"},
	}

	db, err := gorm.Open(postgres.New(postgresCfg), gormCfg)
	if err != nil {
		return nil, fmt.Errorf("DB connection failed: %w", err)
	}
	return db, nil
}
