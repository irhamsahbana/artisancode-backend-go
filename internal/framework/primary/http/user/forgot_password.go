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

func (h *userHandler) forgotPassword(c fiber.Ctx) error {
	tracedCtx, span := tracing.StartSpan(c.Context(), "internal:framework:primary:http:user:forgot_password:forgotPassword")
	defer span.End()
	c.SetContext(tracedCtx)

	var (
		ctx = c.Context()
		req = new(restentity.ForgotPasswordReq)
		v   = adapter.Adapters.Validator
	)

	if err := c.Bind().Body(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Any(common.LogKeyPayload, req).Msg("Invalid request body")
		code, errs := errmsg.Errors(ctx, err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	if err := v.Validate(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Any(common.LogKeyPayload, req).Msg("Invalid request validation")
		code, errs := errmsg.Errors(ctx, err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	if err := h.limitForgotPassword(c, req.Email, req.TenantCode); err != nil {
		return err
	}

	user := mapper.ForgotPasswordReqToCore(ctx, *req)
	err := h.core.ForgotPassword(ctx, user)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, map[string]string{
			"email":       req.Email,
			"tenant_code": req.TenantCode,
		}).Msg("Forgot password service error")
		code, errs := errmsg.Errors[error](ctx, err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(map[string]any{}, errmsg.MessageIfTheEmailIsRegisteredAPasswordResetLinkHasBeenSent))
}
