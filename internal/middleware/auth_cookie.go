package middleware

import (
	"codebase-app/internal/entity/common"
	"codebase-app/pkg/jwthandler"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func AuthMiddleware(c *fiber.Ctx) error {
	cookie := c.Cookies("access_token")

	if cookie == "" {
		log.Error().Msg("Unauthorized - Cookie not set")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
			"success": false,
		})
	}

	claims, err := jwthandler.ParseTokenString(c.UserContext(), cookie)
	if err != nil {
		log.Error().Err(err).Msg("Error while parsing token")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bad request",
			"success": false,
		})
	}

	c.Locals("user_id", claims.UserID)
	c.Locals("roles", claims.Roles)
	c.Locals("tenant_id", claims.TenantID)
	c.Locals("tenant_name", claims.TenantName)
	c.Locals("user_name", claims.UserName)
	c.Locals("company_id", claims.CompanyID)
	c.Locals("company_name", claims.CompanyName)
	userCtx := common.UserContext{
		UserID:      claims.UserID,
		UserName:    claims.UserName,
		Roles:       claims.Roles,
		TenantID:    claims.TenantID,
		TenantName:  claims.TenantName,
		CompanyID:   claims.CompanyID,
		CompanyName: claims.CompanyName,
	}
	c.Context().SetUserValue(common.UserContextKeyClaims, userCtx)

	return c.Next()
}
