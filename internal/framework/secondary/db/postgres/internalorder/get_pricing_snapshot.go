package repository

import (
	"context"
	"database/sql"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
)

func (r *internalOrderRepo) GetPricingSnapshot(
	ctx context.Context,
	productID string,
	pricingID string,
	currencyCode string,
) (map[string]any, string, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalorder:get_pricing_snapshot:GetPricingSnapshot",
	)
	defer span.End()

	type row struct {
		ProductID    string          `db:"product_id"`
		ProductCode  string          `db:"product_code"`
		ProductName  string          `db:"product_name"`
		PricingID    string          `db:"pricing_id"`
		PricingCode  string          `db:"pricing_code"`
		PricingName  string          `db:"pricing_name"`
		CurrencyCode string          `db:"currency_code"`
		Amount       decimal.Decimal `db:"amount"`
	}
	var item row
	query := `
		SELECT
			p.id product_id,
			p.code product_code,
			p.name product_name,
			pp.id pricing_id,
			pp.code pricing_code,
			pp.name pricing_name,
			price.currency_code,
			price.amount
		FROM internal_products p
		JOIN internal_product_pricings pp
			ON pp.internal_product_id = p.id
			AND pp.deleted_at IS NULL
		JOIN internal_product_prices price
			ON price.internal_product_pricing_id = pp.id
			AND price.deleted_at IS NULL
		WHERE p.id = ?
			AND pp.id = ?
			AND price.currency_code = ?
			AND p.deleted_at IS NULL
			AND p.status = 'active'
			AND pp.status = 'active'
			AND price.started_at <= NOW()
			AND (
				price.ended_at IS NULL
				OR price.ended_at > NOW()
			)
		ORDER BY price.started_at DESC
		LIMIT 1
	`
	if err := r.exec(ctx).GetContext(ctx, &item, r.exec(ctx).Rebind(query), productID, pricingID, currencyCode); err != nil {
		payload := map[string]string{
			"product_id":    productID,
			"pricing_id":    pricingID,
			"currency_code": currencyCode,
		}
		if err == sql.ErrNoRows {
			log.Ctx(ctx).Warn().Any(common.LogKeyPayload, payload).Msg("Active product pricing not found")
			return nil, "", errmsg.NewCustomErrors(404).SetMessage("Active product pricing not found")
		}
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, payload).Msg("Failed to get pricing snapshot")
		return nil, "", err
	}
	return map[string]any{
		"product": map[string]any{
			"id":   item.ProductID,
			"code": item.ProductCode,
			"name": item.ProductName,
		},
		"pricing": map[string]any{
			"id":   item.PricingID,
			"code": item.PricingCode,
			"name": item.PricingName,
		},
		"price": map[string]any{
			"currency_code": item.CurrencyCode,
			"amount":        item.Amount.StringFixed(2),
		},
	}, item.Amount.String(), nil
}
