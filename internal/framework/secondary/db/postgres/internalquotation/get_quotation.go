package repository

import (
	"context"
	"database/sql"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
)

func (r *internalQuotationRepo) GetQuotation(
	ctx context.Context,
	tenantID string,
	id string,
) (*coreentity.InternalQuotation, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalquotation:get_quotation:GetQuotation",
	)
	defer span.End()

	var item quotationRow
	query := `
		SELECT
			id,
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
			approved_at,
			converted_to_order_id,
			metadata,
			created_at,
			updated_at
		FROM internal_quotations
		WHERE id = ?
			AND tenant_id = ?
			AND deleted_at IS NULL
	`
	if err := r.exec(ctx).GetContext(ctx, &item, r.exec(ctx).Rebind(query), id, tenantID); err != nil {
		if err == sql.ErrNoRows {
			return nil, errmsg.NewCustomErrors(404).SetMessage("Quotation not found")
		}
		return nil, err
	}
	return mapQuotation(item), nil
}
