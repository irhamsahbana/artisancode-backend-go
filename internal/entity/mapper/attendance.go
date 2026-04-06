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

func AttendanceStatusPtrFromString(value string) *common.AttendanceStatus {
	if value == "" {
		return nil
	}

	result := common.AttendanceStatus(value)
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

func OwnerAttendanceDashboardFromCoreToRest(item coreentity.OwnerAttendanceDashboard) restentity.GetOwnerAttendanceDashboardResp {
	exceptions := make([]restentity.OwnerAttendanceDashboardExceptionResp, 0, len(item.TodayExceptions))
	for _, exception := range item.TodayExceptions {
		exceptions = append(exceptions, restentity.OwnerAttendanceDashboardExceptionResp{
			EmployeeID:     exception.EmployeeID,
			EmployeeNo:     exception.EmployeeNo,
			EmployeeName:   exception.EmployeeName,
			ShiftName:      exception.ShiftName,
			FirstCheckInAt: exception.FirstCheckInAt,
			LastCheckOutAt: exception.LastCheckOutAt,
			ExceptionType:  exception.ExceptionType,
		})
	}

	trend := make([]restentity.OwnerAttendanceDashboardDailyTrendResp, 0, len(item.DailyTrend))
	for _, day := range item.DailyTrend {
		trend = append(trend, restentity.OwnerAttendanceDashboardDailyTrendResp{
			AttendanceDate:       day.AttendanceDate,
			CheckedInCount:       day.CheckedInCount,
			CheckedOutCount:      day.CheckedOutCount,
			LateCheckInCount:     day.LateCheckInCount,
			MissingCheckOutCount: day.MissingCheckOutCount,
		})
	}

	return restentity.GetOwnerAttendanceDashboardResp{
		AttendanceDate: item.AttendanceDate,
		Summary: restentity.OwnerAttendanceDashboardSummaryResp{
			ActiveEmployeeCount:  item.Summary.ActiveEmployeeCount,
			CheckedInCount:       item.Summary.CheckedInCount,
			CheckedOutCount:      item.Summary.CheckedOutCount,
			PendingCheckInCount:  item.Summary.PendingCheckInCount,
			PendingCheckOutCount: item.Summary.PendingCheckOutCount,
			LateCheckInCount:     item.Summary.LateCheckInCount,
		},
		TodayExceptions: exceptions,
		DailyTrend:      trend,
	}
}
