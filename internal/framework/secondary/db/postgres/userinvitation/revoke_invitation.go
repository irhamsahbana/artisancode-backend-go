package repository

import (
	"context"

	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
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

	result, err := r.db.ExecContext(ctx, r.db.Rebind(query), invitationID, tenantID)
	if err != nil {
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
