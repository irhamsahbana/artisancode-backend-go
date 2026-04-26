package repository

import (
	"context"
	"encoding/json"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (r *internalOrderRepo) CreateOrder(
	ctx context.Context,
	data coreentity.InternalOrder,
) (*coreentity.InternalOrder, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalorder:create_order:CreateOrder",
	)
	defer span.End()

	number, err := r.nextNumber(ctx, "ORD")
	if err != nil {
		return nil, err
	}
	pricingSnapshot, _ := json.Marshal(data.PricingSnapshot)
	metadata, _ := json.Marshal(data.Metadata)
	query := `
		INSERT INTO internal_orders (
			tenant_id,
			company_id,
			order_number,
			status,
			currency_code,
			subtotal_amount,
			discount_amount,
			tax_amount,
			total_amount,
			internal_product_id,
			internal_product_pricing_id,
			pricing_snapshot,
			source_type,
			source_reference_id,
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
			order_number,
			ordered_at,
			created_at
	`
	if err := r.exec(ctx).QueryRowxContext(ctx, r.exec(ctx).Rebind(query),
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
		data.SourceType,
		data.SourceReferenceID,
		metadata,
	).Scan(&data.ID, &data.OrderNumber, &data.OrderedAt, &data.CreatedAt); err != nil {
		return nil, err
	}
	return &data, nil
}
