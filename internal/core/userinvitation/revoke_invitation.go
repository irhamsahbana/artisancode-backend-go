package core

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
)

func (c *userInvitationCore) RevokeInvitation(ctx context.Context, data coreentity.UserInvitation) error {
	ctx, span := tracing.StartSpan(ctx, "internal:core:userinvitation:revoke_invitation:RevokeInvitation")
	defer span.End()

	item, err := c.repo.GetInvitationByID(ctx, data.UserCtx.TenantID, data.ID)
	if err != nil {
		return err
	}
	if err := c.authorizeInvitationMutation(item, data.UserCtx); err != nil {
		return err
	}
	if item.Status == coreentity.UserInvitationStatusAccepted {
		return errmsg.NewCustomErrors(400).SetMessage("Accepted invitation cannot be revoked")
	}
	if item.Status == coreentity.UserInvitationStatusRevoked {
		return errmsg.NewCustomErrors(400).SetMessage("Invitation is already revoked")
	}

	return c.repo.RevokeInvitation(ctx, data.UserCtx.TenantID, data.ID)
}

func (c *userInvitationCore) authorizeInvitationMutation(item *coreentity.UserInvitation, userCtx common.UserContext) error {
	if !userCtx.HasRole("owner") && !userCtx.HasRole("admin") {
		return errmsg.NewCustomErrors(403).SetMessage("You are not authorized to manage invitations")
	}
	if item.RoleCode == coreentity.UserInvitationRoleAdmin && !userCtx.HasRole("owner") {
		return errmsg.NewCustomErrors(403).SetMessage("Only owner can manage admin invitations")
	}
	return nil
}
