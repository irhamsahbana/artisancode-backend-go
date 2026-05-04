package handler

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/entity/mapper"
	"codebase-app/internal/entity/restentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (h *internalCurrencyHandler) createInternalCurrency(c *fiber.Ctx) error {
	tracedCtx, span := tracing.StartSpan(c.UserContext(), "internal:framework:primary:http:internalcurrency:create_currency:createInternalCurrency")
	defer span.End()
	c.SetUserContext(tracedCtx)

	ctx := c.UserContext()
	req := new(restentity.CreateInternalCurrencyReq)
	if err := c.BodyParser(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Failed to parse internal currency body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}
	if err := adapter.Adapters.Validator.Validate(req); err != nil {
		code, errors := errmsg.Errors(ctx, err, req)
		return c.Status(code).JSON(response.Error(errors))
	}

	item, err := h.core.CreateInternalCurrency(ctx, mapper.InternalCurrencyFromRestCreateToCore(ctx, *req))
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to create internal currency")
		code, errors := errmsg.Errors[error](ctx, err)
		return c.Status(code).JSON(response.Error(errors))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(mapper.InternalCurrencyFromCoreToRest(*item), ""))
}
