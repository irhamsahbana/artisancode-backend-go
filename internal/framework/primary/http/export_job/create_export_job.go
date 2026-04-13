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

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (h *exportJobHandler) createExportJob(c *fiber.Ctx) error {
	tracedCtx, span := tracing.StartSpan(c.UserContext(), "internal:framework:primary:http:export_job:create_export_job:createExportJob")
	defer span.End()
	c.SetUserContext(tracedCtx)

	var (
		ctx = c.UserContext()
		req = new(restentity.CreateExportJobReq)
		v   = adapter.Adapters.Validator
	)

	if err := c.BodyParser(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Failed to parse export job request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := v.Validate(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Failed to validate export job request body")
		code, errors := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errors))
	}

	paramsJSON, err := mapper.ExportJobParamsToJSON(*req)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to serialize export job params")
		return c.Status(fiber.StatusInternalServerError).JSON(response.Error(err))
	}

	uc := common.GetUserContext(ctx)
	resourceLabel := "Attendance Logs"
	processorKey := "attendance_logs"
	item, err := h.core.CreateExportJob(ctx, coreentity.ExportJobCreate{
		UserCtx:       uc,
		TenantID:      uc.TenantID,
		RequestedBy:   uc.UserID,
		ResourceType:  req.ResourceType,
		ResourceLabel: resourceLabel,
		ProcessorKey:  processorKey,
		Format:        coreentity.ExportJobFormat(req.Format),
		ParamsJSON:    paramsJSON,
	})
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to create export job")
		code, errors := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errors))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(restentity.CreateExportJobResp{
		ID: item.ID,
	}, ""))
}
