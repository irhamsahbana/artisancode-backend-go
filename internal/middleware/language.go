package middleware

import (
	"codebase-app/pkg/errmsg"

	"github.com/gofiber/fiber/v2"
)

func WithRequestLanguage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		language := string(errmsg.ResolveLanguage(c.Get("Accept-Language")))
		ctx := errmsg.ContextWithLanguage(c.UserContext(), language)
		c.SetUserContext(ctx)

		return c.Next()
	}
}
