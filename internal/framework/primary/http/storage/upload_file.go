package handler

import (
	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"
	"errors"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
)

func (h *storageHandler) uploadFile(c fiber.Ctx) error {
	tracedCtx, span := tracing.StartSpan(c.Context(), "internal:framework:primary:http:storage:upload_file:uploadFile")
	defer span.End()
	c.SetContext(tracedCtx)

	ctx := c.Context()

	file, err := c.FormFile("file")
	if err != nil {
		if errors.Is(err, fasthttp.ErrMissingFile) {
			message := errmsg.NewCustomErrors(http.StatusBadRequest).
				Add("file", errmsg.MessageFileIsRequired).
				SetMessage(errmsg.MessageFileIsRequired)
			return c.Status(message.Code).JSON(response.Error(message))
		}
		log.Ctx(ctx).Error().Err(err).Msg("Failed to parse form file")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	req := &coreentity.UploadFileReq{
		File: file,
	}
	req.Filename = c.FormValue("filename")
	req.IsPublic = parseBool(c.FormValue("is_public"))
	req.GeneratePresignedURL = parseBool(c.FormValue("generate_presigned_url"))
	if folder := c.FormValue("folder"); folder != "" {
		req.Folder = common.S3Folder(folder)
	}

	resp, err := h.core.UploadFile(ctx, req)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to upload file")
		code, errs := errmsg.Errors[error](ctx, err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(resp, ""))
}

func parseBool(value string) bool {
	if value == "" {
		return false
	}

	result, err := strconv.ParseBool(value)
	if err != nil {
		return false
	}

	return result
}
