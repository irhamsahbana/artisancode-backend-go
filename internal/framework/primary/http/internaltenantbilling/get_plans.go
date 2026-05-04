package handler

import (
	"codebase-app/internal/entity/mapper"
	"codebase-app/internal/entity/restentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (h *internalTenantBillingHandler) getPlans(c *fiber.Ctx) error {
	tracedCtx, span := tracing.StartSpan(c.UserContext(), "internal:framework:primary:http:internaltenantbilling:get_plans:getPlans")
	defer span.End()
	c.SetUserContext(tracedCtx)

	ctx := c.UserContext()
	items, err := h.core.GetPlans(ctx)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get tenant billing plans")
		code, errors := errmsg.Errors[error](ctx, err)
		return c.Status(code).JSON(response.Error(errors))
	}

	restItems := make([]restentity.TenantBillingPlan, 0, len(items))
	for _, item := range items {
		restItems = append(restItems, mapper.TenantBillingPlanFromCoreToRest(item))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(restItems, ""))
}
