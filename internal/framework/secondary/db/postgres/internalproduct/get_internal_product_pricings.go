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

func (r *internalProductRepo) GetInternalProductPricings(
	ctx context.Context,
	filter coreentity.InternalProductPricingListFilter,
) ([]coreentity.InternalProductPricing, int, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalproduct:repo:GetInternalProductPricings",
	)
	defer span.End()

	type dao struct {
		TotalData         int             `db:"total_data"`
		ID                string          `db:"id"`
		InternalProductID string          `db:"internal_product_id"`
		Code              string          `db:"code"`
		Name              string          `db:"name"`
		Description       string          `db:"description"`
		Status            string          `db:"status"`
		Metadata          json.RawMessage `db:"metadata"`
		CreatedAt         string          `db:"created_at"`
		UpdatedAt         *string         `db:"updated_at"`
	}

	var (
		rows  = make([]dao, 0)
		items = make([]coreentity.InternalProductPricing, 0)
		args  = []any{filter.InternalProductID}
		total int
	)

	query := `
		SELECT
			COUNT(*) OVER() AS total_data,
			id,
			internal_product_id,
			code,
			name,
			description,
			status,
			metadata,
			created_at,
			updated_at
		FROM internal_product_pricings
		WHERE internal_product_id = ?
			AND deleted_at IS NULL
	`
	if filter.Q != "" {
		query += ` AND (name ILIKE '%' || ? || '%' OR code ILIKE '%' || ? || '%')`
		args = append(args, filter.Q, filter.Q)
	}
	query += ` ORDER BY name ASC LIMIT ? OFFSET ?`
	args = append(args, filter.Paginate, (filter.Page-1)*filter.Paginate)

	if err := r.db.SelectContext(ctx, &rows, r.db.Rebind(query), args...); err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to query internal product pricings")
		return nil, 0, err
	}

	for _, row := range rows {
		total = row.TotalData
		item, err := mapInternalProductPricingDAO(
			row.ID,
			row.InternalProductID,
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

func (r *internalProductRepo) GetInternalProductPricing(
	ctx context.Context,
	filter coreentity.InternalProductPricingFilter,
) (*coreentity.InternalProductPricing, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalproduct:repo:GetInternalProductPricing",
	)
	defer span.End()

	type dao struct {
		ID                string          `db:"id"`
		InternalProductID string          `db:"internal_product_id"`
		Code              string          `db:"code"`
		Name              string          `db:"name"`
		Description       string          `db:"description"`
		Status            string          `db:"status"`
		Metadata          json.RawMessage `db:"metadata"`
		CreatedAt         string          `db:"created_at"`
		UpdatedAt         *string         `db:"updated_at"`
	}

	var row dao
	query := `
		SELECT
			id,
			internal_product_id,
			code,
			name,
			description,
			status,
			metadata,
			created_at,
			updated_at
		FROM internal_product_pricings
		WHERE id = ?
			AND deleted_at IS NULL
	`
	if err := r.db.GetContext(ctx, &row, r.db.Rebind(query), filter.ID); err != nil {
		if err == sql.ErrNoRows {
			log.Ctx(ctx).Warn().Any(common.LogKeyPayload, filter).Msg(errmsg.MessageInternalProductPricingNotFound)
			return nil, errmsg.NewCustomErrors(404).SetMessage(errmsg.MessageInternalProductPricingNotFound)
		}
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to get internal product pricing")
		return nil, err
	}

	item, err := mapInternalProductPricingDAO(
		row.ID,
		row.InternalProductID,
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
