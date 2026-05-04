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

func (r *internalCurrencyRepo) GetInternalCurrencies(
	ctx context.Context,
	filter coreentity.InternalCurrencyListFilter,
) ([]coreentity.InternalCurrency, int, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalcurrency:get_currencies:GetInternalCurrencies",
	)
	defer span.End()

	type dao struct {
		TotalData     int             `db:"total_data"`
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

	args := []any{}
	query := `
		SELECT
			COUNT(*) OVER() AS total_data,
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
		WHERE deleted_at IS NULL
	`

	if trimmedQ := strings.TrimSpace(filter.Q); trimmedQ != "" {
		query += ` AND (code ILIKE ? OR name ILIKE ?)`
		args = append(args, "%"+trimmedQ+"%", "%"+trimmedQ+"%")
	}
	if filter.IsActive != nil {
		query += ` AND is_active = ?`
		args = append(args, *filter.IsActive)
	}

	query += `
		ORDER BY
			is_default DESC,
			sort_order ASC,
			code ASC
		LIMIT ?
		OFFSET ?
	`
	args = append(args, filter.Paginate, (filter.Page-1)*filter.Paginate)

	rows := make([]dao, 0)
	if err := r.db.SelectContext(ctx, &rows, r.db.Rebind(query), args...); err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to query internal currencies")
		return nil, 0, err
	}

	items := make([]coreentity.InternalCurrency, 0, len(rows))
	total := 0
	for _, row := range rows {
		total = row.TotalData
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
			return nil, 0, err
		}
		items = append(items, item)
	}

	return items, total, nil
}
