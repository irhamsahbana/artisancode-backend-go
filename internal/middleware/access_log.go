package middleware

import (
	"codebase-app/internal/infrastructure"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	oteltrace "go.opentelemetry.io/otel/trace"
)

func WithAccessLog(logger zerolog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		if err != nil {
			return err
		}

		// get request id from context
		requestId, _ := c.Context().UserValue("request_id").(string)
		spanCtx := oteltrace.SpanFromContext(c.UserContext()).SpanContext()

		event := infrastructure.AccessLogger.Info().Ctx(c.UserContext()).
			Str("method", c.Method()).
			Str("path", c.Path()).
			Any("query", c.Queries()).
			Str("ip", c.IP()).
			Str("user_agent", c.Get("User-Agent")).
			Dur("duration", time.Since(start)). // duration in ms
			Str("request_id", requestId)

		if spanCtx.HasTraceID() {
			event = event.
				Str("trace_id", spanCtx.TraceID().String()).
				Str("span_id", spanCtx.SpanID().String())
		}

		event.Msg("access log")
		return nil
	}
}
