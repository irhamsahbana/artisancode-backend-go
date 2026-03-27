package handler

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/entity/mapper"
	"codebase-app/internal/entity/restentity"
	portsCore "codebase-app/internal/ports/core"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"
	"codebase-app/pkg/types"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type rbacHandler struct {
	core portsCore.RbacCore
}

// RbacHandlerConfig is the config struct for the RBAC handler
type RbacHandlerConfig struct {
	Core portsCore.RbacCore
}

func NewRbacHandler(cfg RbacHandlerConfig) *rbacHandler {
	return &rbacHandler{core: cfg.Core}
}

func (h *rbacHandler) Register(router fiber.Router) {
	router.Get("/roles", h.getRoles)
	router.Get("/roles/:id", h.getRole)
	router.Post("/roles", h.createRole)
	router.Put("/roles/:id", h.updateRole)
	router.Delete("/roles/:id", h.deleteRole)
	router.Get("/permissions", h.getPermissions)
}

func (h *rbacHandler) getRoles(c *fiber.Ctx) error {
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
		code, errors := errmsg.Errors(err, req)
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
		code, errors := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errors))
	}

	items := make([]restentity.RoleWithPermissions, 0, len(roles))
	for _, role := range roles {
		perms, err := h.core.GetPermissionsByRoleID(ctx, role.ID)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Failed to get role permissions")
			code, errors := errmsg.Errors[error](err)
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

func (h *rbacHandler) getRole(c *fiber.Ctx) error {
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
		code, errors := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errors))
	}

	uc := common.GetUserContext(ctx)
	role, err := h.core.GetRoleWithPermissions(ctx, req.ID, uc.TenantID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get role")
		code, errors := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errors))
	}

	perms, err := h.core.GetPermissionsByRoleID(ctx, role.ID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get role permissions")
		code, errors := errmsg.Errors[error](err)
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

func (h *rbacHandler) createRole(c *fiber.Ctx) error {
	var (
		ctx = c.UserContext()
		req = new(restentity.CreateRoleReq)
		v   = adapter.Adapters.Validator
	)

	if err := c.BodyParser(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := v.Validate(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Failed to validate request body")
		code, errors := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errors))
	}

	data := mapper.RoleFromRestCreateToCore(ctx, *req)
	created, err := h.core.CreateRole(ctx, data)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to create role")
		code, errors := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errors))
	}

	// Set permissions if provided
	if len(req.Permissions) > 0 {
		uc := common.GetUserContext(ctx)
		if err := h.core.SetRolePermissions(ctx, created.ID, uc.TenantID, req.Permissions); err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Failed to set role permissions")
			code, errors := errmsg.Errors[error](err)
			return c.Status(code).JSON(response.Error(errors))
		}
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(restentity.CreateRoleResp{ID: created.ID}, ""))
}

func (h *rbacHandler) updateRole(c *fiber.Ctx) error {
	var (
		ctx = c.UserContext()
		req = new(restentity.UpdateRoleReq)
		v   = adapter.Adapters.Validator
	)

	if err := c.ParamsParser(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Failed to parse params")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := c.BodyParser(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := v.Validate(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Failed to validate request body")
		code, errors := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errors))
	}

	data := mapper.RoleFromRestUpdateToCore(ctx, *req)
	if err := h.core.UpdateRole(ctx, data); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to update role")
		code, errors := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errors))
	}

	// Replace permissions
	uc := common.GetUserContext(ctx)
	if err := h.core.SetRolePermissions(ctx, req.ID, uc.TenantID, req.PermissionIDs); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to set role permissions")
		code, errors := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errors))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(nil, ""))
}

func (h *rbacHandler) deleteRole(c *fiber.Ctx) error {
	var (
		ctx = c.UserContext()
		req = new(restentity.DeleteRoleReq)
		v   = adapter.Adapters.Validator
	)

	if err := c.ParamsParser(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Failed to parse params")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := v.Validate(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Failed to validate params")
		code, errors := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errors))
	}

	uc := common.GetUserContext(ctx)
	filter := coreentity.RoleDeleteFilter{
		TenantID: uc.TenantID,
		ID:       req.ID,
	}

	if err := h.core.DeleteRole(ctx, filter); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to delete role")
		code, errors := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errors))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(nil, ""))
}

func (h *rbacHandler) getPermissions(c *fiber.Ctx) error {
	var (
		ctx = c.UserContext()
		req = new(restentity.GetPermissionsReq)
		v   = adapter.Adapters.Validator
	)

	if err := c.QueryParser(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Failed to parse query params")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	req.SetDefault()

	if err := v.Validate(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Failed to validate query params")
		code, errors := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errors))
	}

	uc := common.GetUserContext(ctx)
	filter := coreentity.PermissionListFilter{
		TenantID: uc.TenantID,
		Q:        req.Q,
		Page:     req.Page,
		Paginate: req.Limit,
	}

	permissions, total, err := h.core.GetPermissions(ctx, filter)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get permissions")
		code, errors := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errors))
	}

	items := make([]restentity.Permission, 0, len(permissions))
	for _, p := range permissions {
		items = append(items, mapper.PermissionFromCoreToRest(p))
	}

	meta := types.Meta{}
	meta.CountTotalPage(req.Page, req.Limit, total)

	resp := &restentity.GetPermissionsResp{
		Items:      items,
		Pagination: restentity.NewPaginationResp(meta),
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(resp, ""))
}