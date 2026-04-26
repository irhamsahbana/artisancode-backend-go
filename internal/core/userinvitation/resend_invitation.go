package core

import (
	"context"
	"time"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
)

func (c *userInvitationCore) ResendInvitation(
	ctx context.Context,
	data coreentity.UserInvitation,
) (*coreentity.UserInvitation, error) {
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

	var updated *coreentity.UserInvitation
	err = c.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := c.repo.ResendInvitation(txCtx, *item); err != nil {
			return err
		}

		updatedItem, err := c.repo.GetInvitationByID(txCtx, data.UserCtx.TenantID, data.ID)
		if err != nil {
			return err
		}
		updatedItem.AcceptToken = rawToken
		emailSent, err := c.trySendInvitationEmail(txCtx, updatedItem)
		if err != nil {
			return err
		}
		updatedItem.EmailSent = emailSent
		updated = updatedItem
		return nil
	})
	if err != nil {
		return nil, err
	}

	return updated, nil
}
