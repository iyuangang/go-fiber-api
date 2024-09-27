package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/helmet"
)

func SecurityMiddleware() fiber.Handler {
    return helmet.New()
}
