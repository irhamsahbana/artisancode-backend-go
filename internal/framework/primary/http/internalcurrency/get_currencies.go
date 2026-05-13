package handler

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/entity/mapper"
	"codebase-app/internal/entity/restentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"
	"codebase-app/pkg/types"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
)

func (h *internalCurrencyHandler) getInternalCurrencies(c fiber.Ctx) error {
	tracedCtx, span := tracing.StartSpan(c.Context(), "internal:framework:primary:http:internalcurrency:get_currencies:getInternalCurrencies")
	defer span.End()
	c.SetContext(tracedCtx)

	ctx := c.Context()
	req := new(restentity.GetInternalCurrenciesReq)
	req.SetDefault()

	if err := c.Bind().Query(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Failed to parse internal currency query")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}
	if err := adapter.Adapters.Validator.Validate(req); err != nil {
		code, errors := errmsg.Errors(ctx, err, req)
		return c.Status(code).JSON(response.Error(errors))
	}

	items, total, err := h.core.GetInternalCurrencies(ctx, mapperInternalCurrencyListFilter(*req))
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get internal currencies")
		code, errors := errmsg.Errors[error](ctx, err)
		return c.Status(code).JSON(response.Error(errors))
	}

	respItems := make([]restentity.InternalCurrency, 0, len(items))
	for _, item := range items {
		respItems = append(respItems, mapper.InternalCurrencyFromCoreToRest(item))
	}

	resp := restentity.GetInternalCurrenciesResp{
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
