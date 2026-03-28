package mapper

import (
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/entity/restentity"
)

func MeFromCoreToRest(item coreentity.Me) restentity.MeResp {
	return restentity.MeResp{
		UserID:      item.UserID,
		UserName:    item.UserName,
		TenantID:    item.TenantID,
		TenantName:  item.TenantName,
		Roles:       item.Roles,
		CompanyID:   item.CompanyID,
		CompanyName: item.CompanyName,
	}
}

func MyEmployeeFromCoreToRest(item coreentity.Employee) restentity.GetMyEmployeeResp {
	return restentity.GetMyEmployeeResp{
		ID:            item.ID,
		EmployeeNo:    item.EmployeeNo,
		FullName:      item.FullName,
		Email:         item.Email,
		UserID:        item.UserID,
		OrgUnitID:     item.OrgUnitID,
		JobPositionID: item.JobPositionID,
		LocationID:    item.LocationID,
		ShiftID:       item.ShiftID,
		Status:        item.Status,
		JoinDate:      item.JoinDate,
	}
}

func MyShiftTodayFromCoreToRest(item coreentity.WorkShiftToday) restentity.GetMyShiftTodayResp {
	return restentity.GetMyShiftTodayResp{
		ShiftID:        item.ShiftID,
		ShiftName:      item.ShiftName,
		StartTime:      item.StartTime,
		EndTime:        item.EndTime,
		Timezone:       item.Timezone,
		AttendanceDate: item.AttendanceDate,
	}
}
