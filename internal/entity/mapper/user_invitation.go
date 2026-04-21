package mapper

import (
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/entity/restentity"
)

func UserInvitationFromCoreToRest(item coreentity.UserInvitation) restentity.UserInvitationResource {
	return restentity.UserInvitationResource{
		ID:           item.ID,
		TenantCode:   item.TenantCode,
		TenantName:   item.TenantName,
		EmployeeID:   item.EmployeeID,
		EmployeeNo:   item.EmployeeNo,
		EmployeeName: item.EmployeeName,
		Email:        item.Email,
		RoleCode:     item.RoleCode,
		Status:       item.Status,
		ExpiresAt:    item.ExpiresAt,
		AcceptedAt:   item.AcceptedAt,
		RevokedAt:    item.RevokedAt,
		LastSentAt:   item.LastSentAt,
		CreatedAt:    item.CreatedAt,
	}
}
