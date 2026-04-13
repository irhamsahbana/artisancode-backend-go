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

func (h *exportJobHandler) getExportJobs(c *fiber.Ctx) error {
	tracedCtx, span := tracing.StartSpan(c.UserContext(), "internal:framework:primary:http:export_job:get_export_jobs:getExportJobs")
	defer span.End()
	c.SetUserContext(tracedCtx)

	var (
		ctx = c.UserContext()
		req = new(restentity.GetExportJobsReq)
		v   = adapter.Adapters.Validator
	)

	if err := c.QueryParser(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Failed to parse export job query params")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	req.SetDefault()

	if err := v.Validate(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Failed to validate export job query params")
		code, errors := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errors))
	}

	uc := common.GetUserContext(ctx)
	items, total, err := h.core.GetExportJobs(ctx, coreentity.ExportJobListFilter{
		UserCtx:  uc,
		TenantID: uc.TenantID,
		Page:     req.Page,
		Paginate: req.Limit,
	})
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to get export jobs")
		code, errors := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errors))
	}

	restItems := make([]restentity.ExportJob, 0, len(items))
	for _, item := range items {
		restItems = append(restItems, mapper.ExportJobFromCoreToRest(item))
	}

	meta := types.Meta{
		Page:      req.Page,
		Paginate:  req.Limit,
		TotalData: total,
	}
	meta.CountTotalPage(req.Page, req.Limit, total)

	return c.Status(fiber.StatusOK).JSON(response.Success(restentity.GetExportJobsResp{
		Items:      restItems,
		Pagination: restentity.NewPaginationResp(meta),
	}, ""))
}
