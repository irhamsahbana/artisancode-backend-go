package repository

import (
	"context"
	"encoding/json"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *internalQuotationRepo) CreateQuotation(
	ctx context.Context,
	data coreentity.InternalQuotation,
) (*coreentity.InternalQuotation, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalquotation:create_quotation:CreateQuotation",
	)
	defer span.End()

	exec := r.exec(ctx)
	number, err := r.nextNumber(ctx, "QUO")
	if err != nil {
		return nil, err
	}
	pricingSnapshot, _ := json.Marshal(data.PricingSnapshot)
	quoteSnapshot, _ := json.Marshal(data.QuoteSnapshot)
	metadata, _ := json.Marshal(data.Metadata)
	query := `
		INSERT INTO internal_quotations (
			tenant_id,
			company_id,
			quotation_number,
			status,
			currency_code,
			subtotal_amount,
			discount_amount,
			tax_amount,
			total_amount,
			internal_product_id,
			internal_product_pricing_id,
			pricing_snapshot,
			quote_snapshot,
			expires_at,
			metadata
		)
		VALUES (
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?
		)
		RETURNING
			id,
			quotation_number,
			status,
			created_at
	`
	if err := exec.QueryRowxContext(ctx, exec.Rebind(query),
		data.TenantID,
		data.CompanyID,
		number,
		data.Status,
		data.CurrencyCode,
		data.SubtotalAmount,
		data.DiscountAmount,
		data.TaxAmount,
		data.TotalAmount,
		data.InternalProductID,
		data.InternalProductPricingID,
		pricingSnapshot,
		quoteSnapshot,
		data.ExpiresAt,
		metadata,
	).Scan(&data.ID, &data.QuotationNumber, &data.Status, &data.CreatedAt); err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to create internal quotation")
		return nil, err
	}
	return &data, nil
}
