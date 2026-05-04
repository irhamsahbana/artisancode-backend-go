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

func (h *internalCurrencyHandler) updateInternalCurrency(c *fiber.Ctx) error {
	tracedCtx, span := tracing.StartSpan(c.UserContext(), "internal:framework:primary:http:internalcurrency:update_currency:updateInternalCurrency")
	defer span.End()
	c.SetUserContext(tracedCtx)

	ctx := c.UserContext()
	req := new(restentity.UpdateInternalCurrencyReq)
	if err := c.ParamsParser(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Failed to parse internal currency params")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}
	if err := c.BodyParser(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Failed to parse internal currency body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}
	if err := adapter.Adapters.Validator.Validate(req); err != nil {
		code, errors := errmsg.Errors(ctx, err, req)
		return c.Status(code).JSON(response.Error(errors))
	}

	if err := h.core.UpdateInternalCurrency(ctx, mapper.InternalCurrencyFromRestUpdateToCore(ctx, *req)); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to update internal currency")
		code, errors := errmsg.Errors[error](ctx, err)
		return c.Status(code).JSON(response.Error(errors))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(nil, ""))
}
