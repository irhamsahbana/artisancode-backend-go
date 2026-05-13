package handler

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/entity/mapper"
	"codebase-app/internal/entity/restentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
)

func (h *internalProductHandler) getInternalProduct(c fiber.Ctx) error {
	tracedCtx, span := tracing.StartSpan(c.Context(), "internal:framework:primary:http:internalproduct:handler:getInternalProduct")
	defer span.End()
	c.SetContext(tracedCtx)

	var (
		ctx = c.Context()
		req = new(restentity.GetInternalProductReq)
		v   = adapter.Adapters.Validator
	)

	if err := c.Bind().URI(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Failed to parse internal product params")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := v.Validate(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Failed to validate internal product params")
		code, errors := errmsg.Errors(ctx, err, req)
		return c.Status(code).JSON(response.Error(errors))
	}

	item, err := h.core.GetInternalProduct(ctx, coreentity.InternalProductFilter{ID: req.ID})
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get internal product")
		code, errors := errmsg.Errors[error](ctx, err)
		return c.Status(code).JSON(response.Error(errors))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(restentity.GetInternalProductResp{
		InternalProduct: mapper.InternalProductFromCoreToRest(*item),
	}, ""))
}
