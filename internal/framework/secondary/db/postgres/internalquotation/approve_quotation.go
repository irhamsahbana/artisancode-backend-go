package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/rs/zerolog/log"
)

func (r *internalQuotationRepo) ApproveQuotation(
	ctx context.Context,
	tenantID, id string,
) (*coreentity.InternalQuotation, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalquotation:approve_quotation:ApproveQuotation",
	)
	defer span.End()

	payload := map[string]any{
		"tenant_id": tenantID,
		"id":        id,
	}

	query := `
		UPDATE internal_quotations
		SET
			status = 'approved',
			approved_at = NOW(),
			updated_at = NOW()
			WHERE id = ?
			AND tenant_id = ?
			AND deleted_at IS NULL
	`
	result, err := r.exec(ctx).ExecContext(ctx, r.exec(ctx).Rebind(query), id, tenantID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, payload).Msg("Failed to approve quotation")
		return nil, err
	}
	if affected, _ := result.RowsAffected(); affected == 0 {
		log.Ctx(ctx).
			Warn().
			Any(common.LogKeyPayload, payload).
			Msg("Quotation not found or already deleted when approving")
		return nil, errmsg.NewCustomErrors(404).SetMessage("Quotation not found")
	}
	return r.GetQuotation(ctx, tenantID, id)
}
