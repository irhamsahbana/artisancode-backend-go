package core

import (
	"context"
	"time"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/rs/zerolog/log"
)

func (c *attendanceCore) CheckIn(ctx context.Context, data coreentity.AttendanceLogAction) (*coreentity.AttendanceLog, error) {
	ctx, span := tracing.StartSpan(ctx, "core.CheckIn")
	defer span.End()

	return c.createAttendanceLog(ctx, data, "check_in")
}

func (c *attendanceCore) createAttendanceLog(ctx context.Context, data coreentity.AttendanceLogAction, attendanceType string) (*coreentity.AttendanceLog, error) {
	eventTime, err := resolveAttendanceTime(data.LoggedAt)
	if err != nil {
		log.Ctx(ctx).Warn().Any(common.LogKeyPayload, data).Msg("Invalid attendance timestamp")
		return nil, errmsg.NewCustomErrors(400).SetMessage("Invalid logged_at format")
	}

	employee, err := c.repo.GetEmployeeByUserID(ctx, data.TenantID, data.UserCtx.UserID)
	if err != nil {
		return nil, err
	}

	attendanceDate := eventTime.Format("2006-01-02")
	exists, err := c.repo.ExistsAttendanceByTypeOnDate(ctx, data.TenantID, employee.ID, attendanceDate, attendanceType)
	if err != nil {
		return nil, err
	}
	if exists {
		message := "Check in already recorded for today"
		if attendanceType == "check_out" {
			message = "Check out already recorded for today"
		}
		return nil, errmsg.NewCustomErrors(400).SetMessage(message)
	}

	if attendanceType == "check_out" {
		hasCheckIn, err := c.repo.ExistsAttendanceByTypeOnDate(ctx, data.TenantID, employee.ID, attendanceDate, "check_in")
		if err != nil {
			return nil, err
		}
		if !hasCheckIn {
			return nil, errmsg.NewCustomErrors(400).SetMessage("Check in must be recorded before check out")
		}
	}

	return c.repo.CreateAttendanceLog(ctx, coreentity.AttendanceLog{
		UserCtx:        data.UserCtx,
		TenantID:       data.TenantID,
		EmployeeID:     employee.ID,
		EmployeeNo:     employee.EmployeeNo,
		EmployeeName:   employee.FullName,
		AttendanceDate: attendanceDate,
		Type:           attendanceType,
		Source:         "mobile",
		Status:         "recorded",
		LoggedAt:       eventTime.Format(time.RFC3339),
		Latitude:       data.Latitude,
		Longitude:      data.Longitude,
		Address:        data.Address,
		DeviceID:       data.DeviceID,
		DeviceName:     data.DeviceName,
		Notes:          data.Notes,
	})
}

func resolveAttendanceTime(value string) (time.Time, error) {
	if value == "" {
		return time.Now().UTC(), nil
	}

	return time.Parse(time.RFC3339, value)
}
