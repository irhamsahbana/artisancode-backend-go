package repository

import (
	"context"
	"database/sql"
	"encoding/json"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/rs/zerolog/log"
)

func (r *internalProductRepo) GetInternalProducts(
	ctx context.Context,
	filter coreentity.InternalProductListFilter,
) ([]coreentity.InternalProduct, int, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalproduct:repo:GetInternalProducts",
	)
	defer span.End()

	type dao struct {
		TotalData   int             `db:"total_data"`
		ID          string          `db:"id"`
		Code        string          `db:"code"`
		Name        string          `db:"name"`
		Description string          `db:"description"`
		Status      string          `db:"status"`
		Metadata    json.RawMessage `db:"metadata"`
		CreatedAt   string          `db:"created_at"`
		UpdatedAt   *string         `db:"updated_at"`
	}

	var (
		rows  = make([]dao, 0)
		items = make([]coreentity.InternalProduct, 0)
		args  = make([]any, 0, 4)
		total int
	)

	query := `
		SELECT
			COUNT(*) OVER() AS total_data,
			id,
			code,
			name,
			description,
			status,
			metadata,
			created_at,
			updated_at
		FROM internal_products
		WHERE deleted_at IS NULL
	`
	if filter.Q != "" {
		query += ` AND (name ILIKE '%' || ? || '%' OR code ILIKE '%' || ? || '%')`
		args = append(args, filter.Q, filter.Q)
	}
	if filter.Status != "" {
		query += ` AND status = ?`
		args = append(args, filter.Status)
	}
	query += ` ORDER BY name ASC LIMIT ? OFFSET ?`
	args = append(args, filter.Paginate, (filter.Page-1)*filter.Paginate)

	if err := r.db.SelectContext(ctx, &rows, r.db.Rebind(query), args...); err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to query internal products")
		return nil, 0, err
	}

	for _, row := range rows {
		total = row.TotalData
		item, err := mapInternalProductDAO(
			row.ID,
			row.Code,
			row.Name,
			row.Description,
			row.Status,
			row.Metadata,
			row.CreatedAt,
			row.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}

	return items, total, nil
}

func (r *internalProductRepo) GetInternalProduct(
	ctx context.Context,
	filter coreentity.InternalProductFilter,
) (*coreentity.InternalProduct, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalproduct:repo:GetInternalProduct",
	)
	defer span.End()

	type dao struct {
		ID          string          `db:"id"`
		Code        string          `db:"code"`
		Name        string          `db:"name"`
		Description string          `db:"description"`
		Status      string          `db:"status"`
		Metadata    json.RawMessage `db:"metadata"`
		CreatedAt   string          `db:"created_at"`
		UpdatedAt   *string         `db:"updated_at"`
	}

	var row dao
	query := `
		SELECT
			id,
			code,
			name,
			description,
			status,
			metadata,
			created_at,
			updated_at
		FROM internal_products
		WHERE id = ?
			AND deleted_at IS NULL
	`
	if err := r.db.GetContext(ctx, &row, r.db.Rebind(query), filter.ID); err != nil {
		if err == sql.ErrNoRows {
			log.Ctx(ctx).Warn().Any(common.LogKeyPayload, filter).Msg(errmsg.MessageInternalProductNotFound)
			return nil, errmsg.NewCustomErrors(404).SetMessage(errmsg.MessageInternalProductNotFound)
		}
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to get internal product")
		return nil, err
	}

	item, err := mapInternalProductDAO(
		row.ID,
		row.Code,
		row.Name,
		row.Description,
		row.Status,
		row.Metadata,
		row.CreatedAt,
		row.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &item, nil
}
