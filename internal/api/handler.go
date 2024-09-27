package api

import (
	"encoding/json"
	"go-fiber-api/internal/cache"
	"go-fiber-api/internal/config"
	"go-fiber-api/internal/db"
	"go-fiber-api/internal/logger"
	"go-fiber-api/internal/models"
	"sync"
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
        logger.Log.Info("Failed to set cache:", zap.Error(err))

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
    
    // 删除数据库中的用户
    if err := db.DB.Delete(&models.User{}, id).Error; err != nil {
        logger.Log.Error("Failed to delete user from database", zap.String("id", id), zap.Error(err))
        return fiber.NewError(fiber.StatusNotFound, "User not found")
    }

    // 删除缓存
    err := cache.DeleteCache(id)
    if err != nil {
        logger.Log.Warn("Failed to delete user from cache", zap.String("id", id), zap.Error(err))
        // 注意：我们不因为缓存删除失败而返回错误，因为用户已经从数据库中删除
    }

    logger.Log.Info("User deleted successfully", zap.String("id", id))
    return c.Status(fiber.StatusNoContent).JSON(nil)
}

// 使用并发处理批量请求
func GetMultipleUsers(c *fiber.Ctx) error {
    var userIDs models.UserIDs
    if err := c.BodyParser(&userIDs); err != nil {
        return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
    }
    var wg sync.WaitGroup
    userChan := make(chan *models.User, len(userIDs.IDS))
    errChan := make(chan error, len(userIDs.IDS))

    for _, id := range userIDs.IDS {
        wg.Add(1)
        go func(id string) {
            defer wg.Done()
            user, err := GetUserByID(id)
            if err != nil {
                errChan <- err
                return
            }
            userChan <- user
        }(string(id))
    }

    wg.Wait()
    close(userChan)
    close(errChan)

    users := make([]*models.User, 0, len(userIDs.IDS))
    for user := range userChan {
        users = append(users, user)
    }

    if len(errChan) > 0 {
        return fiber.NewError(fiber.StatusInternalServerError, "Error fetching some users")
    }

    return c.JSON(users)
}

func GetUserByID(id string) (*models.User, error) {
    // Try to get from cache first
    cachedUser, err := cache.GetCache(id)
    if err == nil {
        var user models.User
        if err := json.Unmarshal([]byte(cachedUser), &user); err == nil {
            return &user, nil
        }
    }

    // If not in cache, query the database
    var user models.User
    if err := db.DB.First(&user, id).Error; err != nil {
        return nil, err
    }

    // Cache the result
    cacheExpiration := time.Duration(config.Cfg.Redis.CacheExpirationMinutes) * time.Minute
    cache.SetCache(id, user, cacheExpiration)

    return &user, nil
}
