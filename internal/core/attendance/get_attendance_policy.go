package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
)

func (c *attendanceCore) GetAttendancePolicy(
	ctx context.Context,
	filter coreentity.SelfFilter,
) (*coreentity.AttendancePolicy, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:attendance:get_attendance_policy:GetAttendancePolicy")
	defer span.End()

	employee, err := c.repo.GetEmployeeByUserID(ctx, filter.TenantID, filter.UserID)
	if err != nil {
		return nil, err
	}
	if employee.ShiftID == nil || *employee.ShiftID == "" {
		return nil, errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageWorkShiftIsRequired)
	}

	shift, err := c.repo.GetWorkShift(ctx, coreentity.WorkShift{
		TenantID: filter.TenantID,
		ID:       *employee.ShiftID,
	})
	if err != nil {
		return nil, err
	}

	radiusMeters := 0
	if employee.LocationID != nil && *employee.LocationID != "" {
		location, err := c.repo.GetWorkLocation(ctx, coreentity.WorkLocation{
			TenantID: filter.TenantID,
			ID:       *employee.LocationID,
		})
		if err != nil {
			return nil, err
		}
		if location.RadiusMeters != nil {
			radiusMeters = *location.RadiusMeters
		}
	}

	return &coreentity.AttendancePolicy{
		UserCtx:                 filter.UserCtx,
		Timezone:                shift.Timezone,
		AttendanceRadiusMeters:  radiusMeters,
		AttendanceCheckInStart:  shift.StartTime,
		AttendanceCheckInEnd:    shift.StartTime,
		AttendanceCheckOutStart: shift.EndTime,
		AttendanceCheckOutEnd:   shift.EndTime,
	}, nil
}
