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

func (r *internalCurrencyRepo) UpdateInternalCurrency(
	ctx context.Context,
	data coreentity.InternalCurrency,
) error {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalcurrency:update_currency:UpdateInternalCurrency",
	)
	defer span.End()

	metadata, err := marshalMetadata(data.Metadata)
	if err != nil {
		return err
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to begin internal currency update transaction")
		return err
	}
	defer tx.Rollback()

	code := strings.ToUpper(strings.TrimSpace(data.Code))
	current, err := getInternalCurrencyFromTx(ctx, tx, coreentity.InternalCurrencyFilter{Code: code})
	if err != nil {
		if err == sql.ErrNoRows {
			log.Ctx(ctx).Warn().Any(common.LogKeyPayload, data).Msg("Internal currency not found when updating")
			return errmsg.NewCustomErrors(404).SetMessage(errmsg.MessageCurrencyNotFound)
		}
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to query internal currency before update")
		return err
	}

	if current.IsDefault && !data.IsActive {
		log.Ctx(ctx).Warn().Any(common.LogKeyPayload, data).Msg("Default currency cannot be inactive")
		return errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageDefaultCurrencyCannotBeInactive)
	}

	if data.IsDefault {
		if _, err := tx.ExecContext(
			ctx,
			tx.Rebind(`UPDATE internal_currencies SET is_default = FALSE, updated_at = NOW() WHERE deleted_at IS NULL AND code <> ? AND is_default = TRUE`),
			code,
		); err != nil {
			log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to reset default internal currencies")
			return err
		}
	}

	updateQuery := `
		UPDATE internal_currencies
		SET
			name = ?,
			symbol = ?,
			decimal_places = ?,
			is_active = ?,
			is_default = ?,
			sort_order = ?,
			metadata = ?,
			updated_at = NOW()
		WHERE code = ?
			AND deleted_at IS NULL
	`
	result, err := tx.ExecContext(
		ctx,
		tx.Rebind(updateQuery),
		data.Name,
		data.Symbol,
		data.DecimalPlaces,
		data.IsActive,
		data.IsDefault,
		data.SortOrder,
		metadata,
		code,
	)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to update internal currency")
		return err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		log.Ctx(ctx).Warn().Any(common.LogKeyPayload, data).Msg("Internal currency not found when updating")
		return errmsg.NewCustomErrors(404).SetMessage(errmsg.MessageCurrencyNotFound)
	}

	if err := tx.Commit(); err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to commit internal currency update transaction")
		return err
	}

	return nil
}
