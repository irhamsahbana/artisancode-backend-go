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

func (h *userHandler) googleRegister(c fiber.Ctx) error {
	tracedCtx, span := tracing.StartSpan(
		c.Context(),
		"internal:framework:primary:http:user:google_register:googleRegister",
	)
	defer span.End()
	c.SetContext(tracedCtx)

	var (
		ctx = c.Context()
		req = new(restentity.GoogleRegisterReq)
		v   = adapter.Adapters.Validator
	)

	if err := c.Bind().Body(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Any(common.LogKeyPayload, req.Log()).Msg("Invalid request body")
		code, errs := errmsg.Errors(ctx, err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	if handled, err := h.limitRegistrationRequest(
		c,
		"user-google-register",
		normalizedTenantCode(req.TenantCode),
	); err != nil || handled {
		return err
	}

	if err := v.Validate(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Any(common.LogKeyPayload, req.Log()).Msg("Invalid request validation")
		code, errs := errmsg.Errors(ctx, err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	if req.IDToken == "" && req.RegistrationToken == "" {
		errResp := errmsg.NewCustomErrors(400).
			SetMessage(errmsg.MessageGoogleRegistrationSessionIsInvalidOrExpired).
			SetErrorCode("google_registration_session_invalid")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(errResp))
	}

	result, err := h.core.GoogleRegister(ctx, mapper.GoogleRegisterReqToCore(ctx, *req))
	if err != nil {
		log.Ctx(ctx).Warn().Err(err).Any(common.LogKeyPayload, req.Log()).Msg("Google register service error")
		code, errs := errmsg.Errors[error](ctx, err)
		return c.Status(code).JSON(response.Error(errs))
	}

	resp := mapper.GoogleRegisterResultToResp(*result)
	return c.Status(fiber.StatusCreated).JSON(response.Success(resp, errmsg.MessageGoogleRegistrationSuccessful))
}
