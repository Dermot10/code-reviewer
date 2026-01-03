package database

import (
	"context"
	"fmt"
	"time"

	"github.com/dermot10/code-reviewer/backend_go/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(ctx context.Context, cfg *config.Config) (*gorm.DB, error) {
	dsn := cfg.DatabaseURL
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL not set")
	}

	var db *gorm.DB
	var err error

	for i := 0; i < 5; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Error),
		})
		if err == nil {
			break
		}
		time.Sleep(time.Second * time.Duration(i+1))
	}
	if err != nil {
		return nil, fmt.Errorf("failed to connect after retries: %w", err)
	}

	// pool config
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql DB: %w", err)
	}
	sqlDB.SetMaxOpenConns(cfg.DBMaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.DBMaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Minute * 30)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
