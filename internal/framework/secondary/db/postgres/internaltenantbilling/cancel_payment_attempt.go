package repository

import (
	"context"

	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *internalTenantBillingRepo) CancelPaymentAttempt(
	ctx context.Context,
	tenantID string,
	id string,
) error {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internaltenantbilling:cancel_payment_attempt:CancelPaymentAttempt",
	)
	defer span.End()

	query := `
		UPDATE internal_tenant_payment_attempts
		SET status = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
			AND tenant_id = ?
			AND deleted_at IS NULL
			AND status IN (?, ?)
	`
	result, err := r.exec(ctx).ExecContext(
		ctx,
		r.exec(ctx).Rebind(query),
		coreentity.InternalTenantPaymentAttemptStatusCancelled,
		id,
		tenantID,
		coreentity.InternalTenantPaymentAttemptStatusInitiated,
		coreentity.InternalTenantPaymentAttemptStatusPending,
	)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Str("id", id).Msg("Failed to cancel tenant billing payment attempt")
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return nil
	}

	return nil
}
