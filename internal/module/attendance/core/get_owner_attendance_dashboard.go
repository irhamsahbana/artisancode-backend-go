package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/gofiber/fiber/v2"
)

func (c *attendanceCore) GetOwnerAttendanceDashboard(ctx context.Context, filter coreentity.OwnerAttendanceDashboardFilter) (*coreentity.OwnerAttendanceDashboard, error) {
	ctx, span := tracing.StartSpan(ctx, "core.GetOwnerAttendanceDashboard")
	defer span.End()

	if !filter.UserCtx.HasRole("owner") && !filter.UserCtx.HasRole("admin") {
		return nil, errmsg.NewCustomErrors(fiber.StatusForbidden, errmsg.WithMessage("You are not allowed to access this resource"))
	}

	summary, err := c.repo.GetOwnerAttendanceDashboardSummary(ctx, filter)
	if err != nil {
		return nil, err
	}

	exceptions, err := c.repo.GetOwnerAttendanceDashboardExceptions(ctx, filter)
	if err != nil {
		return nil, err
	}

	trend, err := c.repo.GetOwnerAttendanceDashboardTrend(ctx, filter)
	if err != nil {
		return nil, err
	}

	return &coreentity.OwnerAttendanceDashboard{
		UserCtx:         filter.UserCtx,
		AttendanceDate:  filter.Date,
		Summary:         *summary,
		TodayExceptions: exceptions,
		DailyTrend:      trend,
	}, nil
}
