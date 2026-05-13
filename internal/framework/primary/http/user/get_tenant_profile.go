package handler

import (
	"codebase-app/internal/entity/mapper"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
)

func (h *userHandler) getTenantProfile(c fiber.Ctx) error {
	tracedCtx, span := tracing.StartSpan(
		c.Context(),
		"internal:framework:primary:http:user:get_tenant_profile:getTenantProfile",
	)
	defer span.End()
	c.SetContext(tracedCtx)

	ctx := c.Context()
	profile, err := h.core.GetTenantProfile(ctx)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Get tenant profile service error")
		code, errs := errmsg.Errors[error](ctx, err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(mapper.TenantProfileToResp(*profile), ""))
}
