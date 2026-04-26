package repository

import (
	"context"

	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
)

func (r *userInvitationRepo) MarkInvitationAccepted(ctx context.Context, invitationID string) error {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:userinvitation:mark_invitation_accepted:MarkInvitationAccepted",
	)
	defer span.End()

	query := `
		UPDATE user_invitations
		SET status = 'accepted', accepted_at = NOW(), updated_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`

	exec := r.executor(ctx)
	result, err := exec.ExecContext(ctx, exec.Rebind(query), invitationID)
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
