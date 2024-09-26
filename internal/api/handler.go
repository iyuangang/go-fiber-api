package api

import (
	"go-fiber-api/internal/cache"
	"go-fiber-api/internal/config"
	"go-fiber-api/internal/db"
	"go-fiber-api/internal/logger"
	"go-fiber-api/internal/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

//获取用户信息
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
        logger.Log.Info("Failed to set cache:", zap.String("id", id))

    }

    return c.Status(fiber.StatusOK).JSON(user)
}

// 创建用户
func CreateUser(c *fiber.Ctx) error {
    var user models.User
    if err := c.BodyParser(&user); err != nil {
        return fiber.NewError(fiber.StatusBadRequest, "Invalid input")
    }

    if err := db.DB.Create(&user).Error; err != nil {
        return fiber.NewError(fiber.StatusInternalServerError, "Failed to create user")
    }

    return c.Status(fiber.StatusCreated).JSON(user)
}

// 更新用户
func UpdateUser(c *fiber.Ctx) error {
    id := c.Params("id")
    var user models.User
    if err := c.BodyParser(&user); err != nil {
        return fiber.NewError(fiber.StatusBadRequest, "Invalid input")
    }

    if err := db.DB.Model(&user).Where("id = ?", id).Updates(user).Error; err != nil {
        return fiber.NewError(fiber.StatusInternalServerError, "Failed to update user")
    }

    // 更新缓存
    cache.SetCache(id, user, time.Duration(config.Cfg.Redis.CacheExpirationMinutes)*time.Minute)
    return c.Status(fiber.StatusOK).JSON(user)
}

// 删除用户
func DeleteUser(c *fiber.Ctx) error {
    id := c.Params("id")
    if err := db.DB.Delete(&models.User{}, id).Error; err != nil {
        return fiber.NewError(fiber.StatusNotFound, "User not found")
    }

    // 删除缓存
    // cache.RedisClient.Del(cache.ctx, id)

    return c.Status(fiber.StatusNoContent).JSON(nil)
}
