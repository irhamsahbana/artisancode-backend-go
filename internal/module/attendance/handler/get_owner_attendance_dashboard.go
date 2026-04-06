package handler

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/entity/mapper"
	"codebase-app/internal/entity/restentity"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (h *attendanceHandler) getOwnerAttendanceDashboard(c *fiber.Ctx) error {
	var (
		ctx = c.UserContext()
		req = new(restentity.GetOwnerAttendanceDashboardReq)
		v   = adapter.Adapters.Validator
	)

	if err := c.QueryParser(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Failed to parse owner attendance dashboard query params")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	req.SetDefault()

	if err := v.Validate(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Failed to validate owner attendance dashboard query params")
		code, errors := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errors))
	}

	filter := coreentity.OwnerAttendanceDashboardFilter{
		UserCtx:   common.GetUserContext(ctx),
		TenantID:  common.GetUserContext(ctx).TenantID,
		Date:      *req.Date,
		Timezone:  req.Timezone,
		TrendDays: req.TrendDays,
	}

	data, err := h.core.GetOwnerAttendanceDashboard(ctx, filter)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get owner attendance dashboard")
		code, errors := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errors))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(mapper.OwnerAttendanceDashboardFromCoreToRest(*data), ""))
}
