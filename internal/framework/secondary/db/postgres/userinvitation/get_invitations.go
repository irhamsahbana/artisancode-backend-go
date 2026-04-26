package repository

import (
	"context"
	"database/sql"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (r *userInvitationRepo) GetInvitations(
	ctx context.Context,
	filter coreentity.UserInvitationListFilter,
) ([]coreentity.UserInvitation, int, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:userinvitation:get_invitations:GetInvitations",
	)
	defer span.End()

	type dao struct {
		TotalData    int            `db:"total_data"`
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
		ExpiresAt    timeNull       `db:"expires_at"`
		AcceptedAt   timeNull       `db:"accepted_at"`
		RevokedAt    timeNull       `db:"revoked_at"`
		LastSentAt   timeNull       `db:"last_sent_at"`
		CreatedAt    timeNull       `db:"created_at"`
	}

	var (
		rows  = make([]dao, 0)
		items = make([]coreentity.UserInvitation, 0)
		total = 0
		args  = make([]any, 0, 8)
	)

	query := `
		SELECT
			COUNT(*) OVER() AS total_data,
			ui.id,
			ui.tenant_id,
			t.code AS tenant_code,
			t.name AS tenant_name,
			ui.employee_id,
			e.employee_no,
			e.full_name AS employee_name,
			ui.email,
			ui.role_code,
			CASE
				WHEN ui.status = 'pending' AND ui.expires_at < NOW() THEN 'expired'
				ELSE ui.status
			END AS status,
			ui.expires_at,
			ui.accepted_at,
			ui.revoked_at,
			ui.last_sent_at,
			ui.created_at
		FROM user_invitations ui
		INNER JOIN tenants t ON t.id = ui.tenant_id
		LEFT JOIN employees e ON e.id = ui.employee_id AND e.deleted_at IS NULL
		WHERE ui.deleted_at IS NULL AND ui.tenant_id = ?
	`
	args = append(args, filter.TenantID)

	if filter.RoleCode != "" {
		query += ` AND ui.role_code = ?`
		args = append(args, filter.RoleCode)
	}
	if len(filter.EmployeeIDs) == 1 {
		query += ` AND ui.employee_id = ?`
		args = append(args, filter.EmployeeIDs[0])
	}
	if len(filter.EmployeeIDs) > 1 {
		inQuery, inArgs, err := sqlx.In(` AND ui.employee_id IN (?)`, filter.EmployeeIDs)
		if err != nil {
			log.Ctx(ctx).
				Error().
				Err(err).
				Any(common.LogKeyPayload, filter.EmployeeIDs).
				Msg("Failed to build employee invitation filter")
			return nil, 0, err
		}
		query += inQuery
		args = append(args, inArgs...)
	}
	if filter.Status != "" {
		if filter.Status == coreentity.UserInvitationStatusExpired {
			query += ` AND ui.status = 'pending' AND ui.expires_at < NOW()`
		} else {
			query += ` AND ui.status = ?`
			args = append(args, filter.Status)
		}
	}
	if filter.Q != "" {
		query += `
			AND (
				ui.email ILIKE '%' || ? || '%'
				OR COALESCE(e.full_name, '') ILIKE '%' || ? || '%'
				OR COALESCE(e.employee_no, '') ILIKE '%' || ? || '%'
			)
		`
		args = append(args, filter.Q, filter.Q, filter.Q)
	}

	query += ` ORDER BY ui.created_at DESC LIMIT ? OFFSET ?`
	args = append(args, filter.Paginate, (filter.Page-1)*filter.Paginate)

	if err := r.db.SelectContext(ctx, &rows, r.db.Rebind(query), args...); err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to query invitations")
		return nil, 0, err
	}

	for _, row := range rows {
		total = row.TotalData
		item := coreentity.UserInvitation{
			ID:         row.ID,
			TenantID:   row.TenantID,
			TenantCode: row.TenantCode,
			TenantName: row.TenantName,
			EmployeeID: row.EmployeeID,
			Email:      row.Email,
			RoleCode:   row.RoleCode,
			Status:     row.Status,
			ExpiresAt:  row.ExpiresAt.Time,
			LastSentAt: row.LastSentAt.Time,
			CreatedAt:  row.CreatedAt.Time,
		}

		if row.EmployeeNo.Valid {
			value := row.EmployeeNo.String
			item.EmployeeNo = &value
		}
		if row.EmployeeName.Valid {
			value := row.EmployeeName.String
			item.EmployeeName = &value
		}
		if row.AcceptedAt.Valid {
			value := row.AcceptedAt.Time
			item.AcceptedAt = &value
		}
		if row.RevokedAt.Valid {
			value := row.RevokedAt.Time
			item.RevokedAt = &value
		}

		items = append(items, item)
	}

	return items, total, nil
}

type timeNull = sql.NullTime
