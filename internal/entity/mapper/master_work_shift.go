package mapper

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/entity/repoentity"
	"codebase-app/internal/entity/restentity"
)

func WorkShiftFromRepoToCore(ctx context.Context, item repoentity.WorkShift) coreentity.WorkShift {
	uc := common.GetUserContext(ctx)
	return coreentity.WorkShift{
		UserCtx:            uc,
		ID:                 item.ID,
		TenantID:           item.TenantID,
		Name:               item.Name,
		Timezone:           item.Timezone,
		StartTime:          item.StartTime,
		EndTime:            item.EndTime,
		GracePeriodMinutes: item.GracePeriodMinutes,
	}
}

func WorkShiftFromCoreToRest(item coreentity.WorkShift) restentity.WorkShift {
	return restentity.WorkShift{
		ID:                 item.ID,
		Name:               item.Name,
		Timezone:           item.Timezone,
		StartTime:          item.StartTime,
		EndTime:            item.EndTime,
		GracePeriodMinutes: item.GracePeriodMinutes,
	}
}

func WorkShiftFromRestCreateToCore(ctx context.Context, req restentity.CreateWorkShiftReq) coreentity.WorkShift {
	uc := common.GetUserContext(ctx)
	return coreentity.WorkShift{
		UserCtx:            uc,
		TenantID:           uc.TenantID,
		Name:               req.Name,
		Timezone:           req.Timezone,
		StartTime:          req.StartTime,
		EndTime:            req.EndTime,
		GracePeriodMinutes: req.GracePeriodMinutes,
	}
}

func WorkShiftFromRestUpdateToCore(ctx context.Context, req restentity.UpdateWorkShiftReq) coreentity.WorkShift {
	uc := common.GetUserContext(ctx)
	return coreentity.WorkShift{
		UserCtx:            uc,
		ID:                 req.ID,
		TenantID:           uc.TenantID,
		Name:               req.Name,
		Timezone:           req.Timezone,
		StartTime:          req.StartTime,
		EndTime:            req.EndTime,
		GracePeriodMinutes: req.GracePeriodMinutes,
	}
}
