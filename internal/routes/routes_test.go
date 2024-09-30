package routes

import (
	"go-fiber-api/internal/api"
	"go-fiber-api/internal/middleware"
	"go-fiber-api/internal/monitor"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAPI 是一个模拟的 API 实现
type MockAPI struct {
	mock.Mock
}

func (m *MockAPI) GetUser(c *fiber.Ctx) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m *MockAPI) CreateUser(c *fiber.Ctx) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m *MockAPI) UpdateUser(c *fiber.Ctx) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m *MockAPI) DeleteUser(c *fiber.Ctx) error {
	args := m.Called(c)
	return args.Error(0)
}

func (m *MockAPI) GetMultipleUsers(c *fiber.Ctx) error {
	args := m.Called(c)
	return args.Error(0)
}

// TestSetupRoutes 测试路由设置
func TestSetupRoutes(t *testing.T) {
	app := fiber.New()
	mockAPI := new(MockAPI)

	// 替换真实的 API 处理函数为 mock
	api.GetUser = mockAPI.GetUser
	api.CreateUser = mockAPI.CreateUser
	api.UpdateUser = mockAPI.UpdateUser
	api.DeleteUser = mockAPI.DeleteUser
	api.GetMultipleUsers = mockAPI.GetMultipleUsers

	SetupRoutes(app)

	// 测试中间件
	t.Run("Middleware", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	// 测试用户路由
	t.Run("UserRoutes", func(t *testing.T) {
		testCases := []struct {
			name           string
			method         string
			path           string
			expectedStatus int
			mockSetup      func()
		}{
			{
				name:           "GetUser",
				method:         "GET",
				path:           "/user/1",
				expectedStatus: fiber.StatusOK,
				mockSetup: func() {
					mockAPI.On("GetUser", mock.Anything).Return(nil)
				},
			},
			{
				name:           "CreateUser",
				method:         "POST",
				path:           "/user",
				expectedStatus: fiber.StatusOK,
				mockSetup: func() {
					mockAPI.On("CreateUser", mock.Anything).Return(nil)
				},
			},
			{
				name:           "UpdateUser",
				method:         "PUT",
				path:           "/user/1",
				expectedStatus: fiber.StatusOK,
				mockSetup: func() {
					mockAPI.On("UpdateUser", mock.Anything).Return(nil)
				},
			},
			{
				name:           "DeleteUser",
				method:         "DELETE",
				path:           "/user/1",
				expectedStatus: fiber.StatusOK,
				mockSetup: func() {
					mockAPI.On("DeleteUser", mock.Anything).Return(nil)
				},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				tc.mockSetup()
				req := httptest.NewRequest(tc.method, tc.path, nil)
				resp, err := app.Test(req)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedStatus, resp.StatusCode)
			})
		}
	})

	// 测试批量获取用户
	t.Run("GetMultipleUsers", func(t *testing.T) {
		mockAPI.On("GetMultipleUsers", mock.Anything).Return(nil)
		req := httptest.NewRequest("GET", "/users", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	// 测试健康检查路由
	t.Run("HealthCheck", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/health", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		assert.Equal(t, "OK", string(body))
	})

	// 测试 Prometheus 端点
	t.Run("PrometheusEndpoint", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/metrics", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})
}

// TestSecurityMiddleware 测试安全中间件
func TestSecurityMiddleware(t *testing.T) {
	app := fiber.New()
	app.Use(middleware.SecurityMiddleware())
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	req := httptest.NewRequest("GET", "/", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// 检查安全头部
	assert.Equal(t, "nosniff", resp.Header.Get("X-Content-Type-Options"))
	assert.Equal(t, "SAMEORIGIN", resp.Header.Get("X-Frame-Options"))
	assert.Equal(t, "1; mode=block", resp.Header.Get("X-XSS-Protection"))
}

// TestPrometheusMiddleware 测试 Prometheus 中间件
func TestPrometheusMiddleware(t *testing.T) {
	app := fiber.New()
	app.Use(monitor.PrometheusMiddleware())
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	req := httptest.NewRequest("GET", "/", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	// 检查 Prometheus 指标
	req = httptest.NewRequest("GET", "/metrics", nil)
	resp, err = app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "http_requests_total")
}
