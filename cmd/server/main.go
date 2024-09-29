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
	initializeApp()
	app := createApp()
	routes.SetupRoutes(app)
	startServer(app)
}

func initializeApp() {
	logger.InitLogger()
	config.InitConfig()
	db.InitDB()
	cache.InitRedis()
}

func createApp() *fiber.App {
	app := fiber.New(fiber.Config{
		EnablePrintRoutes: true,
		ErrorHandler: customErrorHandler,
	})
	app.Use(requestTimeMiddleware)
	return app
}

func customErrorHandler(c *fiber.Ctx, err error) error {
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
}

func requestTimeMiddleware(c *fiber.Ctx) error {
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
}

func startServer(app *fiber.App) error {
	logger.Log.Info(fmt.Sprintf("Starting go-fiber-api version %s (commit %s)", Version, Commit),
		zap.Int("port", config.Cfg.Server.Port),
	)
	return app.Listen(fmt.Sprintf(":%d", config.Cfg.Server.Port))
}
