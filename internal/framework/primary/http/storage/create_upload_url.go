package handler

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/entity/restentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (h *storageHandler) createUploadURL(c *fiber.Ctx) error {
	tracedCtx, span := tracing.StartSpan(c.UserContext(), "internal:framework:primary:http:storage:create_upload_url:createUploadURL")
	defer span.End()
	c.SetUserContext(tracedCtx)

	var (
		ctx = c.UserContext()
		req = new(restentity.CreateUploadURLReq)
		v   = adapter.Adapters.Validator
		uc  = common.GetUserContext(ctx)
	)

	if err := c.BodyParser(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Invalid request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := v.Validate(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Failed to validate request body")
		code, errors := errmsg.Errors(ctx, err, req)
		return c.Status(code).JSON(response.Error(errors))
	}

	item, err := h.core.CreateUploadURL(ctx, coreentity.PresignUploadURLReq{
		UserCtx:          uc,
		TenantID:         uc.TenantID,
		CreatedBy:        uc.UserID,
		Filename:         req.Filename,
		OriginalFilename: req.OriginalFilename,
		ContentType:      req.ContentType,
		Folder:           req.Folder,
		IsPublic:         req.IsPublic,
	})
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to create upload URL")
		code, errors := errmsg.Errors[error](ctx, err)
		return c.Status(code).JSON(response.Error(errors))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(restentity.CreateUploadURLResp{
		FileID:    item.FileID,
		ObjectKey: item.Filename,
		UploadURL: item.URL,
		Method:    item.Method,
		Headers:   item.Headers,
	}, ""))
}
