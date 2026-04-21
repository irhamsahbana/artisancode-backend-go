package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
)

func (c *userInvitationCore) GetInvitationByToken(ctx context.Context, token string) (*coreentity.UserInvitation, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:userinvitation:get_invitation_by_token:GetInvitationByToken")
	defer span.End()

	item, err := c.repo.GetInvitationByTokenHash(ctx, hashInvitationToken(token))
	if err != nil {
		return nil, err
	}

	return validateAcceptableInvitation(item)
}

func validateAcceptableInvitation(item *coreentity.UserInvitation) (*coreentity.UserInvitation, error) {
	if item.Status == coreentity.UserInvitationStatusAccepted {
		return nil, errmsg.NewCustomErrors(400).SetMessage("Invitation has already been accepted")
	}
	if item.Status == coreentity.UserInvitationStatusRevoked {
		return nil, errmsg.NewCustomErrors(400).SetMessage("Invitation has been revoked")
	}
	if isExpiredInvitation(item) {
		return nil, errmsg.NewCustomErrors(410).SetMessage("Invitation has expired")
	}

	return item, nil
}
