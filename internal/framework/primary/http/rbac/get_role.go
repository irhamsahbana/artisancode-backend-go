package handler

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/mapper"
	"codebase-app/internal/entity/restentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (h *rbacHandler) getRole(c *fiber.Ctx) error {
	tracedCtx, span := tracing.StartSpan(c.UserContext(), "internal:framework:primary:http:rbac:handler:getRole")
	defer span.End()
	c.SetUserContext(tracedCtx)

	var (
		ctx = c.UserContext()
		req = new(restentity.GetRoleReq)
		v   = adapter.Adapters.Validator
	)

	if err := c.ParamsParser(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Failed to parse params")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := v.Validate(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Failed to validate params")
		code, errors := errmsg.Errors(ctx, err, req)
		return c.Status(code).JSON(response.Error(errors))
	}

	uc := common.GetUserContext(ctx)
	role, err := h.core.GetRoleWithPermissions(ctx, req.ID, uc.TenantID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get role")
		code, errors := errmsg.Errors[error](ctx, err)
		return c.Status(code).JSON(response.Error(errors))
	}

	perms, err := h.core.GetPermissionsByRoleID(ctx, role.ID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get role permissions")
		code, errors := errmsg.Errors[error](ctx, err)
		return c.Status(code).JSON(response.Error(errors))
	}

	restPerms := make([]restentity.Permission, 0, len(perms))
	for _, p := range perms {
		restPerms = append(restPerms, mapper.PermissionFromCoreToRest(p))
	}

	resp := restentity.RoleWithPermissions{
		ID:          role.ID,
		Name:        role.Name,
		Permissions: restPerms,
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(resp, ""))
}
