package repository

import (
	"context"
	"database/sql"
	"time"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/rs/zerolog/log"
)

func (r *userRepo) FindActiveUserByEmailAndTenantID(ctx context.Context, email, tenantID string) (*coreentity.User, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:user:find_active_user_by_email:FindActiveUserByEmailAndTenantID")
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
		JOIN tenants t ON t.id = u.tenant_id
		LEFT JOIN org_units c ON c.id = u.company_id AND c.category = 'company'
		WHERE u.email = ? AND u.tenant_id = ? AND u.deleted_at IS NULL
		LIMIT 1
	`

	var row struct {
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

	err := r.db.GetContext(ctx, &row, r.db.Rebind(query), email, tenantID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Ctx(ctx).Warn().Any(common.LogKeyPayload, map[string]string{
				"email":     email,
				"tenant_id": tenantID,
			}).Msg("User not found")
			return nil, errmsg.NewCustomErrors(400).SetMessage("User not found")
		}

		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, map[string]string{
			"email":     email,
			"tenant_id": tenantID,
		}).Msg("Failed to get user by email and tenant")
		return nil, err
	}

	var emailVerifiedAt *time.Time
	if row.EmailVerifiedAt.Valid {
		emailVerifiedAt = &row.EmailVerifiedAt.Time
	}

	return &coreentity.User{
		ID:              row.ID,
		Name:            row.Name,
		Email:           row.Email,
		Password:        row.Password,
		TenantID:        row.TenantID,
		TenantName:      row.TenantName,
		UserName:        row.UserName,
		CompanyID:       row.CompanyID,
		CompanyName:     row.CompanyName,
		EmailVerifiedAt: emailVerifiedAt,
	}, nil
}
