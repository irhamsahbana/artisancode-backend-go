package mapper

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/entity/repoentity"
	"codebase-app/internal/entity/restentity"
)

func WorkLocationFromRepoToCore(ctx context.Context, item repoentity.WorkLocation) coreentity.WorkLocation {
	uc := common.GetUserContext(ctx)
	return coreentity.WorkLocation{
		UserCtx:      uc,
		ID:           item.ID,
		TenantID:     item.TenantID,
		Name:         item.Name,
		Address:      item.Address,
		Timezone:     item.Timezone,
		Latitude:     item.Latitude,
		Longitude:    item.Longitude,
		RadiusMeters: item.RadiusMeters,
	}
}

func WorkLocationFromCoreToRest(item coreentity.WorkLocation) restentity.WorkLocation {
	return restentity.WorkLocation{
		ID:           item.ID,
		Name:         item.Name,
		Address:      item.Address,
		Timezone:     item.Timezone,
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
		Address:      req.Address,
		Timezone:     req.Timezone,
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
		Address:      req.Address,
		Timezone:     req.Timezone,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		RadiusMeters: req.RadiusMeters,
	}
}