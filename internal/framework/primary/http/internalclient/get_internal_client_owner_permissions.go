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

func (h *internalClientHandler) getInternalClientOwnerPermissions(c *fiber.Ctx) error {
	tracedCtx, span := tracing.StartSpan(c.UserContext(), "internal:framework:primary:http:internalclient:handler:getInternalClientOwnerPermissions")
	defer span.End()
	c.SetUserContext(tracedCtx)

	var (
		ctx = c.UserContext()
		req = new(restentity.GetInternalClientOwnerPermissionsReq)
		v   = adapter.Adapters.Validator
	)

	if err := c.ParamsParser(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Failed to parse internal client owner permissions params")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := v.Validate(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Failed to validate internal client owner permissions params")
		code, errors := errmsg.Errors(ctx, err, req)
		return c.Status(code).JSON(response.Error(errors))
	}

	item, err := h.core.GetInternalClientOwnerPermissions(ctx, req.ID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get internal client owner permissions")
		code, errors := errmsg.Errors[error](ctx, err)
		return c.Status(code).JSON(response.Error(errors))
	}

	available := make([]restentity.Permission, 0, len(item.Available))
	for _, permission := range item.Available {
		available = append(available, mapper.PermissionFromCoreToRest(permission))
	}

	resp := &restentity.GetInternalClientOwnerPermissionsResp{
		ClientID:           item.ClientID,
		Available:          available,
		OwnerPermissionIDs: item.OwnerPermissionIDs,
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(resp, ""))
}
