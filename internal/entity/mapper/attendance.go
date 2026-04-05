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
		Type:           string(item.Type),
		Source:         string(item.Source),
		Status:         string(item.Status),
		LoggedAt:       item.LoggedAt,
		Latitude:       item.Latitude,
		Longitude:      item.Longitude,
		Address:        item.Address,
		DeviceID:       item.DeviceID,
		DeviceName:     item.DeviceName,
		Notes:          item.Notes,
		SelfieFileID:   item.SelfieFileID,
		SelfieURL:      item.SelfieURL,
		CreatedAt:      item.CreatedAt,
		UpdatedAt:      item.UpdatedAt,
	}
}

func AttendanceActionFromRest(ctx context.Context, req restentity.CheckAttendanceReq, actionType string) coreentity.AttendanceLogAction {
	uc := common.GetUserContext(ctx)

	return coreentity.AttendanceLogAction{
		UserCtx:      uc,
		TenantID:     uc.TenantID,
		Type:         common.AttendanceType(actionType),
		LoggedAt:     req.LoggedAt,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		Address:      req.Address,
		DeviceID:     req.DeviceID,
		DeviceName:   req.DeviceName,
		Notes:        req.Notes,
		SelfieFileID: req.SelfieFileID,
	}
}

func AttendanceSummaryFromCoreToRest(item coreentity.AttendanceSummary) restentity.GetAttendanceSummaryTodayResp {
	var lastLogType *string
	if item.LastLogType != nil {
		value := string(*item.LastLogType)
		lastLogType = &value
	}

	return restentity.GetAttendanceSummaryTodayResp{
		AttendanceDate: item.AttendanceDate,
		CheckedIn:      item.CheckedIn,
		CheckedOut:     item.CheckedOut,
		CheckInLogID:   item.CheckInLogID,
		CheckOutLogID:  item.CheckOutLogID,
		LastLogType:    lastLogType,
		LastLoggedAt:   item.LastLoggedAt,
		CanCheckIn:     item.CanCheckIn,
		CanCheckOut:    item.CanCheckOut,
	}
}

func AttendanceTypePtrFromString(value string) *common.AttendanceType {
	if value == "" {
		return nil
	}

	result := common.AttendanceType(value)
	return &result
}

func AttendanceSourcePtrFromString(value string) *common.AttendanceSource {
	if value == "" {
		return nil
	}

	result := common.AttendanceSource(value)
	return &result
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
