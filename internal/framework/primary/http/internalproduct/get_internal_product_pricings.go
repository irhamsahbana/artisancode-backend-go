package handler

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/entity/mapper"
	"codebase-app/internal/entity/restentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"
	"codebase-app/pkg/types"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
)

func (h *internalProductHandler) getInternalProductPricings(c fiber.Ctx) error {
	tracedCtx, span := tracing.StartSpan(c.Context(), "internal:framework:primary:http:internalproduct:handler:getInternalProductPricings")
	defer span.End()
	c.SetContext(tracedCtx)

	var (
		ctx = c.Context()
		req = new(restentity.GetInternalProductPricingsReq)
		v   = adapter.Adapters.Validator
	)

	if err := c.Bind().URI(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Failed to parse internal product pricing params")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}
	if err := c.Bind().Query(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Failed to parse internal product pricing query")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}
	req.SetDefault()

	if err := v.Validate(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Failed to validate internal product pricing query")
		code, errors := errmsg.Errors(ctx, err, req)
		return c.Status(code).JSON(response.Error(errors))
	}

	items, total, err := h.core.GetInternalProductPricings(ctx, coreentity.InternalProductPricingListFilter{
		InternalProductID: req.ID,
		Q:                 req.Q,
		Page:              req.Page,
		Paginate:          req.Paginate,
	})
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get internal product pricings")
		code, errors := errmsg.Errors[error](ctx, err)
		return c.Status(code).JSON(response.Error(errors))
	}

	restItems := make([]restentity.InternalProductPricing, 0, len(items))
	for _, item := range items {
		restItems = append(restItems, mapper.InternalProductPricingFromCoreToRest(item))
	}

	resp := &restentity.GetInternalProductPricingsResp{
		Items: restItems,
		Meta: types.Meta{
			Page:      req.Page,
			Paginate:  req.Paginate,
			TotalData: total,
		},
	}
	resp.Meta.CountTotalPage(req.Page, req.Paginate, total)

	return c.Status(fiber.StatusOK).JSON(response.Success(resp, ""))
}
