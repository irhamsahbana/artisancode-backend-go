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

func (h *attendanceHandler) getAttendanceLog(c *fiber.Ctx) error {
	var (
		ctx = c.UserContext()
		req = new(restentity.GetAttendanceLogReq)
		v   = adapter.Adapters.Validator
	)

	if err := c.ParamsParser(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Failed to parse params")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := v.Validate(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Failed to validate params")
		code, errors := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errors))
	}

	item, err := h.core.GetAttendanceLog(ctx, coreentity.AttendanceLogDetailFilter{
		UserCtx:  common.GetUserContext(ctx),
		TenantID: common.GetUserContext(ctx).TenantID,
		ID:       req.ID,
	})
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get attendance log")
		code, errors := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errors))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(restentity.GetAttendanceLogResp{
		AttendanceLog: mapper.AttendanceLogFromCoreToRest(*item),
	}, ""))
}
