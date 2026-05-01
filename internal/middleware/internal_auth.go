package middleware

import (
	"context"
	"strings"

	"codebase-app/internal/entity/common"
	"codebase-app/pkg/jwthandler"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func InternalAuth(c *fiber.Ctx) error {
	accessToken := c.Get("Authorization")
	unauthorizedResponse := fiber.Map{
		"message": "Unauthorized",
		"success": false,
	}

	if accessToken == "" {
		log.Error().Msg("Unauthorized internal auth - Header not set")
		return c.Status(fiber.StatusUnauthorized).JSON(unauthorizedResponse)
	}

	accessToken = strings.TrimPrefix(accessToken, "Bearer ")

	claims, err := jwthandler.ParseInternalUserTokenString(c.UserContext(), accessToken)
	if err != nil {
		log.Error().Err(err).Msg("Error while parsing internal auth token")
		return c.Status(fiber.StatusUnauthorized).JSON(unauthorizedResponse)
	}

	userCtx := common.UserContext{
		UserID:   claims.UserID,
		UserName: claims.UserName,
		Roles:    claims.Roles,
	}

	c.Locals("user_id", claims.UserID)
	c.Locals("roles", claims.Roles)
	c.Locals("user_name", claims.UserName)
	c.SetUserContext(context.WithValue(c.UserContext(), common.UserContextKeyClaims, userCtx))

	return c.Next()
}
