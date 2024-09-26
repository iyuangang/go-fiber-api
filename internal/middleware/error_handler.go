package middleware

import (
	"go-fiber-api/internal/logger"

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
    return c.Status(code).JSON(fiber.Map{
        "status":  code,
    })
}
