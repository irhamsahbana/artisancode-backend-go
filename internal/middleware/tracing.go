package middleware

import (
	"codebase-app/internal/infrastructure/tracing"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

func WithTracing(appName string) fiber.Handler {
	return func(c fiber.Ctx) error {
		header := make(propagation.HeaderCarrier)
		for k, v := range c.GetReqHeaders() {
			header[k] = v
		}

		ctx := otel.GetTextMapPropagator().Extract(c.Context(), header)
		spanName := c.Method() + " " + c.Path()
		ctx, span := tracing.StartSpan(ctx, spanName,
			oteltrace.WithSpanKind(oteltrace.SpanKindServer),
			oteltrace.WithAttributes(
				semconv.HTTPRequestMethodKey.String(c.Method()),
				semconv.URLPath(c.Path()),
				semconv.ClientAddress(c.IP()),
				attribute.String("user_agent.original", c.Get("User-Agent")),
			),
		)
		defer span.End()

		c.SetContext(ctx)

		if sc := span.SpanContext(); sc.HasTraceID() {
			c.Set("X-Trace-ID", sc.TraceID().String())
		}

		err := c.Next()

		routePath := c.Route().Path
		if routePath != "" && routePath != "/" {
			span.SetName(c.Method() + " " + routePath)
			span.SetAttributes(semconv.HTTPRoute(routePath))
		}

		status := c.Response().StatusCode()
		span.SetAttributes(semconv.HTTPResponseStatusCode(status))

		if err != nil {
			tracing.RecordError(span, err)
		} else if status >= fiber.StatusInternalServerError {
			span.SetStatus(codes.Error, "internal server error")
		}

		return err
	}
}

func Recover() fiber.Handler {
	return recover.New()
}
