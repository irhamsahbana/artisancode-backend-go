package handler

import (
	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/entity/mapper"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
)

func (h *attendanceHandler) getAttendanceSummaryToday(c fiber.Ctx) error {
	tracedCtx, span := tracing.StartSpan(c.Context(), "internal:framework:primary:http:attendance:get_attendance_summary_today:getAttendanceSummaryToday")
	defer span.End()
	c.SetContext(tracedCtx)

	ctx := c.Context()
	uc := common.GetUserContext(ctx)

	item, err := h.core.GetAttendanceSummaryToday(ctx, coreentity.SelfFilter{
		UserCtx:  uc,
		TenantID: uc.TenantID,
		UserID:   uc.UserID,
	})
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get attendance summary today")
		code, errors := errmsg.Errors[error](ctx, err)
		return c.Status(code).JSON(response.Error(errors))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(mapper.AttendanceSummaryFromCoreToRest(*item), ""))
}
