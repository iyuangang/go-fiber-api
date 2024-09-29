package main

import (
	"fmt"
	"go-fiber-api/internal/cache"
	"go-fiber-api/internal/config"
	"go-fiber-api/internal/db"
	"go-fiber-api/internal/logger"
	"go-fiber-api/internal/routes"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

var (
	Version string
	Commit  string
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
    // 创建 Fiber 应用实例，配置日志
	app := fiber.New(fiber.Config{
		// 启用日志
		EnablePrintRoutes: true,
		// 自定义日志输出
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			logger.Log.Error("Request error",
				zap.Int("status", code),
				zap.String("method", c.Method()),
				zap.String("path", c.Path()),
				zap.String("ip", c.IP()),
				zap.Duration("latency", c.Locals("latency").(time.Duration)),
				zap.Error(err),
			)
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})
    // 添加请求时间中间件
	app.Use(func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		latency := time.Since(start)
		c.Locals("latency", latency)
		
		// 记录成功的请求
		if err == nil {
			logger.Log.Info("Request processed",
				zap.String("method", c.Method()),
				zap.String("path", c.Path()),
				zap.Int("status", c.Response().StatusCode()),
				zap.Duration("latency", latency),
				zap.String("ip", c.IP()),
			)
		}
		return err
	})

    
    routes.SetupRoutes(app)


    // 启动服务器
		// 打印版本信息
    logger.Log.Info(fmt.Sprintf("Starting go-fiber-api version %s (commit %s)", Version, Commit),
		zap.Int("port", config.Cfg.Server.Port),
	  )

    if err := app.Listen(fmt.Sprintf(":%d", config.Cfg.Server.Port)); err != nil {
        logger.Log.Fatal("Error starting server", zap.Error(err))
    }
}
