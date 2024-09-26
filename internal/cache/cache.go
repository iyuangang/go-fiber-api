package cache

import (
	"context"
	"encoding/json"
	"go-fiber-api/internal/config"
	"go-fiber-api/internal/logger"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

var RedisClient *redis.Client
var ctx = context.Background()

func InitRedis() {
    RedisClient = redis.NewClient(&redis.Options{
        Addr:     config.Cfg.Redis.Addr,
        Password: config.Cfg.Redis.Pass,
        DB:       config.Cfg.Redis.DB,
    })

    if err := RedisClient.Ping(ctx).Err(); err != nil {
        logger.Log.Error("Failed to connect to Redis:", zap.Error(err))
    }

    logger.Log.Info("Redis connection initialized") 
}

func GetCache(key string) (string, error) {
    return RedisClient.Get(ctx, key).Result()
}

func GetCacheObject(key string, obj interface{}) error {
    data, err := RedisClient.Get(ctx, key).Bytes()
    if err != nil {
        return err
    }
    
    return json.Unmarshal(data, obj)
}

func SetCache(key string, value interface{}, expiration time.Duration) error {
    jsonValue, err := json.Marshal(value)
    if err != nil {
        return err
    }
    
    // Store the JSON string in Redis
    return RedisClient.Set(ctx, key, jsonValue, expiration).Err()
}

func DeleteCache(key string) error {
    logger.Log.Debug("Deleting cache", zap.String("key", key))
    return RedisClient.Del(ctx, key).Err()
}
