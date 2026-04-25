package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *userInvitationRepo) CreateInvitation(ctx context.Context, data coreentity.UserInvitation) (*coreentity.UserInvitation, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:userinvitation:create_invitation:CreateInvitation")
	defer span.End()

	query := `
		INSERT INTO user_invitations (
			tenant_id, employee_id, email, role_code, status, token_hash, expires_at, invited_by, last_sent_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id, created_at
	`

	var created struct {
		ID        string `db:"id"`
		CreatedAt string `db:"created_at"`
	}

	exec := r.executor(ctx)
	err := exec.GetContext(ctx, &created, exec.Rebind(query),
		data.TenantID,
		data.EmployeeID,
		data.Email,
		data.RoleCode,
		data.Status,
		data.TokenHash,
		data.ExpiresAt,
		data.InvitedBy,
		data.LastSentAt,
	)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to create user invitation")
		return nil, err
	}

	return r.GetInvitationByID(ctx, data.TenantID, created.ID)
}
