package monitor

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )
)

func init() {
    prometheus.MustRegister(httpRequestsTotal)
}

func PrometheusMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        err := c.Next()
        status := c.Response().StatusCode()
        httpRequestsTotal.WithLabelValues(c.Method(), c.Path(), strconv.Itoa(status)).Inc()
        return err
    }
}

func SetupPrometheusEndpoint(app *fiber.App) {
    app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))
}
