package repository

import (
	"context"
	"database/sql"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/rs/zerolog/log"
)

func (r *userInvitationRepo) GetInvitationByID(ctx context.Context, tenantID, invitationID string) (*coreentity.UserInvitation, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:userinvitation:get_invitation:GetInvitationByID")
	defer span.End()

	query := baseInvitationSelectQuery() + `
		AND ui.id = ? AND ui.tenant_id = ?
	`

	item, err := r.getInvitation(ctx, query, invitationID, tenantID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errmsg.NewCustomErrors(404).SetMessage("Invitation not found")
		}
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, map[string]string{
			"tenant_id":     tenantID,
			"invitation_id": invitationID,
		}).Msg("Failed to get invitation by id")
		return nil, err
	}

	return item, nil
}

func (r *userInvitationRepo) GetInvitationByTokenHash(ctx context.Context, tokenHash string) (*coreentity.UserInvitation, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:userinvitation:get_invitation:GetInvitationByTokenHash")
	defer span.End()

	query := baseInvitationSelectQuery() + `
		AND ui.token_hash = ?
	`

	item, err := r.getInvitation(ctx, query, tokenHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errmsg.NewCustomErrors(404).SetMessage("Invitation not found")
		}
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, map[string]string{
			"token_hash": tokenHash,
		}).Msg("Failed to get invitation by token hash")
		return nil, err
	}

	return item, nil
}

func baseInvitationSelectQuery() string {
	return `
		SELECT
			ui.id,
			ui.tenant_id,
			t.code AS tenant_code,
			t.name AS tenant_name,
			ui.employee_id,
			e.employee_no,
			e.full_name AS employee_name,
			ui.email,
			ui.role_code,
			ui.status,
			ui.token_hash,
			ui.expires_at,
			ui.accepted_at,
			ui.revoked_at,
			ui.invited_by,
			ui.last_sent_at,
			ui.created_at,
			ui.updated_at
		FROM user_invitations ui
		INNER JOIN tenants t ON t.id = ui.tenant_id
		LEFT JOIN employees e ON e.id = ui.employee_id AND e.deleted_at IS NULL
		WHERE ui.deleted_at IS NULL
	`
}

func (r *userInvitationRepo) getInvitation(ctx context.Context, query string, args ...any) (*coreentity.UserInvitation, error) {
	type dao struct {
		ID           string         `db:"id"`
		TenantID     string         `db:"tenant_id"`
		TenantCode   string         `db:"tenant_code"`
		TenantName   string         `db:"tenant_name"`
		EmployeeID   *string        `db:"employee_id"`
		EmployeeNo   sql.NullString `db:"employee_no"`
		EmployeeName sql.NullString `db:"employee_name"`
		Email        string         `db:"email"`
		RoleCode     string         `db:"role_code"`
		Status       string         `db:"status"`
		TokenHash    string         `db:"token_hash"`
		ExpiresAt    sql.NullTime   `db:"expires_at"`
		AcceptedAt   sql.NullTime   `db:"accepted_at"`
		RevokedAt    sql.NullTime   `db:"revoked_at"`
		InvitedBy    string         `db:"invited_by"`
		LastSentAt   sql.NullTime   `db:"last_sent_at"`
		CreatedAt    sql.NullTime   `db:"created_at"`
		UpdatedAt    sql.NullTime   `db:"updated_at"`
	}

	var row dao
	if err := r.db.GetContext(ctx, &row, r.db.Rebind(query), args...); err != nil {
		return nil, err
	}

	item := &coreentity.UserInvitation{
		ID:         row.ID,
		TenantID:   row.TenantID,
		TenantCode: row.TenantCode,
		TenantName: row.TenantName,
		EmployeeID: row.EmployeeID,
		Email:      row.Email,
		RoleCode:   row.RoleCode,
		Status:     row.Status,
		TokenHash:  row.TokenHash,
		InvitedBy:  row.InvitedBy,
	}

	if row.EmployeeNo.Valid {
		value := row.EmployeeNo.String
		item.EmployeeNo = &value
	}
	if row.EmployeeName.Valid {
		value := row.EmployeeName.String
		item.EmployeeName = &value
	}
	if row.ExpiresAt.Valid {
		item.ExpiresAt = row.ExpiresAt.Time
	}
	if row.AcceptedAt.Valid {
		value := row.AcceptedAt.Time
		item.AcceptedAt = &value
	}
	if row.RevokedAt.Valid {
		value := row.RevokedAt.Time
		item.RevokedAt = &value
	}
	if row.LastSentAt.Valid {
		item.LastSentAt = row.LastSentAt.Time
	}
	if row.CreatedAt.Valid {
		item.CreatedAt = row.CreatedAt.Time
	}
	if row.UpdatedAt.Valid {
		value := row.UpdatedAt.Time
		item.UpdatedAt = &value
	}

	return item, nil
}
