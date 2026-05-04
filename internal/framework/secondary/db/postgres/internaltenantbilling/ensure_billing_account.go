package repository

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *internalTenantBillingRepo) EnsureBillingAccount(
	ctx context.Context,
	tenantID string,
) (*coreentity.InternalBillingAccount, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internaltenantbilling:ensure_billing_account:EnsureBillingAccount",
	)
	defer span.End()

	var item coreentity.InternalBillingAccount
	query := `
		INSERT INTO internal_billing_accounts (
			tenant_id,
			status,
			metadata
		)
		VALUES (
			?,
			?,
			'{}'::jsonb
		)
		ON CONFLICT (tenant_id)
			WHERE deleted_at IS NULL
		DO UPDATE SET
			updated_at = CURRENT_TIMESTAMP
		RETURNING
			id,
			tenant_id,
			status,
			created_at::text AS created_at,
			COALESCE(updated_at, created_at)::text AS updated_at
	`
	if err := r.exec(ctx).GetContext(
		ctx,
		&item,
		r.exec(ctx).Rebind(query),
		tenantID,
		coreentity.InternalTenantBillingAccountStatusOpen,
	); err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, tenantID).Msg("Failed to ensure tenant billing account")
		return nil, err
	}

	return &item, nil
}
