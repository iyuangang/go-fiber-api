package middleware

import (
	"go-fiber-api/internal/config"
	"go-fiber-api/internal/logger"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	logger.Log.Error("Request error",
		zap.Int("status", code),
		zap.String("path", c.Path()),
		zap.String("ip", c.IP()),
		zap.Error(err),
	)
	
    // 添加更详细的错误响应
    errorResponse := fiber.Map{
			"status":  code,
			"message": err.Error(),
	}

	// 在非生产环境中，可以添加更多调试信息
	if !config.Cfg.Env.IsProduction {
			errorResponse["stack"] = string(debug.Stack())
	}

	return c.Status(code).JSON(errorResponse)
}
