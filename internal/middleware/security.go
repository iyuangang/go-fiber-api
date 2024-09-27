package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/helmet"
)

func SecurityMiddleware() fiber.Handler {
    return helmet.New()
}

// 在 main.go 中添加
// app.Use(middleware.SecurityMiddleware())
