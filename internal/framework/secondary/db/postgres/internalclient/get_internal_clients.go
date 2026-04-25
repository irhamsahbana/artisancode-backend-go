package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *internalClientRepo) GetInternalClients(ctx context.Context, filter coreentity.InternalClientListFilter) ([]coreentity.InternalClient, int, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:framework:secondary:db:postgres:internalclient:get_internal_clients:GetInternalClients")
	defer span.End()

	type dao struct {
		TotalData   int     `db:"total_data"`
		ID          string  `db:"id"`
		Name        string  `db:"name"`
		Code        string  `db:"code"`
		OwnerNames  *string `db:"owner_names"`
		OwnerEmails *string `db:"owner_emails"`
		CreatedAt   string  `db:"created_at"`
		UpdatedAt   *string `db:"updated_at"`
	}

	var (
		rows  = make([]dao, 0)
		items = make([]coreentity.InternalClient, 0)
		args  = make([]any, 0, 6)
		total int
	)

	query := `
		WITH tenant_owner_rows AS (
			SELECT
				t.id,
				t.name,
				t.code,
				t.created_at,
				t.updated_at,
				u.name AS owner_name,
				u.email AS owner_email
			FROM tenants t
			LEFT JOIN roles r
				ON r.tenant_id = t.id
				AND r.name = 'owner'
				AND r.deleted_at IS NULL
			LEFT JOIN user_roles ur
				ON ur.role_id = r.id
			LEFT JOIN users u
				ON u.id = ur.user_id
				AND u.deleted_at IS NULL
			WHERE t.deleted_at IS NULL
		)
		SELECT
			COUNT(*) OVER() AS total_data,
			id,
			name,
			code,
			NULLIF(STRING_AGG(DISTINCT owner_name, ', '), '') AS owner_names,
			NULLIF(STRING_AGG(DISTINCT owner_email, ', '), '') AS owner_emails,
			created_at,
			updated_at
		FROM tenant_owner_rows
		WHERE 1 = 1
	`

	if filter.Q != "" {
		query += ` AND (name ILIKE '%' || ? || '%' OR code ILIKE '%' || ? || '%' OR owner_name ILIKE '%' || ? || '%' OR owner_email ILIKE '%' || ? || '%')`
		args = append(args, filter.Q, filter.Q, filter.Q, filter.Q)
	}

	query += `
		GROUP BY id, name, code, created_at, updated_at
		ORDER BY name ASC
		LIMIT ? OFFSET ?
	`
	args = append(args, filter.Paginate, (filter.Page-1)*filter.Paginate)

	if err := r.db.SelectContext(ctx, &rows, r.db.Rebind(query), args...); err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to query internal clients")
		return nil, 0, err
	}

	for _, row := range rows {
		total = row.TotalData

		item := coreentity.InternalClient{
			ID:          row.ID,
			Name:        row.Name,
			Code:        row.Code,
			OwnerNames:  "-",
			OwnerEmails: "-",
			CreatedAt:   row.CreatedAt,
		}
		if row.OwnerNames != nil && *row.OwnerNames != "" {
			item.OwnerNames = *row.OwnerNames
		}
		if row.OwnerEmails != nil && *row.OwnerEmails != "" {
			item.OwnerEmails = *row.OwnerEmails
		}
		if row.UpdatedAt != nil {
			item.UpdatedAt = *row.UpdatedAt
		}

		items = append(items, item)
	}

	return items, total, nil
}
