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

func (h *jobPositionHandler) getJobPositions(c fiber.Ctx) error {
	tracedCtx, span := tracing.StartSpan(c.Context(), "internal:framework:primary:http:jobposition:handler:getJobPositions")
	defer span.End()
	c.SetContext(tracedCtx)

	var (
		ctx = c.Context()
		req = new(restentity.GetJobPositionsReq)
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

	filter := coreentity.JobPositionListFilter{
		TenantID: common.GetUserContext(ctx).TenantID,
		Q:        req.Q,
		Page:     req.Page,
		Paginate: req.Paginate,
	}

	items, total, err := h.core.GetJobPositions(ctx, filter)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get job positions")
		code, errors := errmsg.Errors[error](ctx, err)
		return c.Status(code).JSON(response.Error(errors))
	}

	restItems := make([]restentity.JobPosition, 0, len(items))
	for _, item := range items {
		restItems = append(restItems, mapper.JobPositionFromCoreToRest(item))
	}

	resp := &restentity.GetJobPositionsResp{
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
