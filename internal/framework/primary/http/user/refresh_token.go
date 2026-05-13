package handler

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/mapper"
	"codebase-app/internal/entity/restentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
)

func (h *userHandler) refreshToken(c fiber.Ctx) error {
	tracedCtx, span := tracing.StartSpan(c.Context(), "internal:framework:primary:http:user:refresh_token:refreshToken")
	defer span.End()
	c.SetContext(tracedCtx)

	var (
		ctx = c.Context()
		req = new(restentity.RefreshTokenReq)
		v   = adapter.Adapters.Validator
	)

	if err := c.Bind().Body(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Any(common.LogKeyPayload, req).Msg("Invalid request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if handled, err := h.limitAuthenticationRequest(c, "user-refresh-token"); err != nil || handled {
		return err
	}

	if err := v.Validate(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Any(common.LogKeyPayload, req).Msg("Invalid request validation")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	user := mapper.RefreshTokenReqToCore(ctx, *req)
	tokens, err := h.core.RefreshToken(ctx, user)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("RefreshToken service error")
		code, errs := errmsg.Errors[error](ctx, err)
		return c.Status(code).JSON(response.Error(errs))
	}

	resp := mapper.AuthTokensToLoginResp(*tokens)
	return c.Status(fiber.StatusOK).JSON(response.Success(resp, ""))
}
