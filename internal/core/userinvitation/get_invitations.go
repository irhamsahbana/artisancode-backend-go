package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
)

func (c *userInvitationCore) GetInvitations(
	ctx context.Context,
	filter coreentity.UserInvitationListFilter,
) ([]coreentity.UserInvitation, int, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:userinvitation:get_invitations:GetInvitations")
	defer span.End()

	if !filter.UserCtx.HasRole("owner") && !filter.UserCtx.HasRole("admin") {
		return nil, 0, errmsg.NewCustomErrors(403).SetMessage("You are not authorized to view invitations")
	}

	return c.repo.GetInvitations(ctx, filter)
}
