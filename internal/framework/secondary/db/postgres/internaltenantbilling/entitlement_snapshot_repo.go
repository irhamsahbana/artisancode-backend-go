package repository

import (
	"context"
	"encoding/json"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *internalTenantBillingRepo) GetLatestEntitlementSnapshot(
	ctx context.Context,
	tenantID string,
) (*coreentity.InternalEntitlementSnapshot, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internaltenantbilling:get_latest_entitlement_snapshot:GetLatestEntitlementSnapshot",
	)
	defer span.End()

	var item coreentity.InternalEntitlementSnapshot
	query := `
		SELECT
			id,
			tenant_id::text,
			internal_tenant_subscription_id::text,
			subscription_status,
			features,
			usage_limits,
			effective_at::text,
			created_at::text,
			COALESCE(updated_at, created_at)::text AS updated_at
		FROM internal_entitlement_snapshots
		WHERE tenant_id = ?
			AND deleted_at IS NULL
		ORDER BY effective_at DESC
		LIMIT 1
	`
	if err := r.exec(ctx).GetContext(
		ctx,
		&item,
		r.exec(ctx).Rebind(query),
		tenantID,
	); err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *internalTenantBillingRepo) CreateEntitlementSnapshot(
	ctx context.Context,
	data coreentity.InternalEntitlementSnapshot,
) error {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internaltenantbilling:create_entitlement_snapshot:CreateEntitlementSnapshot",
	)
	defer span.End()

	features, _ := json.Marshal(data.Features)
	usageLimits, _ := json.Marshal(data.UsageLimits)
	sourceSnapshot, _ := json.Marshal(data.SourceSnapshot)
	metadata, _ := json.Marshal(data.Metadata)

	query := `
		INSERT INTO internal_entitlement_snapshots (
			tenant_id,
			internal_tenant_subscription_id,
			subscription_status,
			features,
			usage_limits,
			source_snapshot,
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
			?
		)
	`
	if _, err := r.exec(ctx).ExecContext(
		ctx,
		r.exec(ctx).Rebind(query),
		data.TenantID,
		data.InternalTenantSubscriptionID,
		data.SubscriptionStatus,
		features,
		usageLimits,
		sourceSnapshot,
		data.EffectiveAt,
		metadata,
	); err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to create tenant billing entitlement snapshot")
		return err
	}

	return nil
}
