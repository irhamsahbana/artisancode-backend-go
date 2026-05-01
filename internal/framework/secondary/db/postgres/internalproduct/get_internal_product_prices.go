package repository

import (
	"context"
	"encoding/json"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
)

func (r *internalProductRepo) GetInternalProductPrices(
	ctx context.Context,
	filter coreentity.InternalProductPriceListFilter,
) ([]coreentity.InternalProductPrice, int, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalproduct:repo:GetInternalProductPrices",
	)
	defer span.End()

	type dao struct {
		TotalData                int             `db:"total_data"`
		ID                       string          `db:"id"`
		InternalProductPricingID string          `db:"internal_product_pricing_id"`
		CurrencyCode             string          `db:"currency_code"`
		Amount                   decimal.Decimal `db:"amount"`
		StartedAt                string          `db:"started_at"`
		EndedAt                  *string         `db:"ended_at"`
		Metadata                 json.RawMessage `db:"metadata"`
		CreatedAt                string          `db:"created_at"`
		UpdatedAt                *string         `db:"updated_at"`
	}

	var (
		rows  = make([]dao, 0)
		items = make([]coreentity.InternalProductPrice, 0)
		args  = []any{filter.InternalProductPricingID}
		total int
	)

	query := `
		SELECT
			COUNT(*) OVER() AS total_data,
			id,
			internal_product_pricing_id,
			currency_code,
			amount::text AS amount,
			started_at::text AS started_at,
			ended_at::text AS ended_at,
			metadata,
			created_at::text AS created_at,
			updated_at::text AS updated_at
		FROM internal_product_prices
		WHERE internal_product_pricing_id = ?
			AND deleted_at IS NULL
	`
	if filter.CurrencyCode != "" {
		query += ` AND currency_code = ?`
		args = append(args, filter.CurrencyCode)
	}
	query += ` ORDER BY currency_code ASC, started_at DESC LIMIT ? OFFSET ?`
	args = append(args, filter.Paginate, (filter.Page-1)*filter.Paginate)

	if err := r.db.SelectContext(ctx, &rows, r.db.Rebind(query), args...); err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to query internal product prices")
		return nil, 0, err
	}

	for _, row := range rows {
		total = row.TotalData
		item, err := mapInternalProductPriceDAO(
			row.ID,
			row.InternalProductPricingID,
			row.CurrencyCode,
			row.Amount,
			row.StartedAt,
			row.EndedAt,
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
