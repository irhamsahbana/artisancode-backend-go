package handler

import (
	"codebase-app/internal/entity/mapper"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
)

func (h *internalTenantBillingHandler) getEntitlements(c fiber.Ctx) error {
	tracedCtx, span := tracing.StartSpan(c.Context(), "internal:framework:primary:http:internaltenantbilling:get_entitlements:getEntitlements")
	defer span.End()
	c.SetContext(tracedCtx)

	ctx := c.Context()
	item, err := h.core.GetEntitlements(ctx)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get tenant billing entitlements")
		code, errors := errmsg.Errors[error](ctx, err)
		return c.Status(code).JSON(response.Error(errors))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(mapper.TenantBillingEntitlementsFromCoreToRest(*item), ""))
}
