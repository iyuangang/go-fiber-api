package routes

import (
	"go-fiber-api/internal/api"
	"go-fiber-api/internal/middleware"
	"go-fiber-api/internal/monitor"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func SetupRoutes(app *fiber.App) {
	// 用户相关路由
	userGroup := app.Group("/user")
	userGroup.Get("/:id", api.GetUser)
	userGroup.Post("/", api.CreateUser)
	userGroup.Put("/:id", api.UpdateUser)
	userGroup.Delete("/:id", api.DeleteUser)

	// 批量获取用户
	app.Get("/users", api.GetMultipleUsers)

	// 健康检查路由
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// 添加中间件
	app.Use(middleware.SecurityMiddleware())
  app.Use(recover.New())
  app.Use(middleware.SecurityMiddleware())
  app.Use(monitor.PrometheusMiddleware())

	// 可以在这里添加更多的路由和中间件
}

