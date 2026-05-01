package repository

import (
	"context"
	"database/sql"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/rs/zerolog/log"
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

	payload := map[string]any{
		"tenant_id": tenantID,
		"id":        id,
	}

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
			log.Ctx(ctx).Warn().Any(common.LogKeyPayload, payload).Msg("Quotation not found when retrieving")
			return nil, errmsg.NewCustomErrors(404).SetMessage(errmsg.MessageQuotationNotFound)
		}
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, payload).Msg("Failed to retrieve quotation")
		return nil, err
	}
	return mapQuotation(item), nil
}
