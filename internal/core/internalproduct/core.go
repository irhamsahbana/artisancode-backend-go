package core

import (
	"context"
	"regexp"
	"strings"
	"time"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	corePorts "codebase-app/internal/ports/core"
	portsRepo "codebase-app/internal/ports/secondary/db"
	"codebase-app/pkg/errmsg"

	"github.com/rs/zerolog/log"
)

var (
	internalCodePattern     = regexp.MustCompile(`^[A-Z0-9_-]+$`)
	internalCurrencyPattern = regexp.MustCompile(`^[A-Z]{3}$`)
)

type internalProductCore struct {
	repo portsRepo.InternalProductRepository
}

type Config struct {
	Repo portsRepo.InternalProductRepository
}

var _ corePorts.InternalProductCore = &internalProductCore{}

func NewInternalProductCore(cfg Config) *internalProductCore {
	return &internalProductCore{repo: cfg.Repo}
}

func (c *internalProductCore) GetInternalProducts(ctx context.Context, filter coreentity.InternalProductListFilter) ([]coreentity.InternalProduct, int, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internalproduct:core:GetInternalProducts")
	defer span.End()

	return c.repo.GetInternalProducts(ctx, filter)
}

func (c *internalProductCore) GetInternalProduct(ctx context.Context, filter coreentity.InternalProductFilter) (*coreentity.InternalProduct, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internalproduct:core:GetInternalProduct")
	defer span.End()

	return c.repo.GetInternalProduct(ctx, filter)
}

func (c *internalProductCore) CreateInternalProduct(ctx context.Context, data coreentity.InternalProduct) (*coreentity.InternalProduct, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internalproduct:core:CreateInternalProduct")
	defer span.End()

	if err := c.authorizeWrite(data.UserCtx); err != nil {
		return nil, err
	}

	if err := c.normalizeAndValidateProduct(ctx, &data); err != nil {
		return nil, err
	}

	exists, err := c.repo.ExistsInternalProductByCode(ctx, data.Code, "")
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errmsg.NewCustomErrors(400).SetMessage("Internal product code already exists")
	}

	return c.repo.CreateInternalProduct(ctx, data)
}

func (c *internalProductCore) UpdateInternalProduct(ctx context.Context, data coreentity.InternalProduct) error {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internalproduct:core:UpdateInternalProduct")
	defer span.End()

	if err := c.authorizeWrite(data.UserCtx); err != nil {
		return err
	}

	if err := c.normalizeAndValidateProduct(ctx, &data); err != nil {
		return err
	}

	exists, err := c.repo.ExistsInternalProductByCode(ctx, data.Code, data.ID)
	if err != nil {
		return err
	}
	if exists {
		return errmsg.NewCustomErrors(400).SetMessage("Internal product code already exists")
	}

	return c.repo.UpdateInternalProduct(ctx, data)
}

func (c *internalProductCore) DeleteInternalProduct(ctx context.Context, filter coreentity.InternalProductDeleteFilter) error {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internalproduct:core:DeleteInternalProduct")
	defer span.End()

	if err := c.authorizeWrite(common.GetUserContext(ctx)); err != nil {
		return err
	}

	return c.repo.DeleteInternalProduct(ctx, filter)
}

func (c *internalProductCore) GetInternalProductPricings(ctx context.Context, filter coreentity.InternalProductPricingListFilter) ([]coreentity.InternalProductPricing, int, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internalproduct:core:GetInternalProductPricings")
	defer span.End()

	return c.repo.GetInternalProductPricings(ctx, filter)
}

func (c *internalProductCore) GetInternalProductPricing(ctx context.Context, filter coreentity.InternalProductPricingFilter) (*coreentity.InternalProductPricing, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internalproduct:core:GetInternalProductPricing")
	defer span.End()

	return c.repo.GetInternalProductPricing(ctx, filter)
}

func (c *internalProductCore) CreateInternalProductPricing(ctx context.Context, data coreentity.InternalProductPricing) (*coreentity.InternalProductPricing, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internalproduct:core:CreateInternalProductPricing")
	defer span.End()

	if err := c.authorizeWrite(data.UserCtx); err != nil {
		return nil, err
	}

	if _, err := c.repo.GetInternalProduct(ctx, coreentity.InternalProductFilter{ID: data.InternalProductID}); err != nil {
		return nil, err
	}

	if err := c.normalizeAndValidatePricing(ctx, &data); err != nil {
		return nil, err
	}

	exists, err := c.repo.ExistsInternalProductPricingByCode(ctx, data.InternalProductID, data.Code, "")
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errmsg.NewCustomErrors(400).SetMessage("Internal product pricing code already exists")
	}

	return c.repo.CreateInternalProductPricing(ctx, data)
}

func (c *internalProductCore) UpdateInternalProductPricing(ctx context.Context, data coreentity.InternalProductPricing) error {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internalproduct:core:UpdateInternalProductPricing")
	defer span.End()

	if err := c.authorizeWrite(data.UserCtx); err != nil {
		return err
	}

	if _, err := c.repo.GetInternalProduct(ctx, coreentity.InternalProductFilter{ID: data.InternalProductID}); err != nil {
		return err
	}

	if err := c.normalizeAndValidatePricing(ctx, &data); err != nil {
		return err
	}

	exists, err := c.repo.ExistsInternalProductPricingByCode(ctx, data.InternalProductID, data.Code, data.ID)
	if err != nil {
		return err
	}
	if exists {
		return errmsg.NewCustomErrors(400).SetMessage("Internal product pricing code already exists")
	}

	return c.repo.UpdateInternalProductPricing(ctx, data)
}

func (c *internalProductCore) DeleteInternalProductPricing(ctx context.Context, filter coreentity.InternalProductPricingDeleteFilter) error {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internalproduct:core:DeleteInternalProductPricing")
	defer span.End()

	if err := c.authorizeWrite(common.GetUserContext(ctx)); err != nil {
		return err
	}

	return c.repo.DeleteInternalProductPricing(ctx, filter)
}

func (c *internalProductCore) GetInternalProductPrices(ctx context.Context, filter coreentity.InternalProductPriceListFilter) ([]coreentity.InternalProductPrice, int, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internalproduct:core:GetInternalProductPrices")
	defer span.End()

	if filter.CurrencyCode != "" {
		filter.CurrencyCode = strings.ToUpper(strings.TrimSpace(filter.CurrencyCode))
	}

	return c.repo.GetInternalProductPrices(ctx, filter)
}

func (c *internalProductCore) CreateInternalProductPrice(ctx context.Context, data coreentity.InternalProductPrice) (*coreentity.InternalProductPrice, error) {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internalproduct:core:CreateInternalProductPrice")
	defer span.End()

	if err := c.authorizeWrite(data.UserCtx); err != nil {
		return nil, err
	}

	if _, err := c.repo.GetInternalProductPricing(ctx, coreentity.InternalProductPricingFilter{ID: data.InternalProductPricingID}); err != nil {
		return nil, err
	}

	if err := c.normalizeAndValidatePrice(ctx, &data); err != nil {
		return nil, err
	}

	overlap, err := c.repo.ExistsOverlappingInternalProductPrice(ctx, coreentity.InternalProductPriceOverlapFilter{
		InternalProductPricingID: data.InternalProductPricingID,
		CurrencyCode:             data.CurrencyCode,
		StartedAt:                data.StartedAt,
		EndedAt:                  data.EndedAt,
	})
	if err != nil {
		return nil, err
	}
	if overlap {
		return nil, errmsg.NewCustomErrors(400).SetMessage("Price period overlaps with an existing active price")
	}

	return c.repo.CreateInternalProductPrice(ctx, data)
}

func (c *internalProductCore) UpdateInternalProductPrice(ctx context.Context, data coreentity.InternalProductPrice) error {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internalproduct:core:UpdateInternalProductPrice")
	defer span.End()

	if err := c.authorizeWrite(data.UserCtx); err != nil {
		return err
	}

	if _, err := c.repo.GetInternalProductPricing(ctx, coreentity.InternalProductPricingFilter{ID: data.InternalProductPricingID}); err != nil {
		return err
	}

	if err := c.normalizeAndValidatePrice(ctx, &data); err != nil {
		return err
	}

	overlap, err := c.repo.ExistsOverlappingInternalProductPrice(ctx, coreentity.InternalProductPriceOverlapFilter{
		InternalProductPricingID: data.InternalProductPricingID,
		CurrencyCode:             data.CurrencyCode,
		StartedAt:                data.StartedAt,
		EndedAt:                  data.EndedAt,
		ExcludeID:                data.ID,
	})
	if err != nil {
		return err
	}
	if overlap {
		return errmsg.NewCustomErrors(400).SetMessage("Price period overlaps with an existing active price")
	}

	return c.repo.UpdateInternalProductPrice(ctx, data)
}

func (c *internalProductCore) DeleteInternalProductPrice(ctx context.Context, filter coreentity.InternalProductPriceDeleteFilter) error {
	ctx, span := tracing.StartSpan(ctx, "internal:core:internalproduct:core:DeleteInternalProductPrice")
	defer span.End()

	if err := c.authorizeWrite(common.GetUserContext(ctx)); err != nil {
		return err
	}

	return c.repo.DeleteInternalProductPrice(ctx, filter)
}

func (c *internalProductCore) authorizeWrite(userCtx common.UserContext) error {
	if !userCtx.HasRole(coreentity.InternalUserRoleSuperAdmin) && !userCtx.HasRole(coreentity.InternalUserRoleOperator) {
		return errmsg.NewCustomErrors(403).SetMessage("You are not authorized to manage internal products")
	}

	return nil
}

func (c *internalProductCore) normalizeAndValidateProduct(ctx context.Context, data *coreentity.InternalProduct) error {
	data.Code = strings.ToUpper(strings.TrimSpace(data.Code))
	data.Name = strings.TrimSpace(data.Name)
	data.Description = strings.TrimSpace(data.Description)
	data.Status = strings.ToLower(strings.TrimSpace(data.Status))
	if data.Metadata == nil {
		data.Metadata = map[string]any{}
	}

	if !internalCodePattern.MatchString(data.Code) {
		return errmsg.NewCustomErrors(400).SetMessage("Internal product code format is invalid")
	}
	if data.Status != coreentity.InternalProductStatusDraft &&
		data.Status != coreentity.InternalProductStatusActive &&
		data.Status != coreentity.InternalProductStatusInactive &&
		data.Status != coreentity.InternalProductStatusArchived {
		return errmsg.NewCustomErrors(400).SetMessage("Internal product status is invalid")
	}

	return nil
}

func (c *internalProductCore) normalizeAndValidatePricing(ctx context.Context, data *coreentity.InternalProductPricing) error {
	data.Code = strings.ToUpper(strings.TrimSpace(data.Code))
	data.Name = strings.TrimSpace(data.Name)
	data.Description = strings.TrimSpace(data.Description)
	data.Status = strings.ToLower(strings.TrimSpace(data.Status))
	if data.Metadata == nil {
		data.Metadata = map[string]any{}
	}

	if !internalCodePattern.MatchString(data.Code) {
		return errmsg.NewCustomErrors(400).SetMessage("Internal product pricing code format is invalid")
	}
	if data.Status != coreentity.InternalProductStatusDraft &&
		data.Status != coreentity.InternalProductStatusActive &&
		data.Status != coreentity.InternalProductStatusInactive &&
		data.Status != coreentity.InternalProductStatusArchived {
		return errmsg.NewCustomErrors(400).SetMessage("Internal product pricing status is invalid")
	}

	return nil
}

func (c *internalProductCore) normalizeAndValidatePrice(ctx context.Context, data *coreentity.InternalProductPrice) error {
	data.CurrencyCode = strings.ToUpper(strings.TrimSpace(data.CurrencyCode))
	data.StartedAt = strings.TrimSpace(data.StartedAt)
	if data.EndedAt != nil {
		trimmed := strings.TrimSpace(*data.EndedAt)
		data.EndedAt = &trimmed
		if trimmed == "" {
			data.EndedAt = nil
		}
	}
	if data.Metadata == nil {
		data.Metadata = map[string]any{}
	}

	if !internalCurrencyPattern.MatchString(data.CurrencyCode) {
		return errmsg.NewCustomErrors(400).SetMessage("Currency code format is invalid")
	}

	if !data.Amount.IsPositive() {
		return errmsg.NewCustomErrors(400).SetMessage("Amount must be positive")
	}

	startTime, err := time.Parse(time.RFC3339, data.StartedAt)
	if err != nil {
		log.Ctx(ctx).Warn().Err(err).Any(common.LogKeyPayload, data.StartedAt).Msg("Invalid started_at")
		return errmsg.NewCustomErrors(400).SetMessage("started_at must use RFC3339 format")
	}

	if data.EndedAt != nil {
		endTime, err := time.Parse(time.RFC3339, *data.EndedAt)
		if err != nil {
			log.Ctx(ctx).Warn().Err(err).Any(common.LogKeyPayload, data.EndedAt).Msg("Invalid ended_at")
			return errmsg.NewCustomErrors(400).SetMessage("ended_at must use RFC3339 format")
		}
		if !endTime.After(startTime) {
			return errmsg.NewCustomErrors(400).SetMessage("ended_at must be later than started_at")
		}
	}

	return nil
}
