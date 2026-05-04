package handler

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/entity/mapper"
	"codebase-app/internal/entity/restentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"
	"codebase-app/pkg/types"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (h *internalCurrencyHandler) getProviderCurrencies(c *fiber.Ctx) error {
	tracedCtx, span := tracing.StartSpan(c.UserContext(), "internal:framework:primary:http:internalcurrency:get_provider_currencies:getProviderCurrencies")
	defer span.End()
	c.SetUserContext(tracedCtx)

	ctx := c.UserContext()
	req := new(restentity.GetInternalPaymentProviderCurrenciesReq)
	req.SetDefault()
	if err := c.ParamsParser(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Failed to parse provider currency params")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}
	if err := c.QueryParser(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Failed to parse provider currency query")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}
	if err := adapter.Adapters.Validator.Validate(req); err != nil {
		code, errors := errmsg.Errors(ctx, err, req)
		return c.Status(code).JSON(response.Error(errors))
	}

	items, total, err := h.core.GetProviderCurrencies(ctx, mapperInternalProviderCurrencyListFilter(*req))
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get provider currencies")
		code, errors := errmsg.Errors[error](ctx, err)
		return c.Status(code).JSON(response.Error(errors))
	}

	respItems := make([]restentity.InternalPaymentProviderCurrency, 0, len(items))
	for _, item := range items {
		respItems = append(respItems, mapper.InternalPaymentProviderCurrencyFromCoreToRest(item))
	}

	resp := restentity.GetInternalPaymentProviderCurrenciesResp{
		Items: respItems,
		Meta: types.Meta{
			Page:      req.Page,
			Paginate:  req.Paginate,
			TotalData: total,
		},
	}
	resp.Meta.CountTotalPage(req.Page, req.Paginate, total)

	return c.Status(fiber.StatusOK).JSON(response.Success(resp, ""))
}
