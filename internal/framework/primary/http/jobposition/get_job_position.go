package handler

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/entity/mapper"
	"codebase-app/internal/entity/restentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
)

func (h *jobPositionHandler) getJobPosition(c fiber.Ctx) error {
	tracedCtx, span := tracing.StartSpan(c.Context(), "internal:framework:primary:http:jobposition:handler:getJobPosition")
	defer span.End()
	c.SetContext(tracedCtx)

	var (
		ctx = c.Context()
		req = new(restentity.GetJobPositionReq)
		v   = adapter.Adapters.Validator
	)

	if err := c.Bind().URI(req); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to parse params")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := v.Validate(req); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to validate params")
		code, errors := errmsg.Errors(ctx, err, req)
		return c.Status(code).JSON(response.Error(errors))
	}

	filter := coreentity.JobPosition{
		TenantID: common.GetUserContext(ctx).TenantID,
		ID:       req.ID,
	}

	item, err := h.core.GetJobPosition(ctx, filter)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get job position")
		code, errors := errmsg.Errors[error](ctx, err)
		return c.Status(code).JSON(response.Error(errors))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(restentity.GetJobPositionResp{
		JobPosition: mapper.JobPositionFromCoreToRest(*item),
	}, ""))
}
