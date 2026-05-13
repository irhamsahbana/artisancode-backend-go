package setup

import (
	"codebase-app/internal/infrastructure"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"

	"github.com/gofiber/fiber/v3"
)

func (deps httpDependencies) registerRouteNotFound() {
	deps.app.Use(func(c fiber.Ctx) error {
		var (
			method = c.Method()
			path   = c.Path()
			query  = c.RequestCtx().QueryArgs().String()
			ua     = c.Get("User-Agent")
			ip     = c.IP()
		)

		infrastructure.AccessLogger.Debug().
			Str("method", method).
			Str("path", path).
			Str("query", query).
			Str("ua", ua).
			Str("ip", ip).
			Msg("route not found")

		return c.Status(fiber.StatusNotFound).JSON(
			response.Error(errmsg.MessageRouteNotFound),
		)
	})
}
