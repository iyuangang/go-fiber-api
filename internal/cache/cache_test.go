package cache

import (
	"context"
	"encoding/json"
	"errors"
	"go-fiber-api/internal/config"
	"go-fiber-api/internal/logger"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRedis(t *testing.T) (*miniredis.Miniredis, func()) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	config.Cfg = config.Config{
		Redis: config.RedisConfig{
			Addr: mr.Addr(),
			Pass: "",
			DB:   0,
		},
	}

	return mr, func() {
		mr.Close()
	}
}

func TestInitRedis(t *testing.T) {
	mr, cleanup := setupTestRedis(t)
	defer cleanup()

	logger.InitLogger()

	t.Run("Successful connection", func(t *testing.T) {
		InitRedis()
		assert.NotNil(t, RedisClient)
		assert.NoError(t, RedisClient.Ping(context.Background()).Err())
	})

	t.Run("Failed connection", func(t *testing.T) {
		mr.Close()
		InitRedis()
		assert.NotNil(t, RedisClient)
		assert.Error(t, RedisClient.Ping(context.Background()).Err())
	})
}

func TestGetCache(t *testing.T) {
	_, cleanup := setupTestRedis(t)
	defer cleanup()

	InitRedis()

	t.Run("Existing key", func(t *testing.T) {
		key := "test_key"
		value := "test_value"
		err := RedisClient.Set(context.Background(), key, value, 0).Err()
		require.NoError(t, err)

		result, err := GetCache(key)
		assert.NoError(t, err)
		assert.Equal(t, value, result)
	})

	t.Run("Non-existing key", func(t *testing.T) {
		key := "non_existing_key"
		result, err := GetCache(key)
		assert.Error(t, err)
		assert.Equal(t, redis.Nil, err)
		assert.Empty(t, result)
	})
}

func TestGetCacheObject(t *testing.T) {
	_, cleanup := setupTestRedis(t)
	defer cleanup()

	InitRedis()

	type TestStruct struct {
		Name string
		Age  int
	}

	t.Run("Existing object", func(t *testing.T) {
		key := "test_object"
		value := TestStruct{Name: "John", Age: 30}
		jsonValue, _ := json.Marshal(value)
		err := RedisClient.Set(context.Background(), key, jsonValue, 0).Err()
		require.NoError(t, err)

		var result TestStruct
		err = GetCacheObject(key, &result)
		assert.NoError(t, err)
		assert.Equal(t, value, result)
	})

	t.Run("Non-existing object", func(t *testing.T) {
		key := "non_existing_object"
		var result TestStruct
		err := GetCacheObject(key, &result)
		assert.Error(t, err)
		assert.Equal(t, redis.Nil, err)
	})
}

func TestSetCache(t *testing.T) {
	_, cleanup := setupTestRedis(t)
	defer cleanup()

	InitRedis()

	t.Run("Set string value", func(t *testing.T) {
		key := "test_set_key"
		value := "test_set_value"
		err := SetCache(key, value, time.Minute)
		assert.NoError(t, err)

		result, err := RedisClient.Get(context.Background(), key).Result()
		assert.NoError(t, err)
		assert.Equal(t, value, result)
	})

	t.Run("Set object value", func(t *testing.T) {
		key := "test_set_object"
		value := struct {
			Name string
			Age  int
		}{
			Name: "Alice",
			Age:  25,
		}
		err := SetCache(key, value, time.Minute)
		assert.NoError(t, err)

		result, err := RedisClient.Get(context.Background(), key).Result()
		assert.NoError(t, err)

		var decodedValue struct {
			Name string
			Age  int
		}
		err = json.Unmarshal([]byte(result), &decodedValue)
		assert.NoError(t, err)
		assert.Equal(t, value, decodedValue)
	})
}

func TestDeleteCache(t *testing.T) {
	_, cleanup := setupTestRedis(t)
	defer cleanup()

	InitRedis()

	t.Run("Delete existing key", func(t *testing.T) {
		key := "test_delete_key"
		value := "test_delete_value"
		err := RedisClient.Set(context.Background(), key, value, 0).Err()
		require.NoError(t, err)

		err = DeleteCache(key)
		assert.NoError(t, err)

		_, err = RedisClient.Get(context.Background(), key).Result()
		assert.Error(t, err)
		assert.Equal(t, redis.Nil, err)
	})

	t.Run("Delete non-existing key", func(t *testing.T) {
		key := "non_existing_delete_key"
		err := DeleteCache(key)
		assert.NoError(t, err)
	})
}

func TestWarmUpCache(t *testing.T) {
	_, cleanup := setupTestRedis(t)
	defer cleanup()

	InitRedis()

	err := WarmUpCache()
	assert.NoError(t, err)
	// Add more specific assertions based on your implementation
}

func TestGetCacheWithFallback(t *testing.T) {
	_, cleanup := setupTestRedis(t)
	defer cleanup()

	InitRedis()

	t.Run("Cache hit", func(t *testing.T) {
		key := "test_fallback_key"
		value := "test_fallback_value"
		err := RedisClient.Set(context.Background(), key, value, time.Minute).Err()
		require.NoError(t, err)

		result, err := GetCacheWithFallback(key, func() (interface{}, error) {
			return "fallback_value", nil
		}, time.Minute)

		assert.NoError(t, err)
		assert.Equal(t, value, result)
	})

	t.Run("Cache miss with successful fallback", func(t *testing.T) {
		key := "test_fallback_miss_key"
		fallbackValue := "fallback_value"

		result, err := GetCacheWithFallback(key, func() (interface{}, error) {
			return fallbackValue, nil
		}, time.Minute)

		assert.NoError(t, err)
		assert.Equal(t, fallbackValue, result)

		// Verify that the fallback value was cached
		cachedValue, err := RedisClient.Get(context.Background(), key).Result()
		assert.NoError(t, err)
		assert.Equal(t, fallbackValue, cachedValue)
	})

	t.Run("Cache miss with fallback error", func(t *testing.T) {
		key := "test_fallback_error_key"
		expectedError := errors.New("fallback error")

		result, err := GetCacheWithFallback(key, func() (interface{}, error) {
			return nil, expectedError
		}, time.Minute)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, result)
	})
}
