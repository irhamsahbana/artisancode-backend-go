package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/rs/zerolog/log"
)

func (r *internalInvoiceRepo) GetPaymentAttempt(
	ctx context.Context,
	tenantID, id string,
) (*coreentity.InternalPaymentAttempt, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalinvoice:get_payment_attempt:GetPaymentAttempt",
	)
	defer span.End()

	items, err := r.getPaymentAttemptsByWhere(ctx, tenantID, "pa.id = ?", id, 1)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		log.Ctx(ctx).Warn().Any(common.LogKeyPayload, map[string]string{
			"tenant_id": tenantID,
			"id":        id,
		}).Msg("Payment attempt not found")
		return nil, errmsg.NewCustomErrors(404).SetMessage("Payment attempt not found")
	}
	return &items[0], nil
}
