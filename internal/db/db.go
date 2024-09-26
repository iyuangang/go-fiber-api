package db

import (
	"go-fiber-api/internal/config"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
    var err error
    dsn := config.Cfg.Postgres.URL
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })

    if err != nil {
        log.Fatalf("Failed to connect to the database: %s", err)
    }

    sqlDB, err := DB.DB()
    if err != nil {
        log.Fatalf("Failed to get database instance: %s", err)
    }

    // 设置连接池参数
    sqlDB.SetMaxIdleConns(config.Cfg.Postgres.MaxIdleConns)
    sqlDB.SetMaxOpenConns(config.Cfg.Postgres.MaxOpenConns)
    sqlDB.SetConnMaxLifetime(time.Duration(config.Cfg.Postgres.ConnMaxLifetime) * time.Second)

    log.Println("Database connection initialized")
}
