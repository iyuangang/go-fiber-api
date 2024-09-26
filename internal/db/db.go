package db

import (
	"go-fiber-api/internal/config"
	"time"

	"go-fiber-api/internal/logger"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
    var err error
    dsn := config.Cfg.Postgres.URL
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
    })

    if err != nil {
        logger.Log.Fatal("Failed to connect to the database", zap.Error(err))

    }

    sqlDB, err := DB.DB()
    if err != nil {
        logger.Log.Fatal("Failed to get database instance", zap.Error(err))

    }

    // 设置连接池参数
    sqlDB.SetMaxIdleConns(config.Cfg.Postgres.MaxIdleConns)
    sqlDB.SetMaxOpenConns(config.Cfg.Postgres.MaxOpenConns)
    sqlDB.SetConnMaxLifetime(time.Duration(config.Cfg.Postgres.ConnMaxLifetime) * time.Second)

    logger.Log.Info("Database connection initialized")
}
