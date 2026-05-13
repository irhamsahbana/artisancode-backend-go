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

func (h *internalTenantBillingHandler) getInvoice(c fiber.Ctx) error {
	tracedCtx, span := tracing.StartSpan(
		c.Context(),
		"internal:framework:primary:http:internaltenantbilling:get_invoice:getInvoice",
	)
	defer span.End()
	c.SetContext(tracedCtx)

	ctx := c.Context()
	req := new(restentity.GetTenantBillingResourceReq)
	if err := c.Bind().URI(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Failed to parse tenant billing invoice params")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}
	if err := adapter.Adapters.Validator.Validate(req); err != nil {
		code, errors := errmsg.Errors(ctx, err, req)
		return c.Status(code).JSON(response.Error(errors))
	}

	detail, err := h.core.GetInvoice(ctx, req.ID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get tenant billing invoice")
		code, errors := errmsg.Errors[error](ctx, err)
		return c.Status(code).JSON(response.Error(errors))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(mapper.TenantBillingInvoiceDetailFromCoreToRest(*detail), ""))
}
