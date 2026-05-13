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

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
)

func (h *rbacHandler) getPermissions(c fiber.Ctx) error {
	tracedCtx, span := tracing.StartSpan(c.Context(), "internal:framework:primary:http:rbac:handler:getPermissions")
	defer span.End()
	c.SetContext(tracedCtx)

	var (
		ctx = c.Context()
		req = new(restentity.GetPermissionsReq)
		v   = adapter.Adapters.Validator
	)

	if err := c.Bind().Query(req); err != nil {
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
	filter := coreentity.PermissionListFilter{
		TenantID: uc.TenantID,
		Q:        req.Q,
		Page:     req.Page,
		Paginate: req.Limit,
	}

	permissions, total, err := h.core.GetPermissions(ctx, filter)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get permissions")
		code, errors := errmsg.Errors[error](ctx, err)
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
