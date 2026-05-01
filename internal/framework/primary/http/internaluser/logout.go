package handler

import (
	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/mapper"
	"codebase-app/internal/entity/restentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (h *internalUserHandler) logout(c *fiber.Ctx) error {
	tracedCtx, span := tracing.StartSpan(c.UserContext(), "internal:framework:primary:http:internaluser:logout:logout")
	defer span.End()
	c.SetUserContext(tracedCtx)

	var (
		ctx = c.UserContext()
		req = new(restentity.InternalUserLogoutReq)
	)

	if err := c.BodyParser(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Any(common.LogKeyPayload, req).Msg("Invalid internal logout request body")
		code, errs := errmsg.Errors(ctx, err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	if err := h.core.Logout(ctx, mapper.InternalUserLogoutReqToCore(ctx, *req)); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Internal logout service error")
		code, errs := errmsg.Errors[error](ctx, err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(nil, ""))
}
