package handler

import (
	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/mapper"
	"codebase-app/internal/entity/restentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
)

func (h *userHandler) logout(c fiber.Ctx) error {
	tracedCtx, span := tracing.StartSpan(c.Context(), "internal:framework:primary:http:user:logout:logout")
	defer span.End()
	c.SetContext(tracedCtx)

	var (
		ctx = c.Context()
		req = new(restentity.LogoutReq)
	)

	if err := c.Bind().Body(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Any(common.LogKeyPayload, req).Msg("Invalid logout request body")
		code, errs := errmsg.Errors(ctx, err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	if err := h.core.Logout(ctx, mapper.LogoutReqToCore(ctx, *req)); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Logout service error")
		code, errs := errmsg.Errors[error](ctx, err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(nil, ""))
}
