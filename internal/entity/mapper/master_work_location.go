package mapper

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/entity/restentity"
)

func WorkLocationFromCoreToRest(item coreentity.WorkLocation) restentity.WorkLocation {
	return restentity.WorkLocation{
		ID:           item.ID,
		OrgUnitID:    item.OrgUnitID,
		OrgUnitName:  item.OrgUnitName,
		Name:         item.Name,
		Address:      item.Address,
		Latitude:     item.Latitude,
		Longitude:    item.Longitude,
		RadiusMeters: item.RadiusMeters,
	}
}

func WorkLocationFromRestCreateToCore(ctx context.Context, req restentity.CreateWorkLocationReq) coreentity.WorkLocation {
	uc := common.GetUserContext(ctx)
	return coreentity.WorkLocation{
		UserCtx:      uc,
		TenantID:     uc.TenantID,
		Name:         req.Name,
		OrgUnitID:    req.OrgUnitID,
		Address:      req.Address,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		RadiusMeters: req.RadiusMeters,
	}
}

func WorkLocationFromRestUpdateToCore(ctx context.Context, req restentity.UpdateWorkLocationReq) coreentity.WorkLocation {
	uc := common.GetUserContext(ctx)
	return coreentity.WorkLocation{
		UserCtx:      uc,
		ID:           req.ID,
		TenantID:     uc.TenantID,
		Name:         req.Name,
		OrgUnitID:    req.OrgUnitID,
		Address:      req.Address,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		RadiusMeters: req.RadiusMeters,
	}
}
