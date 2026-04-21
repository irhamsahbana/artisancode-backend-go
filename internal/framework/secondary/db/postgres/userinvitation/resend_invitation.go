package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/rs/zerolog/log"
)

func (r *userInvitationRepo) ResendInvitation(ctx context.Context, data coreentity.UserInvitation) error {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:userinvitation:resend_invitation:ResendInvitation")
	defer span.End()

	query := `
		UPDATE user_invitations
		SET token_hash = ?, expires_at = ?, last_sent_at = ?, updated_at = NOW()
		WHERE id = ? AND tenant_id = ? AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, r.db.Rebind(query), data.TokenHash, data.ExpiresAt, data.LastSentAt, data.ID, data.TenantID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to resend invitation")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errmsg.NewCustomErrors(404).SetMessage("Invitation not found")
	}

	return nil
}
