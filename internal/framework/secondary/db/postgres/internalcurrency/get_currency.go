package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"

	"github.com/rs/zerolog/log"
)

func (r *internalCurrencyRepo) GetInternalCurrency(
	ctx context.Context,
	filter coreentity.InternalCurrencyFilter,
) (*coreentity.InternalCurrency, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalcurrency:get_currency:GetInternalCurrency",
	)
	defer span.End()

	type dao struct {
		Code          string          `db:"code"`
		Name          string          `db:"name"`
		Symbol        string          `db:"symbol"`
		DecimalPlaces int             `db:"decimal_places"`
		IsActive      bool            `db:"is_active"`
		IsDefault     bool            `db:"is_default"`
		SortOrder     int             `db:"sort_order"`
		Metadata      json.RawMessage `db:"metadata"`
		CreatedAt     string          `db:"created_at"`
		UpdatedAt     *string         `db:"updated_at"`
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
			metadata,
			created_at::text AS created_at,
			updated_at::text AS updated_at
		FROM internal_currencies
		WHERE code = ?
			AND deleted_at IS NULL
	`
	if err := r.db.GetContext(ctx, row, r.db.Rebind(query), strings.ToUpper(strings.TrimSpace(filter.Code))); err != nil {
		if err == sql.ErrNoRows {
			log.Ctx(ctx).Warn().Any(common.LogKeyPayload, filter).Msg("Internal currency not found")
			return nil, errmsg.NewCustomErrors(404).SetMessage(errmsg.MessageCurrencyNotFound)
		}
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to query internal currency")
		return nil, err
	}

	item, err := mapInternalCurrencyDAO(
		row.Code,
		row.Name,
		row.Symbol,
		row.DecimalPlaces,
		row.IsActive,
		row.IsDefault,
		row.SortOrder,
		row.Metadata,
		row.CreatedAt,
		row.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &item, nil
}
