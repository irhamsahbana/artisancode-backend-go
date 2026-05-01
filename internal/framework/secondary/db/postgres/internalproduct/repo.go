package repository

import (
	"context"
	"database/sql"
	"encoding/json"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	portsRepo "codebase-app/internal/ports/secondary/db"
	"codebase-app/pkg/errmsg"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
)

type internalProductRepo struct {
	db *sqlx.DB
}

type Config struct {
	DB *sqlx.DB
}

var _ portsRepo.InternalProductRepository = &internalProductRepo{}

func NewInternalProductRepository(cfg Config) portsRepo.InternalProductRepository {
	return &internalProductRepo{db: cfg.DB}
}

func (r *internalProductRepo) GetInternalProducts(
	ctx context.Context,
	filter coreentity.InternalProductListFilter,
) ([]coreentity.InternalProduct, int, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalproduct:repo:GetInternalProducts",
	)
	defer span.End()

	type dao struct {
		TotalData   int             `db:"total_data"`
		ID          string          `db:"id"`
		Code        string          `db:"code"`
		Name        string          `db:"name"`
		Description string          `db:"description"`
		Status      string          `db:"status"`
		Metadata    json.RawMessage `db:"metadata"`
		CreatedAt   string          `db:"created_at"`
		UpdatedAt   *string         `db:"updated_at"`
	}

	var (
		rows  = make([]dao, 0)
		items = make([]coreentity.InternalProduct, 0)
		args  = make([]any, 0, 4)
		total int
	)

	query := `
		SELECT
			COUNT(*) OVER() AS total_data,
			id,
			code,
			name,
			description,
			status,
			metadata,
			created_at,
			updated_at
		FROM internal_products
		WHERE deleted_at IS NULL
	`
	if filter.Q != "" {
		query += ` AND (name ILIKE '%' || ? || '%' OR code ILIKE '%' || ? || '%')`
		args = append(args, filter.Q, filter.Q)
	}
	if filter.Status != "" {
		query += ` AND status = ?`
		args = append(args, filter.Status)
	}
	query += ` ORDER BY name ASC LIMIT ? OFFSET ?`
	args = append(args, filter.Paginate, (filter.Page-1)*filter.Paginate)

	if err := r.db.SelectContext(ctx, &rows, r.db.Rebind(query), args...); err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to query internal products")
		return nil, 0, err
	}

	for _, row := range rows {
		total = row.TotalData
		item, err := mapInternalProductDAO(
			row.ID,
			row.Code,
			row.Name,
			row.Description,
			row.Status,
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

func (r *internalProductRepo) GetInternalProduct(
	ctx context.Context,
	filter coreentity.InternalProductFilter,
) (*coreentity.InternalProduct, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalproduct:repo:GetInternalProduct",
	)
	defer span.End()

	type dao struct {
		ID          string          `db:"id"`
		Code        string          `db:"code"`
		Name        string          `db:"name"`
		Description string          `db:"description"`
		Status      string          `db:"status"`
		Metadata    json.RawMessage `db:"metadata"`
		CreatedAt   string          `db:"created_at"`
		UpdatedAt   *string         `db:"updated_at"`
	}

	var row dao
	query := `
		SELECT
			id,
			code,
			name,
			description,
			status,
			metadata,
			created_at,
			updated_at
		FROM internal_products
		WHERE id = ?
			AND deleted_at IS NULL
	`
	if err := r.db.GetContext(ctx, &row, r.db.Rebind(query), filter.ID); err != nil {
		if err == sql.ErrNoRows {
			log.Ctx(ctx).Warn().Any(common.LogKeyPayload, filter).Msg(errmsg.MessageInternalProductNotFound)
			return nil, errmsg.NewCustomErrors(404).SetMessage(errmsg.MessageInternalProductNotFound)
		}
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to get internal product")
		return nil, err
	}

	item, err := mapInternalProductDAO(
		row.ID,
		row.Code,
		row.Name,
		row.Description,
		row.Status,
		row.Metadata,
		row.CreatedAt,
		row.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *internalProductRepo) CreateInternalProduct(
	ctx context.Context,
	data coreentity.InternalProduct,
) (*coreentity.InternalProduct, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalproduct:repo:CreateInternalProduct",
	)
	defer span.End()

	metadata, err := marshalMetadata(data.Metadata)
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO internal_products (
			code,
			name,
			description,
			status,
			metadata
		)
		VALUES (
			?,
			?,
			?,
			?,
			?
		)
		RETURNING
			id
	`
	if err := r.db.GetContext(ctx, &data.ID, r.db.Rebind(query),
		data.Code,
		data.Name,
		data.Description,
		data.Status,
		metadata,
	); err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to create internal product")
		return nil, err
	}

	return &data, nil
}

func (r *internalProductRepo) UpdateInternalProduct(ctx context.Context, data coreentity.InternalProduct) error {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalproduct:repo:UpdateInternalProduct",
	)
	defer span.End()

	metadata, err := marshalMetadata(data.Metadata)
	if err != nil {
		return err
	}

	query := `
		UPDATE internal_products
		SET
			code = ?,
			name = ?,
			description = ?,
			status = ?,
			metadata = ?,
			updated_at = NOW()
		WHERE id = ?
			AND deleted_at IS NULL
	`
	result, err := r.db.ExecContext(ctx, r.db.Rebind(query),
		data.Code,
		data.Name,
		data.Description,
		data.Status,
		metadata,
		data.ID,
	)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to update internal product")
		return err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		log.Ctx(ctx).Warn().Any(common.LogKeyPayload, data).Msg("Internal product not found when updating")
		return errmsg.NewCustomErrors(404).SetMessage(errmsg.MessageInternalProductNotFound)
	}
	return nil
}

func (r *internalProductRepo) DeleteInternalProduct(
	ctx context.Context,
	filter coreentity.InternalProductDeleteFilter,
) error {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalproduct:repo:DeleteInternalProduct",
	)
	defer span.End()

	query := `
		UPDATE internal_products
		SET deleted_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`
	result, err := r.db.ExecContext(ctx, r.db.Rebind(query), filter.ID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to delete internal product")
		return err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		log.Ctx(ctx).Warn().Any(common.LogKeyPayload, filter).Msg("Internal product not found when deleting")
		return errmsg.NewCustomErrors(404).SetMessage(errmsg.MessageInternalProductNotFound)
	}
	return nil
}

func (r *internalProductRepo) ExistsInternalProductByCode(ctx context.Context, code, excludeID string) (bool, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalproduct:repo:ExistsInternalProductByCode",
	)
	defer span.End()

	query := `
		SELECT COUNT(*)
		FROM internal_products
		WHERE code = ? AND deleted_at IS NULL
	`
	args := []any{code}
	if excludeID != "" {
		query += ` AND id != ?`
		args = append(args, excludeID)
	}

	var count int
	if err := r.db.GetContext(ctx, &count, r.db.Rebind(query), args...); err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, map[string]any{
			"code":       code,
			"exclude_id": excludeID,
		}).Msg("Failed to check internal product code")
		return false, err
	}

	return count > 0, nil
}

func (r *internalProductRepo) GetInternalProductPricings(
	ctx context.Context,
	filter coreentity.InternalProductPricingListFilter,
) ([]coreentity.InternalProductPricing, int, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalproduct:repo:GetInternalProductPricings",
	)
	defer span.End()

	type dao struct {
		TotalData         int             `db:"total_data"`
		ID                string          `db:"id"`
		InternalProductID string          `db:"internal_product_id"`
		Code              string          `db:"code"`
		Name              string          `db:"name"`
		Description       string          `db:"description"`
		Status            string          `db:"status"`
		Metadata          json.RawMessage `db:"metadata"`
		CreatedAt         string          `db:"created_at"`
		UpdatedAt         *string         `db:"updated_at"`
	}

	var (
		rows  = make([]dao, 0)
		items = make([]coreentity.InternalProductPricing, 0)
		args  = []any{filter.InternalProductID}
		total int
	)

	query := `
		SELECT
			COUNT(*) OVER() AS total_data,
			id,
			internal_product_id,
			code,
			name,
			description,
			status,
			metadata,
			created_at,
			updated_at
		FROM internal_product_pricings
		WHERE internal_product_id = ?
			AND deleted_at IS NULL
	`
	if filter.Q != "" {
		query += ` AND (name ILIKE '%' || ? || '%' OR code ILIKE '%' || ? || '%')`
		args = append(args, filter.Q, filter.Q)
	}
	query += ` ORDER BY name ASC LIMIT ? OFFSET ?`
	args = append(args, filter.Paginate, (filter.Page-1)*filter.Paginate)

	if err := r.db.SelectContext(ctx, &rows, r.db.Rebind(query), args...); err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to query internal product pricings")
		return nil, 0, err
	}

	for _, row := range rows {
		total = row.TotalData
		item, err := mapInternalProductPricingDAO(
			row.ID,
			row.InternalProductID,
			row.Code,
			row.Name,
			row.Description,
			row.Status,
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

func (r *internalProductRepo) GetInternalProductPricing(
	ctx context.Context,
	filter coreentity.InternalProductPricingFilter,
) (*coreentity.InternalProductPricing, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalproduct:repo:GetInternalProductPricing",
	)
	defer span.End()

	type dao struct {
		ID                string          `db:"id"`
		InternalProductID string          `db:"internal_product_id"`
		Code              string          `db:"code"`
		Name              string          `db:"name"`
		Description       string          `db:"description"`
		Status            string          `db:"status"`
		Metadata          json.RawMessage `db:"metadata"`
		CreatedAt         string          `db:"created_at"`
		UpdatedAt         *string         `db:"updated_at"`
	}

	var row dao
	query := `
		SELECT
			id,
			internal_product_id,
			code,
			name,
			description,
			status,
			metadata,
			created_at,
			updated_at
		FROM internal_product_pricings
		WHERE id = ?
			AND deleted_at IS NULL
	`
	if err := r.db.GetContext(ctx, &row, r.db.Rebind(query), filter.ID); err != nil {
		if err == sql.ErrNoRows {
			log.Ctx(ctx).Warn().Any(common.LogKeyPayload, filter).Msg(errmsg.MessageInternalProductPricingNotFound)
			return nil, errmsg.NewCustomErrors(404).SetMessage(errmsg.MessageInternalProductPricingNotFound)
		}
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to get internal product pricing")
		return nil, err
	}

	item, err := mapInternalProductPricingDAO(
		row.ID,
		row.InternalProductID,
		row.Code,
		row.Name,
		row.Description,
		row.Status,
		row.Metadata,
		row.CreatedAt,
		row.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *internalProductRepo) CreateInternalProductPricing(
	ctx context.Context,
	data coreentity.InternalProductPricing,
) (*coreentity.InternalProductPricing, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalproduct:repo:CreateInternalProductPricing",
	)
	defer span.End()

	metadata, err := marshalMetadata(data.Metadata)
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO internal_product_pricings (
			internal_product_id,
			code,
			name,
			description,
			status,
			metadata
		)
		VALUES (
			?,
			?,
			?,
			?,
			?,
			?
		)
		RETURNING
			id
	`
	if err := r.db.GetContext(ctx, &data.ID, r.db.Rebind(query),
		data.InternalProductID,
		data.Code,
		data.Name,
		data.Description,
		data.Status,
		metadata,
	); err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to create internal product pricing")
		return nil, err
	}

	return &data, nil
}

func (r *internalProductRepo) UpdateInternalProductPricing(
	ctx context.Context,
	data coreentity.InternalProductPricing,
) error {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalproduct:repo:UpdateInternalProductPricing",
	)
	defer span.End()

	metadata, err := marshalMetadata(data.Metadata)
	if err != nil {
		return err
	}

	query := `
		UPDATE internal_product_pricings
		SET
			internal_product_id = ?,
			code = ?,
			name = ?,
			description = ?,
			status = ?,
			metadata = ?,
			updated_at = NOW()
		WHERE id = ?
			AND deleted_at IS NULL
	`
	result, err := r.db.ExecContext(ctx, r.db.Rebind(query),
		data.InternalProductID,
		data.Code,
		data.Name,
		data.Description,
		data.Status,
		metadata,
		data.ID,
	)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to update internal product pricing")
		return err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		log.Ctx(ctx).Warn().Any(common.LogKeyPayload, data).Msg("Internal product pricing not found when updating")
		return errmsg.NewCustomErrors(404).SetMessage(errmsg.MessageInternalProductPricingNotFound)
	}
	return nil
}

func (r *internalProductRepo) DeleteInternalProductPricing(
	ctx context.Context,
	filter coreentity.InternalProductPricingDeleteFilter,
) error {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalproduct:repo:DeleteInternalProductPricing",
	)
	defer span.End()

	query := `
		UPDATE internal_product_pricings
		SET deleted_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`
	result, err := r.db.ExecContext(ctx, r.db.Rebind(query), filter.ID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to delete internal product pricing")
		return err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		log.Ctx(ctx).Warn().Any(common.LogKeyPayload, filter).Msg("Internal product pricing not found when deleting")
		return errmsg.NewCustomErrors(404).SetMessage(errmsg.MessageInternalProductPricingNotFound)
	}
	return nil
}

func (r *internalProductRepo) ExistsInternalProductPricingByCode(
	ctx context.Context,
	internalProductID string,
	code string,
	excludeID string,
) (bool, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalproduct:repo:ExistsInternalProductPricingByCode",
	)
	defer span.End()

	query := `
		SELECT COUNT(*)
		FROM internal_product_pricings
		WHERE internal_product_id = ? AND code = ? AND deleted_at IS NULL
	`
	args := []any{internalProductID, code}
	if excludeID != "" {
		query += ` AND id != ?`
		args = append(args, excludeID)
	}

	var count int
	if err := r.db.GetContext(ctx, &count, r.db.Rebind(query), args...); err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, map[string]any{
			"internal_product_id": internalProductID,
			"code":                code,
			"exclude_id":          excludeID,
		}).Msg("Failed to check internal product pricing code")
		return false, err
	}

	return count > 0, nil
}

func (r *internalProductRepo) GetInternalProductPrices(
	ctx context.Context,
	filter coreentity.InternalProductPriceListFilter,
) ([]coreentity.InternalProductPrice, int, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalproduct:repo:GetInternalProductPrices",
	)
	defer span.End()

	type dao struct {
		TotalData                int             `db:"total_data"`
		ID                       string          `db:"id"`
		InternalProductPricingID string          `db:"internal_product_pricing_id"`
		CurrencyCode             string          `db:"currency_code"`
		Amount                   decimal.Decimal `db:"amount"`
		StartedAt                string          `db:"started_at"`
		EndedAt                  *string         `db:"ended_at"`
		Metadata                 json.RawMessage `db:"metadata"`
		CreatedAt                string          `db:"created_at"`
		UpdatedAt                *string         `db:"updated_at"`
	}

	var (
		rows  = make([]dao, 0)
		items = make([]coreentity.InternalProductPrice, 0)
		args  = []any{filter.InternalProductPricingID}
		total int
	)

	query := `
		SELECT
			COUNT(*) OVER() AS total_data,
			id,
			internal_product_pricing_id,
			currency_code,
			amount::text AS amount,
			started_at::text AS started_at,
			ended_at::text AS ended_at,
			metadata,
			created_at::text AS created_at,
			updated_at::text AS updated_at
		FROM internal_product_prices
		WHERE internal_product_pricing_id = ?
			AND deleted_at IS NULL
	`
	if filter.CurrencyCode != "" {
		query += ` AND currency_code = ?`
		args = append(args, filter.CurrencyCode)
	}
	query += ` ORDER BY currency_code ASC, started_at DESC LIMIT ? OFFSET ?`
	args = append(args, filter.Paginate, (filter.Page-1)*filter.Paginate)

	if err := r.db.SelectContext(ctx, &rows, r.db.Rebind(query), args...); err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to query internal product prices")
		return nil, 0, err
	}

	for _, row := range rows {
		total = row.TotalData
		item, err := mapInternalProductPriceDAO(
			row.ID,
			row.InternalProductPricingID,
			row.CurrencyCode,
			row.Amount,
			row.StartedAt,
			row.EndedAt,
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

func (r *internalProductRepo) CreateInternalProductPrice(
	ctx context.Context,
	data coreentity.InternalProductPrice,
) (*coreentity.InternalProductPrice, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalproduct:repo:CreateInternalProductPrice",
	)
	defer span.End()

	metadata, err := marshalMetadata(data.Metadata)
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO internal_product_prices (
			internal_product_pricing_id,
			currency_code,
			amount,
			started_at,
			ended_at,
			metadata
		)
		VALUES (
			?,
			?,
			?,
			?,
			?,
			?
		)
		RETURNING
			id
	`
	if err := r.db.GetContext(ctx, &data.ID, r.db.Rebind(query),
		data.InternalProductPricingID,
		data.CurrencyCode,
		data.Amount,
		data.StartedAt,
		data.EndedAt,
		metadata,
	); err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to create internal product price")
		return nil, err
	}

	return &data, nil
}

func (r *internalProductRepo) UpdateInternalProductPrice(
	ctx context.Context,
	data coreentity.InternalProductPrice,
) error {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalproduct:repo:UpdateInternalProductPrice",
	)
	defer span.End()

	metadata, err := marshalMetadata(data.Metadata)
	if err != nil {
		return err
	}

	query := `
		UPDATE internal_product_prices
		SET
			internal_product_pricing_id = ?,
			currency_code = ?,
			amount = ?,
			started_at = ?,
			ended_at = ?,
			metadata = ?,
			updated_at = NOW()
		WHERE id = ?
			AND deleted_at IS NULL
	`
	result, err := r.db.ExecContext(ctx, r.db.Rebind(query),
		data.InternalProductPricingID,
		data.CurrencyCode,
		data.Amount,
		data.StartedAt,
		data.EndedAt,
		metadata,
		data.ID,
	)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, data).Msg("Failed to update internal product price")
		return err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		log.Ctx(ctx).Warn().Any(common.LogKeyPayload, data).Msg("Internal product price not found when updating")
		return errmsg.NewCustomErrors(404).SetMessage(errmsg.MessageInternalProductPriceNotFound)
	}
	return nil
}

func (r *internalProductRepo) DeleteInternalProductPrice(
	ctx context.Context,
	filter coreentity.InternalProductPriceDeleteFilter,
) error {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalproduct:repo:DeleteInternalProductPrice",
	)
	defer span.End()

	query := `
		UPDATE internal_product_prices
		SET deleted_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`
	result, err := r.db.ExecContext(ctx, r.db.Rebind(query), filter.ID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, filter).Msg("Failed to delete internal product price")
		return err
	}

	if affected, _ := result.RowsAffected(); affected == 0 {
		log.Ctx(ctx).Warn().Any(common.LogKeyPayload, filter).Msg("Internal product price not found when deleting")
		return errmsg.NewCustomErrors(404).SetMessage(errmsg.MessageInternalProductPriceNotFound)
	}
	return nil
}

func (r *internalProductRepo) ExistsOverlappingInternalProductPrice(
	ctx context.Context,
	filter coreentity.InternalProductPriceOverlapFilter,
) (bool, error) {
	ctx, span := tracing.StartSpan(
		ctx,
		"internal:framework:secondary:db:postgres:internalproduct:repo:ExistsOverlappingInternalProductPrice",
	)
	defer span.End()

	query := `
		SELECT COUNT(*)
		FROM internal_product_prices
		WHERE internal_product_pricing_id = ?
		  AND currency_code = ?
		  AND deleted_at IS NULL
		  AND started_at < COALESCE(?, 'infinity'::timestamptz)
		  AND COALESCE(ended_at, 'infinity'::timestamptz) > ?
	`
	args := []any{filter.InternalProductPricingID, filter.CurrencyCode, filter.EndedAt, filter.StartedAt}
	if filter.ExcludeID != "" {
		query += ` AND id != ?`
		args = append(args, filter.ExcludeID)
	}

	var count int
	if err := r.db.GetContext(ctx, &count, r.db.Rebind(query), args...); err != nil {
		log.Ctx(ctx).
			Error().
			Err(err).
			Any(common.LogKeyPayload, filter).
			Msg("Failed to check overlapping internal product price")
		return false, err
	}

	return count > 0, nil
}

func marshalMetadata(metadata map[string]any) ([]byte, error) {
	if metadata == nil {
		metadata = map[string]any{}
	}

	return json.Marshal(metadata)
}

func parseMetadata(ctx context.Context, raw json.RawMessage) (map[string]any, error) {
	if len(raw) == 0 {
		return map[string]any{}, nil
	}

	var metadata map[string]any
	if err := json.Unmarshal(raw, &metadata); err != nil {
		log.Ctx(ctx).Error().Err(err).Any(common.LogKeyPayload, string(raw)).Msg("Failed to unmarshal metadata")
		return nil, err
	}
	if metadata == nil {
		metadata = map[string]any{}
	}

	return metadata, nil
}

func mapInternalProductDAO(
	id string,
	code string,
	name string,
	description string,
	status string,
	metadataRaw json.RawMessage,
	createdAt string,
	updatedAt *string,
) (coreentity.InternalProduct, error) {
	metadata, err := parseMetadata(context.Background(), metadataRaw)
	if err != nil {
		return coreentity.InternalProduct{}, err
	}

	item := coreentity.InternalProduct{
		ID:          id,
		Code:        code,
		Name:        name,
		Description: description,
		Status:      status,
		Metadata:    metadata,
		CreatedAt:   createdAt,
	}
	if updatedAt != nil {
		item.UpdatedAt = *updatedAt
	}

	return item, nil
}

func mapInternalProductPricingDAO(
	id string,
	internalProductID string,
	code string,
	name string,
	description string,
	status string,
	metadataRaw json.RawMessage,
	createdAt string,
	updatedAt *string,
) (coreentity.InternalProductPricing, error) {
	metadata, err := parseMetadata(context.Background(), metadataRaw)
	if err != nil {
		return coreentity.InternalProductPricing{}, err
	}

	item := coreentity.InternalProductPricing{
		ID:                id,
		InternalProductID: internalProductID,
		Code:              code,
		Name:              name,
		Description:       description,
		Status:            status,
		Metadata:          metadata,
		CreatedAt:         createdAt,
	}
	if updatedAt != nil {
		item.UpdatedAt = *updatedAt
	}

	return item, nil
}

func mapInternalProductPriceDAO(
	id string,
	pricingID string,
	currencyCode string,
	amount decimal.Decimal,
	startedAt string,
	endedAt *string,
	metadataRaw json.RawMessage,
	createdAt string,
	updatedAt *string,
) (coreentity.InternalProductPrice, error) {
	metadata, err := parseMetadata(context.Background(), metadataRaw)
	if err != nil {
		return coreentity.InternalProductPrice{}, err
	}

	item := coreentity.InternalProductPrice{
		ID:                       id,
		InternalProductPricingID: pricingID,
		CurrencyCode:             currencyCode,
		Amount:                   amount,
		StartedAt:                startedAt,
		EndedAt:                  endedAt,
		Metadata:                 metadata,
		CreatedAt:                createdAt,
	}
	if updatedAt != nil {
		item.UpdatedAt = *updatedAt
	}

	return item, nil
}
