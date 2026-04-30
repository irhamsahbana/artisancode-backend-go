package repository

import (
	"context"
	"database/sql"
	"time"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *userRepo) FindActiveUsersByEmail(ctx context.Context, email string) ([]coreentity.User, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:user:find_active_users_by_email:FindActiveUsersByEmail",
	)
	defer span.End()

	query := `
		SELECT
			u.id,
			u.name,
			u.email,
			u.password,
			u.tenant_id,
			t.name AS tenant_name,
			u.username,
			u.company_id,
			c.name AS company_name,
			u.email_verified_at
		FROM users u
		JOIN tenants t ON t.id = u.tenant_id AND t.deleted_at IS NULL
		LEFT JOIN org_units c ON c.id = u.company_id AND c.category = 'company'
		WHERE LOWER(u.email) = LOWER(?) AND u.deleted_at IS NULL
		ORDER BY u.created_at ASC
	`

	var rows []struct {
		ID              string       `db:"id"`
		Name            string       `db:"name"`
		Email           string       `db:"email"`
		Password        string       `db:"password"`
		TenantID        string       `db:"tenant_id"`
		TenantName      string       `db:"tenant_name"`
		UserName        string       `db:"username"`
		CompanyID       *string      `db:"company_id"`
		CompanyName     *string      `db:"company_name"`
		EmailVerifiedAt sql.NullTime `db:"email_verified_at"`
	}

	exec := r.executor(ctx)
	if err := exec.SelectContext(ctx, &rows, exec.Rebind(query), email); err != nil {
		log.Ctx(ctx).
			Error().
			Err(err).
			Any(common.LogKeyPayload, map[string]string{"email": email}).
			Msg("Failed to find active users by email")
		return nil, err
	}

	users := make([]coreentity.User, 0, len(rows))
	for _, row := range rows {
		var emailVerifiedAt *time.Time
		if row.EmailVerifiedAt.Valid {
			emailVerifiedAt = &row.EmailVerifiedAt.Time
		}

		roleNames, err := r.GetUserRolesByUserID(ctx, row.ID)
		if err != nil {
			return nil, err
		}

		users = append(users, coreentity.User{
			ID:              row.ID,
			Name:            row.Name,
			Email:           row.Email,
			Password:        row.Password,
			TenantID:        row.TenantID,
			TenantName:      row.TenantName,
			UserName:        row.UserName,
			CompanyID:       row.CompanyID,
			CompanyName:     row.CompanyName,
			RoleNames:       roleNames,
			EmailVerifiedAt: emailVerifiedAt,
		})
	}

	return users, nil
}
