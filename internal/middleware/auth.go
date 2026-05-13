package middleware

import (
	"context"
	"strings"

	"codebase-app/internal/entity/common"
	"codebase-app/pkg/jwthandler"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
)

func Auth(c fiber.Ctx) error {
	AccessToken := c.Get("Authorization")
	unauthorizedResponse := fiber.Map{
		"message": "Unauthorized",
		"success": false,
	}

	if AccessToken == "" {
		log.Error().Msg("Unauthorized - Header not set")
		return c.Status(fiber.StatusUnauthorized).JSON(unauthorizedResponse)
	}

	AccessToken = strings.TrimPrefix(AccessToken, "Bearer ")

	claims, err := jwthandler.ParseTokenString(c.Context(), AccessToken)
	if err != nil {
		log.Error().Err(err).Msg("Error while parsing token")
		return c.Status(fiber.StatusUnauthorized).JSON(unauthorizedResponse)
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
		TenantID:    claims.TenantID,
		TenantName:  claims.TenantName,
		Roles:       claims.Roles,
		CompanyID:   claims.CompanyID,
		CompanyName: claims.CompanyName,
	}

	c.SetContext(context.WithValue(c.Context(), common.UserContextKeyClaims, userCtx))

	return c.Next()
}
