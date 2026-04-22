package core

import (
	"context"
	"time"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
)

func (c *userInvitationCore) ResendInvitation(ctx context.Context, data coreentity.UserInvitation) (*coreentity.UserInvitation, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:userinvitation:resend_invitation:ResendInvitation")
	defer span.End()

	item, err := c.repo.GetInvitationByID(ctx, data.UserCtx.TenantID, data.ID)
	if err != nil {
		return nil, err
	}
	if err := c.authorizeInvitationMutation(item, data.UserCtx); err != nil {
		return nil, err
	}
	if item.Status == coreentity.UserInvitationStatusAccepted {
		return nil, errmsg.NewCustomErrors(400).SetMessage("Accepted invitation cannot be resent")
	}
	if item.Status == coreentity.UserInvitationStatusRevoked {
		return nil, errmsg.NewCustomErrors(400).SetMessage("Revoked invitation cannot be resent")
	}

	rawToken, tokenHash, err := generateInvitationToken()
	if err != nil {
		return nil, errmsg.NewCustomErrors(500).SetMessage("Failed to resend invitation")
	}

	item.TokenHash = tokenHash
	item.AcceptToken = rawToken
	item.ExpiresAt = time.Now().Add(invitationExpiryDuration)
	item.LastSentAt = time.Now()

	if err := c.repo.ResendInvitation(ctx, *item); err != nil {
		return nil, err
	}

	updated, err := c.repo.GetInvitationByID(ctx, data.UserCtx.TenantID, data.ID)
	if err != nil {
		return nil, err
	}
	updated.AcceptToken = rawToken
	updated.EmailSent = c.trySendInvitationEmail(ctx, updated)

	return updated, nil
}
