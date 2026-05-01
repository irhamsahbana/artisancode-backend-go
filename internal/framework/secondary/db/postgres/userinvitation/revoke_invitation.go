package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/rs/zerolog/log"
)

func (r *userInvitationRepo) RevokeInvitation(ctx context.Context, tenantID, invitationID string) error {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:userinvitation:revoke_invitation:RevokeInvitation",
	)
	defer span.End()

	query := `
		UPDATE user_invitations
		SET status = 'revoked', revoked_at = NOW(), updated_at = NOW()
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	payload := map[string]string{
		"tenant_id":     tenantID,
		"invitation_id": invitationID,
	}
	result, err := r.db.ExecContext(ctx, r.db.Rebind(query), invitationID, tenantID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, payload).Msg("Failed to revoke invitation")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, payload).Msg("Failed to get revoke invitation rows affected")
		return err
	}
	if rowsAffected == 0 {
		log.Ctx(ctx).Warn().Any(common.LogKeyPayload, payload).Msg("Invitation not found when revoking")
		return errmsg.NewCustomErrors(404).SetMessage(errmsg.MessageInvitationNotFound)
	}

	return nil
}
