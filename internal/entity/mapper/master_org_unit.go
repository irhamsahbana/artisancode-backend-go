package mapper

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/entity/repoentity"
	"codebase-app/internal/entity/restentity"
)

func OrgUnitFromRepoToCore(ctx context.Context, item repoentity.OrgUnit) coreentity.OrgUnit {
	uc := common.GetUserContext(ctx)
	return coreentity.OrgUnit{
		UserCtx:  uc,
		ID:       item.ID,
		TenantID: item.TenantID,
		Name:     item.Name,
		ParentID: item.ParentID,
		Category: item.Category,
	}
}

func OrgUnitFromCoreToRest(item coreentity.OrgUnit) restentity.OrgUnit {
	return restentity.OrgUnit{
		ID:       item.ID,
		Name:     item.Name,
		ParentID: item.ParentID,
		Category: item.Category,
	}
}

func OrgUnitFromRestCreateToCore(ctx context.Context, req restentity.CreateOrgUnitReq) coreentity.OrgUnit {
	uc := common.GetUserContext(ctx)
	return coreentity.OrgUnit{
		UserCtx:  uc,
		TenantID: uc.TenantID,
		Name:     req.Name,
		ParentID: req.ParentID,
		Category: req.Category,
	}
}

func OrgUnitFromRestUpdateToCore(ctx context.Context, req restentity.UpdateOrgUnitReq) coreentity.OrgUnit {
	uc := common.GetUserContext(ctx)
	return coreentity.OrgUnit{
		UserCtx:  uc,
		ID:       req.ID,
		TenantID: uc.TenantID,
		Name:     req.Name,
		ParentID: req.ParentID,
		Category: req.Category,
	}
}
