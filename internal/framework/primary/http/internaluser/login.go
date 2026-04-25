package handler

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/entity/mapper"
	"codebase-app/internal/entity/restentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (h *internalUserHandler) login(c *fiber.Ctx) error {
	tracedCtx, span := tracing.StartSpan(c.UserContext(), "internal:framework:primary:http:internaluser:handler:login")
	defer span.End()
	c.SetUserContext(tracedCtx)

	var (
		ctx = c.UserContext()
		req = new(restentity.InternalUserLoginReq)
		v   = adapter.Adapters.Validator
	)

	if err := c.BodyParser(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Invalid internal user login body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}
	if err := v.Validate(req); err != nil {
		code, errs := errmsg.Errors(ctx, err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	tokens, err := h.core.Login(ctx, mapper.InternalUserLoginReqToCore(ctx, *req))
	if err != nil {
		code, errs := errmsg.Errors[error](ctx, err)
		return c.Status(code).JSON(response.Error(errs))
	}
	return c.Status(fiber.StatusOK).JSON(response.Success(mapper.AuthTokensToInternalUserLoginResp(*tokens), ""))
}
