package repository

import (
	"context"
	"encoding/json"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *internalTenantBillingRepo) UpdatePaymentAttemptGateway(
	ctx context.Context,
	data coreentity.InternalTenantPaymentAttempt,
) (*coreentity.InternalTenantPaymentAttempt, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internaltenantbilling:update_payment_attempt_gateway:UpdatePaymentAttemptGateway",
	)
	defer span.End()

	payload, _ := json.Marshal(data.ProviderPayloadSnapshot)
	query := `
		UPDATE internal_tenant_payment_attempts
		SET
			provider_request_id = ?,
			provider_payment_url = ?,
			provider_payload_snapshot = ?,
			status = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
			AND tenant_id = ?
			AND deleted_at IS NULL
		RETURNING
			updated_at::text AS updated_at
	`
	if err := r.exec(ctx).QueryRowxContext(
		ctx,
		r.exec(ctx).Rebind(query),
		data.ProviderRequestID,
		data.ProviderPaymentURL,
		payload,
		data.Status,
		data.ID,
		data.TenantID,
	).Scan(&data.UpdatedAt); err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to update tenant billing payment attempt gateway")
		return nil, err
	}

	return &data, nil
}
