package handler

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/mapper"
	"codebase-app/internal/entity/restentity"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (h *userHandler) refreshToken(c *fiber.Ctx) error {
	var (
		ctx = c.UserContext()
		req = new(restentity.RefreshTokenReq)
		v   = adapter.Adapters.Validator
	)

	if err := c.BodyParser(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Any(common.LogKeyPayload, req).Msg("Invalid request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := v.Validate(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Any(common.LogKeyPayload, req).Msg("Invalid request validation")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	user := mapper.RefreshTokenReqToCore(ctx, *req)
	tokens, err := h.core.RefreshToken(ctx, user)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("RefreshToken service error")
		code, errs := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	resp := mapper.AuthTokensToLoginResp(*tokens)
	return c.Status(fiber.StatusOK).JSON(response.Success(resp, ""))
}
