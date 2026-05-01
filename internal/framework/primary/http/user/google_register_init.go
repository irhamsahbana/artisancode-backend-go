package handler

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/mapper"
	"codebase-app/internal/entity/restentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (h *userHandler) googleRegisterInit(c *fiber.Ctx) error {
	tracedCtx, span := tracing.StartSpan(
		c.UserContext(),
		"internal:framework:primary:http:user:google_register_init:googleRegisterInit",
	)
	defer span.End()
	c.SetUserContext(tracedCtx)

	var (
		ctx = c.UserContext()
		req = new(restentity.GoogleRegisterInitReq)
		v   = adapter.Adapters.Validator
	)

	if err := c.BodyParser(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Any(common.LogKeyPayload, req.Log()).Msg("Invalid request body")
		code, errs := errmsg.Errors(ctx, err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	if err := v.Validate(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Any(common.LogKeyPayload, req.Log()).Msg("Invalid request validation")
		code, errs := errmsg.Errors(ctx, err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	result, err := h.core.GoogleRegisterInit(ctx, mapper.GoogleRegisterInitReqToCore(ctx, *req))
	if err != nil {
		log.Ctx(ctx).Warn().Err(err).Any(common.LogKeyPayload, req.Log()).Msg("Google register init service error")
		code, errs := errmsg.Errors[error](ctx, err)
		return c.Status(code).JSON(response.Error(errs))
	}

	resp := mapper.GoogleRegisterInitResultToResp(*result)
	return c.Status(fiber.StatusOK).JSON(response.Success(resp, errmsg.MessageGoogleRegistrationSessionCreated))
}
