package mapper

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/entity/restentity"
)

// PermissionFromCoreToRest maps coreentity.Permission to restentity.Permission
func PermissionFromCoreToRest(item coreentity.Permission) restentity.Permission {
	return restentity.Permission{
		ID:          item.ID,
		Name:        item.Name,
		Description: item.Description,
	}
}

// RoleFromRestCreateToCore maps create request to coreentity.Role
func RoleFromRestCreateToCore(ctx context.Context, req restentity.CreateRoleReq) coreentity.Role {
	uc := common.GetUserContext(ctx)
	tenantID := tenantIDPtr(uc)
	return coreentity.Role{
		UserCtx:  uc,
		TenantID: tenantID,
		Name:     req.Name,
	}
}

// RoleFromRestUpdateToCore maps update request to coreentity.Role
func RoleFromRestUpdateToCore(ctx context.Context, req restentity.UpdateRoleReq) coreentity.Role {
	uc := common.GetUserContext(ctx)
	tenantID := tenantIDPtr(uc)
	return coreentity.Role{
		UserCtx:  uc,
		TenantID: tenantID,
		ID:       req.ID,
		Name:     req.Name,
	}
}

// tenantIDPtr returns a pointer to the TenantID string, or nil if empty
func tenantIDPtr(uc common.UserContext) *string {
	if uc.TenantID != "" {
		return &uc.TenantID
	}
	return nil
}
