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

func (h *internalTenantBillingHandler) getPaymentAttempts(c fiber.Ctx) error {
	tracedCtx, span := tracing.StartSpan(
		c.Context(),
		"internal:framework:primary:http:internaltenantbilling:get_payment_attempts:getPaymentAttempts",
	)
	defer span.End()
	c.SetContext(tracedCtx)

	ctx := c.Context()
	req := new(restentity.GetTenantBillingResourceReq)
	if err := c.Bind().URI(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Failed to parse tenant billing payment attempts params")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}
	if err := adapter.Adapters.Validator.Validate(req); err != nil {
		code, errors := errmsg.Errors(ctx, err, req)
		return c.Status(code).JSON(response.Error(errors))
	}

	items, err := h.core.GetPaymentAttempts(ctx, req.ID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get tenant billing payment attempts")
		code, errors := errmsg.Errors[error](ctx, err)
		return c.Status(code).JSON(response.Error(errors))
	}

	restItems := make([]restentity.TenantBillingPaymentAttempt, 0, len(items))
	for _, item := range items {
		restItems = append(restItems, mapper.TenantBillingPaymentAttemptFromCoreToRest(item))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(restentity.GetTenantBillingPaymentAttemptsResp{
		Items: restItems,
	}, ""))
}
