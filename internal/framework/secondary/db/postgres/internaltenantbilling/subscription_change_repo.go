package repository

import (
	"context"
	"encoding/json"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *internalTenantBillingRepo) CreateSubscriptionChange(
	ctx context.Context,
	data coreentity.InternalTenantSubscriptionChange,
) error {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internaltenantbilling:create_subscription_change:CreateSubscriptionChange",
	)
	defer span.End()

	metadata, _ := json.Marshal(data.Metadata)

	query := `
		INSERT INTO internal_tenant_subscription_changes (
			tenant_id,
			internal_tenant_subscription_id,
			change_type,
			from_status,
			to_status,
			trigger,
			source_type,
			source_reference_id,
			effective_at,
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
			?
		)
	`
	if _, err := r.exec(ctx).ExecContext(
		ctx,
		r.exec(ctx).Rebind(query),
		data.TenantID,
		data.InternalTenantSubscriptionID,
		data.ChangeType,
		data.FromStatus,
		data.ToStatus,
		data.Trigger,
		data.SourceType,
		data.SourceReferenceID,
		data.EffectiveAt,
		metadata,
	); err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to create tenant billing subscription change")
		return err
	}

	return nil
}
