package cache

import (
	"context"
	"go-fiber-api/internal/config"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
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
        log.Fatalf("Failed to connect to Redis: %s", err)
    }

    log.Println("Redis connection initialized")
}

func GetCache(key string) (string, error) {
    return RedisClient.Get(ctx, key).Result()
}

func SetCache(key string, value interface{}, expiration time.Duration) error {
    return RedisClient.Set(ctx, key, value, expiration).Err()
}
