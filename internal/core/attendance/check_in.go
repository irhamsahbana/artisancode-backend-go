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
	ctx, span := tracing.StartSpan(ctx, "internal:core:attendance:check_in:CheckIn")
	defer span.End()

	return c.createAttendanceLog(ctx, data, common.AttendanceTypeCheckIn)
}

func (c *attendanceCore) createAttendanceLog(ctx context.Context, data coreentity.AttendanceLogAction, attendanceType common.AttendanceType) (*coreentity.AttendanceLog, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:attendance:check_in:createAttendanceLog")
	defer span.End()

	eventTime, err := resolveAttendanceTime(data.LoggedAt)
	if err != nil {
		log.Ctx(ctx).Warn().Any(common.LogKeyPayload, data).Msg("Invalid attendance timestamp")
		return nil, errmsg.NewCustomErrors(400).SetMessage("Invalid logged_at format")
	}

	attendanceDate := eventTime.Format("2006-01-02")
	var item *coreentity.AttendanceLog
	err = c.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		employee, err := c.repo.GetEmployeeByUserID(txCtx, data.TenantID, data.UserCtx.UserID)
		if err != nil {
			return err
		}
		if employee.ShiftID == nil || *employee.ShiftID == "" {
			log.Ctx(txCtx).Warn().Any(common.LogKeyPayload, map[string]string{
				"employee_id": employee.ID,
				"tenant_id":   data.TenantID,
			}).Msg("Employee does not have a work shift")
			return errmsg.NewCustomErrors(400).SetMessage("Work shift is required")
		}

		shift, err := c.repo.GetWorkShift(txCtx, coreentity.WorkShift{
			TenantID: data.TenantID,
			ID:       *employee.ShiftID,
		})
		if err != nil {
			return err
		}

		exists, err := c.repo.ExistsAttendanceByTypeOnDate(txCtx, data.TenantID, employee.ID, attendanceDate, string(attendanceType))
		if err != nil {
			return err
		}
		if exists {
			message := "Check in already recorded for today"
			if attendanceType == common.AttendanceTypeCheckOut {
				message = "Check out already recorded for today"
			}
			log.Ctx(txCtx).Warn().Any(common.LogKeyPayload, data).Msg(message)
			return errmsg.NewCustomErrors(400).SetMessage(message)
		}

		if attendanceType == common.AttendanceTypeCheckOut {
			hasCheckIn, err := c.repo.ExistsAttendanceByTypeOnDate(txCtx, data.TenantID, employee.ID, attendanceDate, string(common.AttendanceTypeCheckIn))
			if err != nil {
				return err
			}
			if !hasCheckIn {
				log.Ctx(txCtx).Warn().Any(common.LogKeyPayload, data).Msg("Check out attempted without a corresponding check in")
				return errmsg.NewCustomErrors(400).SetMessage("Check in must be recorded before check out")
			}
		}

		file, err := c.storageRepo.GetFile(txCtx, coreentity.FileFilter{
			TenantID: data.TenantID,
			ID:       data.SelfieFileID,
		})
		if err != nil {
			log.Ctx(txCtx).Warn().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to get selfie file for attendance log")
			return errmsg.NewCustomErrors(400).SetMessage("Selfie file not found")
		}
		if file.Folder != common.S3FolderAttendanceFace {
			log.Ctx(txCtx).Warn().Any(common.LogKeyPayload, data).Msg("Invalid selfie file folder")
			return errmsg.NewCustomErrors(400).SetMessage("Invalid selfie file")
		}
		if file.Status != common.FileStatusPending {
			log.Ctx(txCtx).Warn().Any(common.LogKeyPayload, data).Msg("Selfie file is no longer pending")
			return errmsg.NewCustomErrors(400).SetMessage("Selfie file is no longer available")
		}

		item, err = c.repo.CreateAttendanceLog(txCtx, coreentity.AttendanceLog{
			UserCtx:        data.UserCtx,
			TenantID:       data.TenantID,
			EmployeeID:     employee.ID,
			EmployeeNo:     employee.EmployeeNo,
			EmployeeName:   employee.FullName,
			AttendanceDate: attendanceDate,
			ShiftID:        &shift.ID,
			ShiftName:      &shift.Name,
			ShiftTimezone:  &shift.Timezone,
			ShiftStartTime: &shift.StartTime,
			ShiftEndTime:   &shift.EndTime,
			ShiftGraceMins: &shift.GracePeriodMinutes,
			Type:           attendanceType,
			Source:         common.AttendanceSourceMobile,
			Status:         common.AttendanceStatusRecorded,
			LoggedAt:       eventTime.Format(time.RFC3339),
			Latitude:       data.Latitude,
			Longitude:      data.Longitude,
			Address:        data.Address,
			DeviceID:       data.DeviceID,
			DeviceName:     data.DeviceName,
			Notes:          data.Notes,
		})
		if err != nil {
			return err
		}

		_, err = c.storageRepo.CreateFileLink(txCtx, coreentity.CreateStorageFileLinkReq{
			TenantID:      data.TenantID,
			StorageFileID: data.SelfieFileID,
			ResourceType:  "attendance_log",
			ResourceID:    item.ID,
			FieldName:     "selfie",
			SortOrder:     1,
		})
		if err != nil {
			log.Ctx(txCtx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to create storage file link for attendance log")
			return err
		}

		err = c.storageRepo.MarkFileAttached(txCtx, data.TenantID, data.SelfieFileID)
		if err != nil {
			log.Ctx(txCtx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to mark storage file attached for attendance log")
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return item, nil
}

func resolveAttendanceTime(value string) (time.Time, error) {
	if value == "" {
		return time.Now().UTC(), nil
	}

	return time.Parse(time.RFC3339, value)
}
