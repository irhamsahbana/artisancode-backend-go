package repository

import (
	"context"
	"database/sql"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

func (r *userRepo) GetUser(ctx context.Context, filter coreentity.User) (*coreentity.User, error) {
	ctx, span := tracing.StartSpan(ctx, "repo.GetUser")
	defer span.End()

	payload := map[string]string{
		"id":        filter.ID,
		"tenant_id": filter.TenantID,
	}

	var row struct {
		ID          string         `db:"id"`
		Name        string         `db:"name"`
		UserName    string         `db:"username"`
		Email       string         `db:"email"`
		CompanyID   *string        `db:"company_id"`
		CompanyName *string        `db:"company_name"`
		RoleIDs     pq.StringArray `db:"role_ids"`
		RoleNames   pq.StringArray `db:"role_names"`
	}

	query := `
		SELECT
			u.id,
			u.name,
			u.username,
			u.email,
			u.company_id,
			c.name AS company_name,
			COALESCE(array_agg(DISTINCT ur.role_id::text) FILTER (WHERE ur.role_id IS NOT NULL), '{}') AS role_ids,
			COALESCE(array_agg(DISTINCT r.name) FILTER (WHERE r.name IS NOT NULL), '{}') AS role_names
		FROM users u
		LEFT JOIN org_units c ON c.id = u.company_id AND c.deleted_at IS NULL
		LEFT JOIN user_roles ur ON ur.user_id = u.id
		LEFT JOIN roles r ON r.id = ur.role_id AND r.deleted_at IS NULL
		WHERE u.id = ? AND u.tenant_id = ? AND u.deleted_at IS NULL
		GROUP BY u.id, u.name, u.username, u.email, u.company_id, c.name
	`

	if err := r.db.GetContext(ctx, &row, r.db.Rebind(query), filter.ID, filter.TenantID); err != nil {
		if err == sql.ErrNoRows {
			log.Ctx(ctx).Warn().Any(common.LogKeyPayload, payload).Msg("User not found")
			return nil, errmsg.NewCustomErrors(404).SetMessage("User not found")
		}
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, payload).Msg("Failed to get user")
		return nil, err
	}

	return &coreentity.User{
		ID:          row.ID,
		Name:        row.Name,
		UserName:    row.UserName,
		Email:       row.Email,
		CompanyID:   row.CompanyID,
		CompanyName: row.CompanyName,
		RoleIDs:     []string(row.RoleIDs),
		RoleNames:   []string(row.RoleNames),
		TenantID:    filter.TenantID,
	}, nil
}
