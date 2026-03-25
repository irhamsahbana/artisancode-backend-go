package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.opentelemetry.io/otel"
)

func WithTracing(appName string) fiber.Handler {
	tracer := otel.Tracer(appName)
	return func(c *fiber.Ctx) error {
		ctx, span := tracer.Start(c.UserContext(), c.Path())
		defer span.End()
		c.SetUserContext(ctx)
		return c.Next()
	}
}

func Recover() fiber.Handler {
	return recover.New()
}
