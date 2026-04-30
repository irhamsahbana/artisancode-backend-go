package repository

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
	"github.com/rs/zerolog/log"
)

func (r *userRepo) GetUserByEmailAndTenant(ctx context.Context, email, tenantCode string) (*coreentity.User, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:user:login:GetUserByEmailAndTenant")
	defer span.End()

	tenantCode = strings.ToUpper(tenantCode)
	payload := map[string]string{"email": email, "tenantCode": tenantCode}

	query := `
		SELECT
			u.id,
			u.email,
			u.password,
			u.tenant_id,
			t.name as tenant_name,
			u.username,
			u.company_id,
			c.name as company_name,
			u.email_verified_at
		FROM
			users u
		JOIN tenants t ON t.id = u.tenant_id AND t.deleted_at IS NULL
		LEFT JOIN org_units c ON c.id = u.company_id AND c.category = 'company'
		WHERE
			u.email = ?
			AND t.code = ?
			AND u.deleted_at IS NULL
	`

	var row struct {
		ID              string       `db:"id"`
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
	err := exec.GetContext(ctx, &row, exec.Rebind(query), email, tenantCode)
	if err != nil {
		if err == sql.ErrNoRows {
			tracing.RecordError(span, err)
			log.Ctx(ctx).Warn().Any(common.LogKeyPayload, payload).Msg("User not found")
			return nil, errmsg.NewCustomErrors(400).SetMessage("User not found")
		}
		tracing.RecordError(span, err)
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, payload).Msg("Failed to get user details")
		return nil, err
	}

	var emailVerifiedAt *time.Time
	if row.EmailVerifiedAt.Valid {
		emailVerifiedAt = &row.EmailVerifiedAt.Time
	}

	return &coreentity.User{
		ID:              row.ID,
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
