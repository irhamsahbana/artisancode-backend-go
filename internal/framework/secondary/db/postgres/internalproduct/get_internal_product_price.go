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
	"github.com/shopspring/decimal"
)

func (r *internalProductRepo) GetInternalProductPrice(
	ctx context.Context,
	filter coreentity.InternalProductPriceFilter,
) (*coreentity.InternalProductPrice, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalproduct:get_internal_product_price:GetInternalProductPrice",
	)
	defer span.End()

	type dao struct {
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

	row := new(dao)
	query := `
		SELECT
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
		WHERE id = ?
			AND deleted_at IS NULL
	`
	if err := r.db.GetContext(ctx, row, r.db.Rebind(query), filter.ID); err != nil {
		if err == sql.ErrNoRows {
			log.Ctx(ctx).Warn().Any(common.LogKeyPayload, filter).Msg("Internal product price not found")
			return nil, errmsg.NewCustomErrors(404).SetMessage(errmsg.MessageInternalProductPriceNotFound)
		}
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to query internal product price")
		return nil, err
	}

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
		return nil, err
	}

	return &item, nil
}
