package mapper

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/entity/repoentity"
	"codebase-app/internal/entity/restentity"
)

func JobPositionFromRepoToCore(ctx context.Context, item repoentity.JobPosition) coreentity.JobPosition {
	uc := common.GetUserContext(ctx)
	return coreentity.JobPosition{
		UserCtx:  uc,
		ID:       item.ID,
		TenantID: item.TenantID,
		Name:     item.Name,
		Grade:    item.Grade,
	}
}

func JobPositionFromCoreToRest(item coreentity.JobPosition) restentity.JobPosition {
	return restentity.JobPosition{
		ID:    item.ID,
		Name:  item.Name,
		Grade: item.Grade,
	}
}

func JobPositionFromRestCreateToCore(ctx context.Context, req restentity.CreateJobPositionReq) coreentity.JobPosition {
	uc := common.GetUserContext(ctx)
	return coreentity.JobPosition{
		UserCtx:  uc,
		TenantID: uc.TenantID,
		Name:     req.Name,
		Grade:    req.Grade,
	}
}

func JobPositionFromRestUpdateToCore(ctx context.Context, req restentity.UpdateJobPositionReq) coreentity.JobPosition {
	uc := common.GetUserContext(ctx)
	return coreentity.JobPosition{
		UserCtx:  uc,
		ID:       req.ID,
		TenantID: uc.TenantID,
		Name:     req.Name,
		Grade:    req.Grade,
	}
}
