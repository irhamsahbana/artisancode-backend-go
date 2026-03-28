package mapper

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/entity/restentity"
)

func AttendanceLogFromCoreToRest(item coreentity.AttendanceLog) restentity.AttendanceLog {
	return restentity.AttendanceLog{
		ID:             item.ID,
		EmployeeID:     item.EmployeeID,
		EmployeeNo:     item.EmployeeNo,
		EmployeeName:   item.EmployeeName,
		AttendanceDate: item.AttendanceDate,
		Type:           item.Type,
		Source:         item.Source,
		Status:         item.Status,
		LoggedAt:       item.LoggedAt,
		Latitude:       item.Latitude,
		Longitude:      item.Longitude,
		Address:        item.Address,
		DeviceID:       item.DeviceID,
		DeviceName:     item.DeviceName,
		Notes:          item.Notes,
		CreatedAt:      item.CreatedAt,
		UpdatedAt:      item.UpdatedAt,
	}
}

func AttendanceActionFromRest(ctx context.Context, req restentity.CheckAttendanceReq, actionType string) coreentity.AttendanceLogAction {
	uc := common.GetUserContext(ctx)

	return coreentity.AttendanceLogAction{
		UserCtx:    uc,
		TenantID:   uc.TenantID,
		Type:       actionType,
		LoggedAt:   req.LoggedAt,
		Latitude:   req.Latitude,
		Longitude:  req.Longitude,
		Address:    req.Address,
		DeviceID:   req.DeviceID,
		DeviceName: req.DeviceName,
		Notes:      req.Notes,
	}
}

func AttendanceSummaryFromCoreToRest(item coreentity.AttendanceSummary) restentity.GetAttendanceSummaryTodayResp {
	return restentity.GetAttendanceSummaryTodayResp{
		AttendanceDate: item.AttendanceDate,
		CheckedIn:      item.CheckedIn,
		CheckedOut:     item.CheckedOut,
		CheckInLogID:   item.CheckInLogID,
		CheckOutLogID:  item.CheckOutLogID,
		LastLogType:    item.LastLogType,
		LastLoggedAt:   item.LastLoggedAt,
		CanCheckIn:     item.CanCheckIn,
		CanCheckOut:    item.CanCheckOut,
	}
}

func AttendancePolicyFromCoreToRest(item coreentity.AttendancePolicy) restentity.GetAttendancePolicyResp {
	return restentity.GetAttendancePolicyResp{
		Timezone:                item.Timezone,
		AttendanceRadiusMeters:  item.AttendanceRadiusMeters,
		AttendanceCheckInStart:  item.AttendanceCheckInStart,
		AttendanceCheckInEnd:    item.AttendanceCheckInEnd,
		AttendanceCheckOutStart: item.AttendanceCheckOutStart,
		AttendanceCheckOutEnd:   item.AttendanceCheckOutEnd,
	}
}
