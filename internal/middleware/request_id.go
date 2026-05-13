package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog"
)

func RequestID(c fiber.Ctx) error {
	requestID := c.Get("X-Request-ID")
	if requestID == "" {
		requestID = ulid.Make().String()
	}
	c.Set("X-Request-ID", requestID)
	c.Locals("request_id", requestID)

	return c.Next()
}

func WithAppLogger(base zerolog.Logger) fiber.Handler {
	return func(c fiber.Ctx) error {
		requestId, _ := c.Locals("request_id").(string)
		subLogger := base.With().Str("request_id", requestId).Logger()

		// Inject request-scoped logger into the Go context used downstream.
		ctx := subLogger.WithContext(c.Context())
		c.SetContext(ctx)

		return c.Next()
	}
}
