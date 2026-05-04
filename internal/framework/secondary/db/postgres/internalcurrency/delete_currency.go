package repository

import (
	"context"
	"database/sql"
	"strings"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/rs/zerolog/log"
)

func (r *internalCurrencyRepo) DeleteInternalCurrency(
	ctx context.Context,
	filter coreentity.InternalCurrencyDeleteFilter,
) error {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalcurrency:delete_currency:DeleteInternalCurrency",
	)
	defer span.End()

	code := strings.ToUpper(strings.TrimSpace(filter.Code))
	current, err := r.GetInternalCurrency(ctx, coreentity.InternalCurrencyFilter{Code: code})
	if err != nil {
		return err
	}
	if current.IsDefault {
		log.Ctx(ctx).Warn().Any(common.LogKeyPayload, filter).Msg("Default currency cannot be deleted")
		return errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageDefaultCurrencyCannotBeDeleted)
	}

	usageQuery := `
		SELECT
			COALESCE((
				SELECT COUNT(*)
				FROM internal_product_prices
				WHERE currency_code = ?
					AND deleted_at IS NULL
			), 0) +
			COALESCE((
				SELECT COUNT(*)
				FROM internal_tenant_invoices
				WHERE currency_code = ?
					AND deleted_at IS NULL
			), 0) +
			COALESCE((
				SELECT COUNT(*)
				FROM internal_payment_provider_currencies
				WHERE currency_code = ?
					AND deleted_at IS NULL
			), 0) AS total_usage
	`
	var totalUsage int
	if err := r.db.GetContext(ctx, &totalUsage, r.db.Rebind(usageQuery), code, code, code); err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to check internal currency usage")
		return err
	}
	if totalUsage > 0 {
		log.Ctx(ctx).Warn().Any(common.LogKeyPayload, filter).Msg("Internal currency is still used")
		return errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageCurrencyIsStillUsed)
	}

	result, err := r.db.ExecContext(
		ctx,
		r.db.Rebind(`UPDATE internal_currencies SET deleted_at = NOW(), updated_at = NOW() WHERE code = ? AND deleted_at IS NULL`),
		code,
	)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to delete internal currency")
		return err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		log.Ctx(ctx).Warn().Any(common.LogKeyPayload, filter).Msg("Internal currency not found when deleting")
		return errmsg.NewCustomErrors(404).SetMessage(errmsg.MessageCurrencyNotFound)
	}

	return nil
}

func (r *internalCurrencyRepo) IsCurrencyActive(ctx context.Context, code string) (bool, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalcurrency:is_currency_active:IsCurrencyActive",
	)
	defer span.End()

	var exists sql.NullBool
	query := `
		SELECT is_active
		FROM internal_currencies
		WHERE code = ?
			AND deleted_at IS NULL
	`
	if err := r.db.GetContext(ctx, &exists, r.db.Rebind(query), strings.ToUpper(strings.TrimSpace(code))); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, code).Msg("Failed to query internal currency active flag")
		return false, err
	}

	return exists.Valid && exists.Bool, nil
}

func (r *internalCurrencyRepo) GetDefaultCurrency(ctx context.Context) (*coreentity.InternalCurrency, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalcurrency:get_default_currency:GetDefaultCurrency",
	)
	defer span.End()

	var code string
	query := `
		SELECT code
		FROM internal_currencies
		WHERE is_default = TRUE
			AND is_active = TRUE
			AND deleted_at IS NULL
		LIMIT 1
	`
	if err := r.db.GetContext(ctx, &code, r.db.Rebind(query)); err != nil {
		if err == sql.ErrNoRows {
			log.Ctx(ctx).Warn().Msg("Default internal currency not found")
			return nil, errmsg.NewCustomErrors(404).SetMessage(errmsg.MessageCurrencyNotFound)
		}
		log.Ctx(ctx).Error().Err(err).Msg("Failed to query default internal currency")
		return nil, err
	}

	return r.GetInternalCurrency(ctx, coreentity.InternalCurrencyFilter{Code: code})
}
