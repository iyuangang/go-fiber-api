package db

import (
	"go-fiber-api/internal/config"
	"go-fiber-api/internal/logger"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

var DB *gorm.DB

func InitDB() {
    var err error
    dsn := config.Cfg.Postgres.Master
    dsn_slave := config.Cfg.Postgres.Slave
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
    })
    if err != nil {
        logger.Log.Fatal("Failed to connect to the database", zap.Error(err))
    }

    // 添加数据库读写分离
    err = DB.Use(dbresolver.Register(dbresolver.Config{
        Sources:  []gorm.Dialector{postgres.Open(dsn)},
        Replicas: []gorm.Dialector{postgres.Open(dsn_slave)},
        Policy:   dbresolver.RandomPolicy{},
    }))
    if err != nil {
        logger.Log.Fatal("Failed to configure database read/write splitting", zap.Error(err))
    }
    
    // 添加数据库连接池配置
    sqlDB, err := DB.DB()
    if err != nil {
        logger.Log.Fatal("Failed to get database connection", zap.Error(err))
    }

    // 设置连接池参数
    sqlDB.SetMaxIdleConns(config.Cfg.Postgres.MaxIdleConns)
    sqlDB.SetMaxOpenConns(config.Cfg.Postgres.MaxOpenConns)
    sqlDB.SetConnMaxLifetime(time.Duration(config.Cfg.Postgres.ConnMaxLifetime) * time.Second)

    logger.Log.Info("Database connection initialized")
}
