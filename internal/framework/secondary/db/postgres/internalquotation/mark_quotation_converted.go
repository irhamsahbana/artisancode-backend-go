package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/rs/zerolog/log"
)

func (r *internalQuotationRepo) MarkQuotationConverted(
	ctx context.Context,
	tenantID string,
	quotationID string,
	orderID string,
) (*coreentity.InternalQuotation, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalquotation:mark_quotation_converted:MarkQuotationConverted",
	)
	defer span.End()

	payload := map[string]any{
		"tenant_id":    tenantID,
		"quotation_id": quotationID,
		"order_id":     orderID,
	}

	query := `
		UPDATE internal_quotations
		SET
			status = 'converted',
			converted_to_order_id = ?,
			updated_at = NOW()
		WHERE id = ?
			AND tenant_id = ?
			AND deleted_at IS NULL
	`
	result, err := r.exec(ctx).ExecContext(ctx, r.exec(ctx).Rebind(query), orderID, quotationID, tenantID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, payload).Msg("Failed to mark quotation converted")
		return nil, err
	}
	if affected, _ := result.RowsAffected(); affected == 0 {
		log.Ctx(ctx).Warn().Any(common.LogKeyPayload, payload).Msg("Quotation not found when marking converted")
		return nil, errmsg.NewCustomErrors(404).SetMessage(errmsg.MessageQuotationNotFound)
	}
	return r.GetQuotation(ctx, tenantID, quotationID)
}
