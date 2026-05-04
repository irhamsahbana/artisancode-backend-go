package repository

import (
	"context"
	"strings"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/rs/zerolog/log"
)

func (r *internalCurrencyRepo) DeleteProviderCurrency(
	ctx context.Context,
	filter coreentity.InternalPaymentProviderCurrencyFilter,
) error {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalcurrency:delete_provider_currency:DeleteProviderCurrency",
	)
	defer span.End()

	result, err := r.db.ExecContext(
		ctx,
		r.db.Rebind(`
			UPDATE internal_payment_provider_currencies
			SET deleted_at = NOW(), updated_at = NOW()
			WHERE provider = ?
				AND currency_code = ?
				AND deleted_at IS NULL
		`),
		strings.ToLower(strings.TrimSpace(filter.Provider)),
		strings.ToUpper(strings.TrimSpace(filter.CurrencyCode)),
	)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to delete provider currency")
		return err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		log.Ctx(ctx).Warn().Any(common.LogKeyPayload, filter).Msg("Provider currency not found when deleting")
		return errmsg.NewCustomErrors(404).SetMessage(errmsg.MessageCurrencyNotFound)
	}

	return nil
}

func (r *internalCurrencyRepo) IsProviderCurrencyActive(
	ctx context.Context,
	provider string,
	currencyCode string,
	amount string,
) (bool, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalcurrency:is_provider_currency_active:IsProviderCurrencyActive",
	)
	defer span.End()

	var total int
	query := `
		SELECT COUNT(*)
		FROM internal_payment_provider_currencies ipc
		JOIN internal_currencies ic
			ON ic.code = ipc.currency_code
			AND ic.deleted_at IS NULL
		WHERE ipc.provider = ?
			AND ipc.currency_code = ?
			AND ipc.deleted_at IS NULL
			AND ipc.is_active = TRUE
			AND ic.is_active = TRUE
			AND (ipc.min_amount IS NULL OR ipc.min_amount <= ?::numeric)
			AND (ipc.max_amount IS NULL OR ipc.max_amount >= ?::numeric)
	`
	if err := r.db.GetContext(
		ctx,
		&total,
		r.db.Rebind(query),
		strings.ToLower(strings.TrimSpace(provider)),
		strings.ToUpper(strings.TrimSpace(currencyCode)),
		amount,
		amount,
	); err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, map[string]any{
			"provider":      provider,
			"currency_code": currencyCode,
			"amount":        amount,
		}).Msg("Failed to query provider currency active flag")
		return false, err
	}

	return total > 0, nil
}
