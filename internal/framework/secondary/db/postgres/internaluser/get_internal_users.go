package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"github.com/rs/zerolog/log"
)

func (r *internalUserRepo) GetInternalUsers(ctx context.Context, filter coreentity.InternalUserListFilter) ([]coreentity.InternalUser, int, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:internaluser:repo:GetInternalUsers")
	defer span.End()

	type dao struct {
		TotalData   int     `db:"total_data"`
		ID          string  `db:"id"`
		FullName    string  `db:"full_name"`
		Email       string  `db:"email"`
		RoleCode    string  `db:"role_code"`
		Status      string  `db:"status"`
		LastLoginAt *string `db:"last_login_at"`
		CreatedAt   string  `db:"created_at"`
		UpdatedAt   *string `db:"updated_at"`
	}

	var (
		rows  = make([]dao, 0)
		items = make([]coreentity.InternalUser, 0)
		args  = make([]any, 0, 4)
		total int
	)

	query := `
		SELECT COUNT(*) OVER() AS total_data,
		       id, full_name, email, role_code, status,
		       last_login_at::text AS last_login_at, created_at::text AS created_at, updated_at::text AS updated_at
		FROM internal_users
		WHERE deleted_at IS NULL
	`
	if filter.Q != "" {
		query += ` AND (full_name ILIKE '%' || ? || '%' OR email ILIKE '%' || ? || '%')`
		args = append(args, filter.Q, filter.Q)
	}
	query += ` ORDER BY full_name ASC LIMIT ? OFFSET ?`
	args = append(args, filter.Paginate, (filter.Page-1)*filter.Paginate)

	if err := r.db.SelectContext(ctx, &rows, r.db.Rebind(query), args...); err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to query internal users")
		return nil, 0, err
	}

	for _, row := range rows {
		total = row.TotalData
		item := coreentity.InternalUser{
			ID:          row.ID,
			FullName:    row.FullName,
			Email:       row.Email,
			RoleCode:    row.RoleCode,
			Status:      row.Status,
			LastLoginAt: row.LastLoginAt,
			CreatedAt:   row.CreatedAt,
		}
		if row.UpdatedAt != nil {
			item.UpdatedAt = *row.UpdatedAt
		}
		items = append(items, item)
	}

	return items, total, nil
}
