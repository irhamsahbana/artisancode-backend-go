package repository

import (
	"context"
	"strings"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (r *internalCurrencyRepo) CreateInternalCurrency(
	ctx context.Context,
	data coreentity.InternalCurrency,
) (*coreentity.InternalCurrency, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalcurrency:create_currency:CreateInternalCurrency",
	)
	defer span.End()

	metadata, err := marshalMetadata(data.Metadata)
	if err != nil {
		return nil, err
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to begin internal currency transaction")
		return nil, err
	}
	defer tx.Rollback()

	code := strings.ToUpper(strings.TrimSpace(data.Code))
	existing := 0
	checkQuery := `
		SELECT COUNT(*)
		FROM internal_currencies
		WHERE code = ?
			AND deleted_at IS NULL
	`
	if err := tx.GetContext(ctx, &existing, tx.Rebind(checkQuery), code); err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to check existing internal currency")
		return nil, err
	}
	if existing > 0 {
		log.Ctx(ctx).Warn().Any(common.LogKeyPayload, data).Msg("Internal currency code already exists")
		return nil, errmsg.NewCustomErrors(400).SetMessage(errmsg.MessageCurrencyCodeAlreadyExists)
	}

	if data.IsDefault {
		if _, err := tx.ExecContext(
			ctx,
			tx.Rebind(`UPDATE internal_currencies SET is_default = FALSE, updated_at = NOW() WHERE deleted_at IS NULL AND is_default = TRUE`),
		); err != nil {
			log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to reset default internal currency")
			return nil, err
		}
	}

	insertQuery := `
		INSERT INTO internal_currencies (
			code,
			name,
			symbol,
			decimal_places,
			is_active,
			is_default,
			sort_order,
			metadata
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	if _, err := tx.ExecContext(
		ctx,
		tx.Rebind(insertQuery),
		code,
		data.Name,
		data.Symbol,
		data.DecimalPlaces,
		data.IsActive,
		data.IsDefault,
		data.SortOrder,
		metadata,
	); err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to create internal currency")
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to commit internal currency transaction")
		return nil, err
	}

	return r.GetInternalCurrency(ctx, coreentity.InternalCurrencyFilter{Code: code})
}

func getInternalCurrencyFromTx(
	ctx context.Context,
	tx *sqlx.Tx,
	filter coreentity.InternalCurrencyFilter,
) (*coreentity.InternalCurrency, error) {
	type dao struct {
		Code          string  `db:"code"`
		Name          string  `db:"name"`
		Symbol        string  `db:"symbol"`
		DecimalPlaces int     `db:"decimal_places"`
		IsActive      bool    `db:"is_active"`
		IsDefault     bool    `db:"is_default"`
		SortOrder     int     `db:"sort_order"`
		CreatedAt     string  `db:"created_at"`
		UpdatedAt     *string `db:"updated_at"`
	}

	row := new(dao)
	query := `
		SELECT
			code,
			name,
			symbol,
			decimal_places,
			is_active,
			is_default,
			sort_order,
			created_at::text AS created_at,
			updated_at::text AS updated_at
		FROM internal_currencies
		WHERE code = ?
			AND deleted_at IS NULL
	`
	if err := tx.GetContext(ctx, row, tx.Rebind(query), strings.ToUpper(strings.TrimSpace(filter.Code))); err != nil {
		return nil, err
	}

	item := &coreentity.InternalCurrency{
		Code:          row.Code,
		Name:          row.Name,
		Symbol:        row.Symbol,
		DecimalPlaces: row.DecimalPlaces,
		IsActive:      row.IsActive,
		IsDefault:     row.IsDefault,
		SortOrder:     row.SortOrder,
		Metadata:      map[string]any{},
		CreatedAt:     row.CreatedAt,
	}
	if row.UpdatedAt != nil {
		item.UpdatedAt = *row.UpdatedAt
	}

	return item, nil
}
