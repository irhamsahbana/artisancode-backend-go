package repository

import (
	"context"
	"encoding/json"
	"time"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *internalTenantBillingRepo) GetSubscription(
	ctx context.Context,
	tenantID string,
) (*coreentity.InternalTenantSubscription, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internaltenantbilling:get_subscription:GetSubscription",
	)
	defer span.End()

	var item coreentity.InternalTenantSubscription
	query := `
		SELECT
			s.id,
			s.internal_billing_account_id,
			s.tenant_id::text,
			s.product_family,
			s.status,
			s.internal_product_id::text,
			s.internal_product_pricing_id::text,
			s.plan_snapshot,
			s.add_on_snapshots,
			s.current_period_started_at::text,
			s.current_period_ended_at::text,
			s.grace_ended_at::text,
			s.cancelled_at::text,
			s.expired_at::text,
			s.created_at::text,
			COALESCE(s.updated_at, s.created_at)::text AS updated_at
		FROM internal_tenant_subscriptions s
		WHERE s.tenant_id = ?
			AND s.deleted_at IS NULL
			AND s.status != ?
		ORDER BY s.created_at DESC
		LIMIT 1
	`
	if err := r.exec(ctx).GetContext(
		ctx,
		&item,
		r.exec(ctx).Rebind(query),
		tenantID,
		coreentity.InternalTenantSubscriptionStatusExpired,
	); err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *internalTenantBillingRepo) GetActiveSubscription(
	ctx context.Context,
	tenantID string,
) (*coreentity.InternalTenantSubscription, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internaltenantbilling:get_active_subscription:GetActiveSubscription",
	)
	defer span.End()

	var item coreentity.InternalTenantSubscription
	query := `
		SELECT
			s.id,
			s.tenant_id::text,
			s.product_family,
			s.status,
			s.internal_product_id::text,
			s.internal_product_pricing_id::text,
			s.created_at::text,
			COALESCE(s.updated_at, s.created_at)::text AS updated_at
		FROM internal_tenant_subscriptions s
		WHERE s.tenant_id = ?
			AND s.deleted_at IS NULL
			AND s.status IN (?, ?, ?, ?)
		ORDER BY s.created_at DESC
		LIMIT 1
	`
	if err := r.exec(ctx).GetContext(
		ctx,
		&item,
		r.exec(ctx).Rebind(query),
		tenantID,
		coreentity.InternalTenantSubscriptionStatusActive,
		coreentity.InternalTenantSubscriptionStatusGracePeriod,
		coreentity.InternalTenantSubscriptionStatusCancelled,
		coreentity.InternalTenantSubscriptionStatusSuspended,
	); err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *internalTenantBillingRepo) UpsertSubscription(
	ctx context.Context,
	data coreentity.InternalTenantSubscription,
) (*coreentity.InternalTenantSubscription, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internaltenantbilling:upsert_subscription:UpsertSubscription",
	)
	defer span.End()

	if data.ID == "" {
		return r.createSubscription(ctx, data)
	}

	metadata, _ := json.Marshal(data.Metadata)
	planSnapshot, _ := json.Marshal(data.PlanSnapshot)
	addOnSnapshots, _ := json.Marshal(data.AddOnSnapshots)

	now := time.Now().UTC().Format(time.RFC3339)
	query := `
		UPDATE internal_tenant_subscriptions
		SET
			status = ?,
			internal_product_id = ?,
			internal_product_pricing_id = ?,
			plan_snapshot = ?,
			add_on_snapshots = ?,
			current_period_started_at = ?,
			current_period_ended_at = ?,
			grace_ended_at = ?,
			cancelled_at = ?,
			expired_at = ?,
			metadata = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
			AND deleted_at IS NULL
		RETURNING
			id,
			internal_billing_account_id,
			tenant_id::text,
			product_family,
			status,
			internal_product_id::text,
			internal_product_pricing_id::text,
			current_period_started_at::text,
			current_period_ended_at::text,
			created_at::text,
			?::text AS updated_at
	`
	if err := r.exec(ctx).GetContext(
		ctx,
		&data,
		r.exec(ctx).Rebind(query),
		data.Status,
		data.InternalProductID,
		data.InternalProductPricingID,
		planSnapshot,
		addOnSnapshots,
		data.CurrentPeriodStartedAt,
		data.CurrentPeriodEndedAt,
		data.GraceEndedAt,
		data.CancelledAt,
		data.ExpiredAt,
		metadata,
		data.ID,
		now,
	); err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to update tenant billing subscription")
		return nil, err
	}

	return &data, nil
}

func (r *internalTenantBillingRepo) createSubscription(
	ctx context.Context,
	data coreentity.InternalTenantSubscription,
) (*coreentity.InternalTenantSubscription, error) {
	metadata, _ := json.Marshal(data.Metadata)
	planSnapshot, _ := json.Marshal(data.PlanSnapshot)
	addOnSnapshots, _ := json.Marshal(data.AddOnSnapshots)

	query := `
		INSERT INTO internal_tenant_subscriptions (
			internal_billing_account_id,
			tenant_id,
			product_family,
			status,
			internal_product_id,
			internal_product_pricing_id,
			plan_snapshot,
			add_on_snapshots,
			current_period_started_at,
			current_period_ended_at,
			grace_ended_at,
			cancelled_at,
			expired_at,
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
			?,
			?,
			?,
			?,
			?
		)
		RETURNING
			id,
			internal_billing_account_id,
			tenant_id::text,
			product_family,
			status,
			internal_product_id::text,
			internal_product_pricing_id::text,
			current_period_started_at::text,
			current_period_ended_at::text,
			created_at::text,
			COALESCE(updated_at, created_at)::text AS updated_at
	`
	if err := r.exec(ctx).QueryRowxContext(
		ctx,
		r.exec(ctx).Rebind(query),
		data.InternalBillingAccountID,
		data.TenantID,
		data.ProductFamily,
		data.Status,
		data.InternalProductID,
		data.InternalProductPricingID,
		planSnapshot,
		addOnSnapshots,
		data.CurrentPeriodStartedAt,
		data.CurrentPeriodEndedAt,
		data.GraceEndedAt,
		data.CancelledAt,
		data.ExpiredAt,
		metadata,
	).Scan(
		&data.ID,
		&data.InternalBillingAccountID,
		&data.TenantID,
		&data.ProductFamily,
		&data.Status,
		&data.InternalProductID,
		&data.InternalProductPricingID,
		&data.CurrentPeriodStartedAt,
		&data.CurrentPeriodEndedAt,
		&data.CreatedAt,
		&data.UpdatedAt,
	); err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to create tenant billing subscription")
		return nil, err
	}

	return &data, nil
}
