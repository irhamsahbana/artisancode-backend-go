package handler

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/mapper"
	"codebase-app/internal/entity/restentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
)

func (h *rbacHandler) createRole(c fiber.Ctx) error {
	tracedCtx, span := tracing.StartSpan(c.Context(), "internal:framework:primary:http:rbac:handler:createRole")
	defer span.End()
	c.SetContext(tracedCtx)

	var (
		ctx = c.Context()
		req = new(restentity.CreateRoleReq)
		v   = adapter.Adapters.Validator
	)

	if err := c.Bind().Body(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := v.Validate(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Failed to validate request body")
		code, errors := errmsg.Errors(ctx, err, req)
		return c.Status(code).JSON(response.Error(errors))
	}

	data := mapper.RoleFromRestCreateToCore(ctx, *req)
	created, err := h.core.CreateRole(ctx, data)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to create role")
		code, errors := errmsg.Errors[error](ctx, err)
		return c.Status(code).JSON(response.Error(errors))
	}

	if len(req.Permissions) > 0 {
		uc := common.GetUserContext(ctx)
		if err := h.core.SetRolePermissions(ctx, created.ID, uc.TenantID, req.Permissions); err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Failed to set role permissions")
			code, errors := errmsg.Errors[error](ctx, err)
			return c.Status(code).JSON(response.Error(errors))
		}
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(restentity.CreateRoleResp{ID: created.ID}, ""))
}
