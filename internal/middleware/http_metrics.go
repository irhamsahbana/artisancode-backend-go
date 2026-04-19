package middleware

import (
	"codebase-app/internal/infrastructure/metrics"
	"errors"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func WithHTTPMetrics(metricRegistry *metrics.Metrics) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Path() == "/metrics" {
			return c.Next()
		}

		startedAt := time.Now()
		method := c.Method()
		initialRoute := normalizedMetricRoute(c, normalizeMetricRoute(c.Path()))
		metricRegistry.AddInFlight(method, initialRoute, 1)
		defer metricRegistry.AddInFlight(method, initialRoute, -1)

		requestSizeBytes := approximateRequestSize(c)
		err := c.Next()

		route := normalizedMetricRoute(c, initialRoute)
		status := resolveMetricStatusCode(c, err)
		statusCode := strconv.Itoa(status)

		metricRegistry.ObserveRequest(
			method,
			route,
			statusCode,
			time.Since(startedAt).Seconds(),
			requestSizeBytes,
			len(c.Response().Body()),
		)

		return err
	}
}

func approximateRequestSize(c *fiber.Ctx) int {
	size := len(c.Method()) + len(c.Path()) + len(c.Context().QueryArgs().String()) + len(c.Request().Body())
	c.Request().Header.VisitAll(func(key []byte, value []byte) {
		size += len(key) + len(value)
	})

	return size
}

func normalizedMetricRoute(c *fiber.Ctx, fallback string) string {
	if c.Route() == nil {
		return fallback
	}

	return normalizeMetricRoute(c.Route().Path)
}

func normalizeMetricRoute(route string) string {
	if route == "" {
		return "unknown"
	}

	if route == "/" {
		return route
	}

	return route
}

func resolveMetricStatusCode(c *fiber.Ctx, err error) int {
	status := c.Response().StatusCode()
	if status >= fiber.StatusBadRequest {
		return status
	}

	if err == nil {
		if status == 0 {
			return fiber.StatusOK
		}

		return status
	}

	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		return fiberErr.Code
	}

	return fiber.StatusInternalServerError
}
