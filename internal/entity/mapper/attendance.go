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
