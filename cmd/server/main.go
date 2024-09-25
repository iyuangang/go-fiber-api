package main

import (
	"fmt"
	"go-fiber-api/internal/api"
	"go-fiber-api/internal/cache"
	"go-fiber-api/internal/config"
	"go-fiber-api/internal/db"
	"go-fiber-api/internal/middleware"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
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

    // 日志中间件
    app.Use(logger.New())

    // 设置路由
    app.Get("/user/:id", api.GetUser)
    app.Post("/user", api.CreateUser)
    app.Put("/user/:id", api.UpdateUser)
    app.Delete("/user/:id", api.DeleteUser)

    // 启动服务器
    log.Printf("Starting server on port %d", config.Cfg.Server.Port)
    if err := app.Listen(fmt.Sprintf(":%d", config.Cfg.Server.Port)); err != nil {
        log.Fatalf("Error starting server: %s", err)
    }
}
