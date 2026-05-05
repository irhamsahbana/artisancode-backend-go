package repository

import (
	"context"
	"encoding/json"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *internalTenantBillingRepo) CreateLedgerEntry(
	ctx context.Context,
	data coreentity.InternalBillingLedgerEntry,
) error {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internaltenantbilling:create_ledger_entry:CreateLedgerEntry",
	)
	defer span.End()

	metadata, _ := json.Marshal(data.Metadata)

	query := `
		INSERT INTO internal_billing_ledger_entries (
			tenant_id,
			internal_tenant_subscription_id,
			entry_type,
			source_type,
			source_reference_id,
			occurred_at,
			metadata
		)
		VALUES (
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
		data.EntryType,
		data.SourceType,
		data.SourceReferenceID,
		data.OccurredAt,
		metadata,
	); err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to create tenant billing ledger entry")
		return err
	}

	return nil
}
