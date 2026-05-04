package repository

import (
	"context"
	"strings"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *internalCurrencyRepo) UpsertProviderCurrency(
	ctx context.Context,
	data coreentity.InternalPaymentProviderCurrency,
) error {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalcurrency:upsert_provider_currency:UpsertProviderCurrency",
	)
	defer span.End()

	metadata, err := marshalMetadata(data.Metadata)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO internal_payment_provider_currencies (
			provider,
			currency_code,
			is_active,
			min_amount,
			max_amount,
			metadata,
			deleted_at
		)
		VALUES (?, ?, ?, ?, ?, ?, NULL)
		ON CONFLICT (provider, currency_code) DO UPDATE SET
			is_active = EXCLUDED.is_active,
			min_amount = EXCLUDED.min_amount,
			max_amount = EXCLUDED.max_amount,
			metadata = EXCLUDED.metadata,
			updated_at = NOW(),
			deleted_at = NULL
	`
	if _, err := r.db.ExecContext(
		ctx,
		r.db.Rebind(query),
		strings.ToLower(strings.TrimSpace(data.Provider)),
		strings.ToUpper(strings.TrimSpace(data.CurrencyCode)),
		data.IsActive,
		data.MinAmount,
		data.MaxAmount,
		metadata,
	); err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to upsert provider currency")
		return err
	}

	return nil
}
