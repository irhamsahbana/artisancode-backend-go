package handler

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/entity/mapper"
	"codebase-app/internal/entity/restentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"
	"codebase-app/pkg/types"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (h *rbacHandler) getRoles(c *fiber.Ctx) error {
	tracedCtx, span := tracing.StartSpan(c.UserContext(), "internal:framework:primary:http:rbac:handler:getRoles")
	defer span.End()
	c.SetUserContext(tracedCtx)

	var (
		ctx = c.UserContext()
		req = new(restentity.GetRolesReq)
		v   = adapter.Adapters.Validator
	)

	if err := c.QueryParser(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Failed to parse query params")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	req.SetDefault()

	if err := v.Validate(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Failed to validate query params")
		code, errors := errmsg.Errors(ctx, err, req)
		return c.Status(code).JSON(response.Error(errors))
	}

	uc := common.GetUserContext(ctx)
	filter := coreentity.RoleListFilter{
		TenantID: uc.TenantID,
		Q:        req.Q,
		Page:     req.Page,
		Paginate: req.Limit,
	}

	roles, total, err := h.core.GetRoles(ctx, filter)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get roles")
		code, errors := errmsg.Errors[error](ctx, err)
		return c.Status(code).JSON(response.Error(errors))
	}

	items := make([]restentity.RoleWithPermissions, 0, len(roles))
	for _, role := range roles {
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

		items = append(items, restentity.RoleWithPermissions{
			ID:          role.ID,
			Name:        role.Name,
			Permissions: restPerms,
		})
	}

	meta := types.Meta{}
	meta.CountTotalPage(req.Page, req.Limit, total)

	resp := &restentity.GetRolesResp{
		Items:      items,
		Pagination: restentity.NewPaginationResp(meta),
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(resp, ""))
}
