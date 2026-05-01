package middleware

import (
	"codebase-app/internal/infrastructure"
	"encoding/json"
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

		if c.Path() == "/metrics" {
			return nil
		}

		event := infrastructure.AccessLogger.Info().Ctx(c.UserContext()).
			Str("method", c.Method()).
			Str("path", c.Path()).
			Any("query", c.Queries()).
			Str("ip", c.IP()).
			Str("user_agent", c.Get("User-Agent")).
			Dur("duration", time.Since(start)). // duration in ms
			Str("request_id", requestId).
			Any("body", redactAccessLogBody(c.Body()))

		if spanCtx.HasTraceID() {
			event = event.
				Str("trace_id", spanCtx.TraceID().String()).
				Str("span_id", spanCtx.SpanID().String())
		}

		event.Msg("access log")
		return nil
	}
}

var redactedAccessLogKeys = map[string]struct{}{
	"password":           {},
	"token":              {},
	"id_token":           {},
	"access_token":       {},
	"refresh_token":      {},
	"registration_token": {},
	"authorization":      {},
	"authorization_code": {},
}

func redactAccessLogBody(body []byte) any {
	if len(body) == 0 {
		return ""
	}

	var payload any
	if err := json.Unmarshal(body, &payload); err != nil {
		return string(body)
	}

	return redactAccessLogValue(payload)
}

func redactAccessLogValue(value any) any {
	switch typed := value.(type) {
	case map[string]any:
		redacted := make(map[string]any, len(typed))
		for key, item := range typed {
			if _, sensitive := redactedAccessLogKeys[key]; sensitive {
				redacted[key] = "[redacted]"
				continue
			}
			redacted[key] = redactAccessLogValue(item)
		}
		return redacted
	case []any:
		redacted := make([]any, 0, len(typed))
		for _, item := range typed {
			redacted = append(redacted, redactAccessLogValue(item))
		}
		return redacted
	default:
		return value
	}
}
