package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

func CORS() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodPut,
			fiber.MethodDelete,
			fiber.MethodPatch,
			fiber.MethodOptions,
			fiber.MethodHead,
		},
		AllowHeaders: []string{
			fiber.HeaderOrigin,
			fiber.HeaderContentType,
			fiber.HeaderAccept,
			fiber.HeaderContentLength,
			fiber.HeaderAcceptLanguage,
			fiber.HeaderAcceptEncoding,
			fiber.HeaderConnection,
			fiber.HeaderAccessControlAllowOrigin,
			fiber.HeaderAuthorization,
		},
	})
}
