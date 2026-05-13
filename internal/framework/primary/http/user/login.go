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

func (h *userHandler) login(c fiber.Ctx) error {
	tracedCtx, span := tracing.StartSpan(c.Context(), "internal:framework:primary:http:user:login:login")
	defer span.End()
	c.SetContext(tracedCtx)

	var (
		ctx = c.Context()
		req = new(restentity.LoginReq)
		v   = adapter.Adapters.Validator
	)

	if err := c.Bind().Body(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Any(common.LogKeyPayload, req).Msg("Invalid request body")
		code, errs := errmsg.Errors(ctx, err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	if handled, err := h.limitAuthenticationRequest(
		c,
		"user-login",
		normalizedTenantCode(req.TenantCode),
		normalizedEmail(req.Email),
	); err != nil || handled {
		return err
	}

	if err := v.Validate(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Any(common.LogKeyPayload, req).Msg("Invalid request validation")
		code, errs := errmsg.Errors(ctx, err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	user := mapper.LoginReqToCore(ctx, *req)
	tokens, err := h.core.Login(ctx, user)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, map[string]string{"email": req.Email}).Msg("Login service error")
		code, errs := errmsg.Errors[error](ctx, err)
		return c.Status(code).JSON(response.Error(errs))
	}

	resp := mapper.AuthTokensToLoginResp(*tokens)
	return c.Status(fiber.StatusOK).JSON(response.Success(resp, ""))
}
