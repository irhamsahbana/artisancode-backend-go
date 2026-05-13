package metrics

import (
	"github.com/gofiber/fiber/v3"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

func (m *Metrics) Handler() fiber.Handler {
	handler := fasthttpadaptor.NewFastHTTPHandler(
		promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{}),
	)

	return func(c fiber.Ctx) error {
		handler(c.RequestCtx())
		return nil
	}
}
