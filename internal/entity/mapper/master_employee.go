package mapper

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/entity/restentity"
)

func EmployeeFromCoreToRest(item coreentity.Employee) restentity.Employee {
	return restentity.Employee{
		ID:               item.ID,
		EmployeeNo:       item.EmployeeNo,
		FullName:         item.FullName,
		Email:            item.Email,
		OrgUnitID:        item.OrgUnitID,
		JobPositionID:    item.JobPositionID,
		LocationID:       item.LocationID,
		ShiftID:          item.ShiftID,
		Status:           item.Status,
		AccessStatus:     item.AccessStatus,
		JoinDate:         item.JoinDate,
		JoinDateTimezone: item.JoinDateTimezone,
	}
}

func EmployeeFromRestCreateToCore(ctx context.Context, req restentity.CreateEmployeeReq) coreentity.Employee {
	uc := common.GetUserContext(ctx)
	return coreentity.Employee{
		UserCtx:          uc,
		TenantID:         uc.TenantID,
		EmployeeNo:       req.EmployeeNo,
		FullName:         req.FullName,
		Email:            req.Email,
		Password:         req.Password,
		OrgUnitID:        req.OrgUnitID,
		JobPositionID:    req.JobPositionID,
		LocationID:       req.LocationID,
		ShiftID:          req.ShiftID,
		Status:           req.Status,
		JoinDate:         req.JoinDate,
		JoinDateTimezone: &req.JoinDateTimezone,
	}
}

func EmployeeFromRestUpdateToCore(ctx context.Context, req restentity.UpdateEmployeeReq) coreentity.Employee {
	uc := common.GetUserContext(ctx)
	return coreentity.Employee{
		UserCtx:          uc,
		ID:               req.ID,
		TenantID:         uc.TenantID,
		EmployeeNo:       req.EmployeeNo,
		FullName:         req.FullName,
		Email:            req.Email,
		Password:         req.Password,
		OrgUnitID:        req.OrgUnitID,
		JobPositionID:    req.JobPositionID,
		LocationID:       req.LocationID,
		ShiftID:          req.ShiftID,
		Status:           req.Status,
		JoinDate:         req.JoinDate,
		JoinDateTimezone: &req.JoinDateTimezone,
	}
}
