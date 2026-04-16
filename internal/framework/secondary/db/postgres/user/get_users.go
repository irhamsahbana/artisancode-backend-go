package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *userRepo) GetUsers(ctx context.Context, filter coreentity.UserListFilter) ([]coreentity.User, int, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:user:get_users:GetUsers")
	defer span.End()

	type row struct {
		TotalData    int     `db:"total_data"`
		ID           string  `db:"id"`
		Name         string  `db:"name"`
		UserName     string  `db:"username"`
		Email        string  `db:"email"`
		CompanyID    *string `db:"company_id"`
		CompanyName  *string `db:"company_name"`
		RoleIDsRaw   string  `db:"role_ids"`
		RoleNamesRaw string  `db:"role_names"`
	}

	rows := make([]row, 0)
	args := make([]any, 0, 4)
	query := `
		SELECT
			COUNT(*) OVER() AS total_data,
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
		WHERE u.tenant_id = ? AND u.deleted_at IS NULL
	`
	args = append(args, filter.TenantID)

	if filter.Q != "" {
		query += ` AND (u.name ILIKE '%' || ? || '%' OR u.username ILIKE '%' || ? || '%' OR u.email ILIKE '%' || ? || '%')`
		args = append(args, filter.Q, filter.Q, filter.Q)
	}

	query += `
		GROUP BY u.id, u.name, u.username, u.email, u.company_id, c.name
		ORDER BY u.name ASC, u.email ASC
		LIMIT ? OFFSET ?
	`
	args = append(args, filter.Paginate, (filter.Page-1)*filter.Paginate)

	if err := r.db.SelectContext(ctx, &rows, r.db.Rebind(query), args...); err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to query users")
		return nil, 0, err
	}

	items := make([]coreentity.User, 0, len(rows))
	total := 0
	for _, item := range rows {
		roleIDs, err := parsePostgresTextArray(item.RoleIDsRaw)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to parse user role ids")
			return nil, 0, err
		}

		roleNames, err := parsePostgresTextArray(item.RoleNamesRaw)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to parse user role names")
			return nil, 0, err
		}

		total = item.TotalData
		items = append(items, coreentity.User{
			ID:          item.ID,
			Name:        item.Name,
			UserName:    item.UserName,
			Email:       item.Email,
			CompanyID:   item.CompanyID,
			CompanyName: item.CompanyName,
			RoleIDs:     roleIDs,
			RoleNames:   roleNames,
			TenantID:    filter.TenantID,
		})
	}

	return items, total, nil
}
