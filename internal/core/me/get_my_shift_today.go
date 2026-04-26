package core

import (
	"context"
	"time"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (c *meCore) GetMyShiftToday(
	ctx context.Context,
	filter coreentity.SelfFilter,
) (*coreentity.WorkShiftToday, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:me:get_my_shift_today:GetMyShiftToday")
	defer span.End()

	shift, err := c.repo.GetWorkShiftByUserID(ctx, filter.TenantID, filter.UserID)
	if err != nil {
		return nil, err
	}
	if shift == nil {
		return nil, nil
	}

	location := time.Now().Location()
	if shift.Timezone != "" {
		loadedLocation, loadErr := time.LoadLocation(shift.Timezone)
		if loadErr == nil {
			location = loadedLocation
		}
	}

	return &coreentity.WorkShiftToday{
		UserCtx:        filter.UserCtx,
		ShiftID:        shift.ID,
		ShiftName:      shift.Name,
		StartTime:      shift.StartTime,
		EndTime:        shift.EndTime,
		Timezone:       shift.Timezone,
		AttendanceDate: time.Now().In(location).Format("2006-01-02"),
	}, nil
}
