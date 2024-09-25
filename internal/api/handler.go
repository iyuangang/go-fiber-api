package api

import (
    "github.com/gofiber/fiber/v2"
    "go-fiber-api/internal/cache"
    "go-fiber-api/internal/db"
    "go-fiber-api/internal/models"
    "go-fiber-api/internal/config"
    "time"
    "log"
)

func GetUser(c *fiber.Ctx) error {
    id := c.Params("id")

    // 从 Redis 获取缓存
    cachedUser, err := cache.GetCache(id)
    if err == nil {
        return c.Status(fiber.StatusOK).SendString(cachedUser)
    }

    // 查询 PostgreSQL
    var user models.User
    if err := db.DB.First(&user, id).Error; err != nil {
        return fiber.NewError(fiber.StatusNotFound, "User not found")
    }

    // 将查询结果存入 Redis
    cacheExpiration := time.Duration(config.Cfg.Redis.CacheExpirationMinutes) * time.Minute
    if err := cache.SetCache(id, user, cacheExpiration); err != nil {
        log.Println("Failed to set cache:", err)
    }

    return c.Status(fiber.StatusOK).JSON(user)
}
