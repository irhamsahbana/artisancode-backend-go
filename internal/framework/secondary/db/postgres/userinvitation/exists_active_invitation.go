package repository

import (
	"context"
	"database/sql"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *userInvitationRepo) ExistsActiveInvitation(ctx context.Context, tenantID, email, roleCode string, employeeID *string) (bool, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:userinvitation:exists_active_invitation:ExistsActiveInvitation")
	defer span.End()

	query := `
		SELECT EXISTS(
			SELECT 1
			FROM user_invitations
			WHERE tenant_id = ?
				AND email = ?
				AND role_code = ?
				AND status = 'pending'
				AND expires_at >= NOW()
				AND deleted_at IS NULL
				AND (
					(employee_id IS NULL AND ? IS NULL)
					OR employee_id = ?
				)
		)
	`

	var exists bool
	if err := r.db.GetContext(ctx, &exists, r.db.Rebind(query), tenantID, email, roleCode, employeeID, employeeID); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, map[string]any{
			"tenant_id":   tenantID,
			"email":       email,
			"role_code":   roleCode,
			"employee_id": employeeID,
		}).Msg("Failed to check active invitation")
		return false, err
	}

	return exists, nil
}
