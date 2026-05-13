package middleware

import (
	"codebase-app/internal/infrastructure/metrics"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/require"
)

func TestWithHTTPMetricsExposesPrometheusMetrics(t *testing.T) {
	app := fiber.New()
	metricRegistry := metrics.New("artisan-backend", "1.0.0", "test")

	app.Use(WithHTTPMetrics(metricRegistry))
	app.Get("/metrics", metricRegistry.Handler())
	app.Get("/companies/:id", func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusCreated)
	})
	app.Get("/fail", func(c fiber.Ctx) error {
		return fiber.NewError(fiber.StatusBadRequest, "bad request")
	})

	resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/companies/123", nil))
	require.NoError(t, err)
	require.Equal(t, fiber.StatusCreated, resp.StatusCode)

	resp, err = app.Test(httptest.NewRequest(http.MethodGet, "/fail", nil))
	require.NoError(t, err)
	require.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	resp, err = app.Test(httptest.NewRequest(http.MethodGet, "/metrics", nil))
	require.NoError(t, err)
	require.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	metricsText := string(body)
	require.Contains(t, metricsText, "artisancode_backend_http_request_duration_seconds_bucket")
	require.Contains(t, metricsText, "artisancode_backend_http_requests_total")
	require.Contains(t, metricsText, "route=\"/companies/:id\"")
	require.Contains(t, metricsText, "status=\"201\"")
	require.Contains(t, metricsText, "route=\"/fail\"")
	require.Contains(t, metricsText, "status=\"400\"")
	require.False(t, strings.Contains(metricsText, "route=\"/metrics\""))
}
