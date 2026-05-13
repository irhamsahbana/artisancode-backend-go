package handler

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/entity/mapper"
	"codebase-app/internal/entity/restentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
)

func (h *internalProductHandler) createInternalProductPrice(c fiber.Ctx) error {
	tracedCtx, span := tracing.StartSpan(c.Context(), "internal:framework:primary:http:internalproduct:handler:createInternalProductPrice")
	defer span.End()
	c.SetContext(tracedCtx)

	var (
		ctx = c.Context()
		req = new(restentity.CreateInternalProductPriceReq)
		v   = adapter.Adapters.Validator
	)

	if err := c.Bind().URI(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Failed to parse internal product price params")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}
	if err := c.Bind().Body(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Failed to parse internal product price body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := v.Validate(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Failed to validate internal product price body")
		code, errors := errmsg.Errors(ctx, err, req)
		return c.Status(code).JSON(response.Error(errors))
	}

	item, err := h.core.CreateInternalProductPrice(ctx, mapper.InternalProductPriceFromRestCreateToCore(ctx, *req))
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to create internal product price")
		code, errors := errmsg.Errors[error](ctx, err)
		return c.Status(code).JSON(response.Error(errors))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(restentity.CreateInternalProductPriceResp{ID: item.ID}, ""))
}
