package repository

import (
	"context"
	"encoding/json"
	"strings"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"

	"github.com/rs/zerolog/log"
)

func (r *internalCurrencyRepo) GetProviderCurrencies(
	ctx context.Context,
	filter coreentity.InternalPaymentProviderCurrencyListFilter,
) ([]coreentity.InternalPaymentProviderCurrency, int, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalcurrency:get_provider_currencies:GetProviderCurrencies",
	)
	defer span.End()

	type dao struct {
		TotalData    int             `db:"total_data"`
		Provider     string          `db:"provider"`
		CurrencyCode string          `db:"currency_code"`
		IsActive     bool            `db:"is_active"`
		MinAmount    *string         `db:"min_amount"`
		MaxAmount    *string         `db:"max_amount"`
		Metadata     json.RawMessage `db:"metadata"`
		CreatedAt    string          `db:"created_at"`
		UpdatedAt    *string         `db:"updated_at"`
	}

	rows := make([]dao, 0)
	query := `
		SELECT
			COUNT(*) OVER() AS total_data,
			provider,
			currency_code,
			is_active,
			min_amount::text AS min_amount,
			max_amount::text AS max_amount,
			metadata,
			created_at::text AS created_at,
			updated_at::text AS updated_at
		FROM internal_payment_provider_currencies
		WHERE provider = ?
			AND deleted_at IS NULL
		ORDER BY currency_code ASC
		LIMIT ?
		OFFSET ?
	`
	args := []any{
		strings.ToLower(strings.TrimSpace(filter.Provider)),
		filter.Paginate,
		(filter.Page - 1) * filter.Paginate,
	}

	if err := r.db.SelectContext(ctx, &rows, r.db.Rebind(query), args...); err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to query provider currencies")
		return nil, 0, err
	}

	items := make([]coreentity.InternalPaymentProviderCurrency, 0, len(rows))
	total := 0
	for _, row := range rows {
		total = row.TotalData
		item, err := mapInternalPaymentProviderCurrencyDAO(
			row.Provider,
			row.CurrencyCode,
			row.IsActive,
			row.MinAmount,
			row.MaxAmount,
			row.Metadata,
			row.CreatedAt,
			row.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}

	return items, total, nil
}
