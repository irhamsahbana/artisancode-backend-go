package handler

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/entity/mapper"
	"codebase-app/internal/entity/restentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
)

func (h *jobPositionHandler) updateJobPosition(c fiber.Ctx) error {
	tracedCtx, span := tracing.StartSpan(c.Context(), "internal:framework:primary:http:jobposition:handler:updateJobPosition")
	defer span.End()
	c.SetContext(tracedCtx)

	var (
		ctx = c.Context()
		req = new(restentity.UpdateJobPositionReq)
		v   = adapter.Adapters.Validator
	)

	if err := c.Bind().URI(req); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to parse params")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := c.Bind().Body(req); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := v.Validate(req); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to validate request body")
		code, errors := errmsg.Errors(ctx, err, req)
		return c.Status(code).JSON(response.Error(errors))
	}

	data := mapper.JobPositionFromRestUpdateToCore(ctx, *req)
	if err := h.core.UpdateJobPosition(ctx, data); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to update job position")
		code, errors := errmsg.Errors[error](ctx, err)
		return c.Status(code).JSON(response.Error(errors))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(nil, ""))
}
