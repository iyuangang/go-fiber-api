package main

import (
	"fmt"
	"go-fiber-api/internal/api"
	"go-fiber-api/internal/cache"
	"go-fiber-api/internal/config"
	"go-fiber-api/internal/db"
	"go-fiber-api/internal/logger"
	"go-fiber-api/internal/middleware"
	"go-fiber-api/internal/monitor"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func main() {
    // 初始化日志
    logger.InitLogger()
    defer logger.Log.Sync()

    // 初始化配置
    config.InitConfig()

    // 初始化数据库
    db.InitDB()

    // 初始化 Redis
    cache.InitRedis()

    // 初始化 Fiber
    app := fiber.New(fiber.Config{
        ErrorHandler: middleware.ErrorHandler,
        ReadTimeout:  time.Duration(config.Cfg.Server.ReadTimeout) * time.Second,
    })
    app.Use(middleware.SecurityMiddleware())
    app.Use(monitor.PrometheusMiddleware())

    // 设置路由
    app.Get("/user/:id", api.GetUser)
    app.Get("/users/:id", api.GetMultipleUsers)
    app.Post("/user", api.CreateUser)
    app.Put("/user/:id", api.UpdateUser)
    app.Delete("/user/:id", api.DeleteUser)

    monitor.SetupPrometheusEndpoint(app)

    // 启动服务器
    logger.Log.Info("Starting server",
        zap.Int("port", config.Cfg.Server.Port),
    )

    if err := app.Listen(fmt.Sprintf(":%d", config.Cfg.Server.Port)); err != nil {
        logger.Log.Fatal("Error starting server", zap.Error(err))
    }
}
