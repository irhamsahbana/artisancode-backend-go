package repository

import (
	"context"
	"encoding/json"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
)

func (r *internalInvoiceRepo) UpdatePaymentAttemptGateway(
	ctx context.Context,
	data coreentity.InternalPaymentAttempt,
) (*coreentity.InternalPaymentAttempt, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalinvoice:update_payment_attempt_gateway:UpdatePaymentAttemptGateway",
	)
	defer span.End()

	payload, _ := json.Marshal(data.ProviderPayloadSnapshot)
	query := `
		UPDATE internal_payment_attempts
		SET
			provider_request_id = ?,
			provider_payment_url = ?,
			provider_payload_snapshot = ?,
			status = ?,
			updated_at = NOW()
		WHERE id = ?
			AND deleted_at IS NULL
	`
	if _, err := r.exec(ctx).ExecContext(ctx, r.exec(ctx).Rebind(query),
		data.ProviderRequestID,
		data.ProviderPaymentURL,
		payload,
		data.Status,
		data.ID,
	); err != nil {
		return nil, err
	}
	return r.GetPaymentAttempt(ctx, "", data.ID)
}
