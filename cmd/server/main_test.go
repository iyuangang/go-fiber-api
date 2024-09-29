package main

import (
	"fmt"
	"go-fiber-api/internal/cache"
	"go-fiber-api/internal/config"
	"go-fiber-api/internal/db"
	"go-fiber-api/internal/logger"
	"go-fiber-api/internal/routes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// 设置测试环境
	os.Setenv("APP_ENV", "test")
	
	// 运行测试
	code := m.Run()
	
	// 清理测试环境
	os.Exit(code)
}

func TestInitialization(t *testing.T) {
	// 测试初始化函数
	t.Run("InitLogger", testInitLogger)
	t.Run("InitConfig", testInitConfig)
	t.Run("InitDB", testInitDB)
	t.Run("InitRedis", testInitRedis)
}

func testInitLogger(t *testing.T) {
	logger.InitLogger()
	assert.NotNil(t, logger.Log)
}

func testInitConfig(t *testing.T) {
	config.InitConfig()
	assert.NotNil(t, config.Cfg)
	assert.Greater(t, config.Cfg.Server.Port, 0)
}

func testInitDB(t *testing.T) {
	db.InitDB()
	// 添加数据库连接检查
}

func testInitRedis(t *testing.T) {
	cache.InitRedis()
	// 添加 Redis 连接检查
}

func TestCreateApp(t *testing.T) {
	app := createApp()
	assert.NotNil(t, app)
	
	// 测试错误处理中间件
	t.Run("ErrorHandler", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/not-found", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})
	
	// 测试请求时间中间件
	t.Run("RequestTimeMiddleware", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		// 检查日志输出中是否包含请求时间
	})
}

func TestSetupRoutes(t *testing.T) {
	app := fiber.New()
	routes.SetupRoutes(app)
	
	// 测试各个路由是否正确设置
	// 例如：
	t.Run("HealthCheck", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/health", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})
	
	// 添加其他路由的测试...
}

func TestServerStart(t *testing.T) {
	// 模拟服务器启动
	app := createApp()
	go func() {
		err := startServer(app)
		assert.NoError(t, err)
	}()
	
	// 等待服务器启动
	time.Sleep(100 * time.Millisecond)
	
	// 测试服务器是否正常响应
	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/health", config.Cfg.Server.Port))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

