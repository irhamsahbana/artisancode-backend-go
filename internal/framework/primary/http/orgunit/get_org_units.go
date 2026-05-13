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

func (h *orgUnitHandler) getOrgUnits(c fiber.Ctx) error {
	tracedCtx, span := tracing.StartSpan(c.Context(), "internal:framework:primary:http:orgunit:handler:getOrgUnits")
	defer span.End()
	c.SetContext(tracedCtx)

	var (
		ctx = c.Context()
		req = new(restentity.GetOrgUnitsReq)
		v   = adapter.Adapters.Validator
	)

	if err := c.Bind().Query(req); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to parse query params")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	req.SetDefault()

	if err := v.Validate(req); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to validate query params")
		code, errors := errmsg.Errors(ctx, err, req)
		return c.Status(code).JSON(response.Error(errors))
	}

	filter := coreentity.OrgUnitListFilter{
		TenantID: common.GetUserContext(ctx).TenantID,
		Q:        req.Q,
		Category: req.Category,
		Page:     req.Page,
		Paginate: req.Paginate,
	}

	items, total, err := h.core.GetOrgUnits(ctx, filter)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get org units")
		code, errors := errmsg.Errors[error](ctx, err)
		return c.Status(code).JSON(response.Error(errors))
	}

	restItems := make([]restentity.OrgUnit, 0, len(items))
	for _, item := range items {
		restItems = append(restItems, mapper.OrgUnitFromCoreToRest(item))
	}

	resp := &restentity.GetOrgUnitsResp{
		Items: restItems,
		Meta: types.Meta{
			Page:      req.Page,
			Paginate:  req.Paginate,
			TotalData: total,
		},
	}
	resp.Meta.CountTotalPage(req.Page, req.Paginate, total)

	return c.Status(fiber.StatusOK).JSON(response.Success(resp, ""))
}
