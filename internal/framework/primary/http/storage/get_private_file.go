package handler

import (
	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (h *storageHandler) getPrivateFile(c *fiber.Ctx) error {
	tracedCtx, span := tracing.StartSpan(c.UserContext(), "internal:framework:primary:http:storage:get_private_file:getPrivateFile")
	defer span.End()
	c.SetUserContext(tracedCtx)

	var (
		ctx      = c.UserContext()
		filename = c.Params("*")
		tenantID = c.Query("tenant_id")
		folder   = c.Query("folder")
	)

	filename = "private/" + filename

	filter := coreentity.FileFilter{
		TenantID: tenantID,
		Folder:   common.S3Folder(folder),
		Filename: filename,
	}

	url, err := h.core.GetFileURL(ctx, filter)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Any("filename", filename).Msg("Failed to get file URL")
		code, errs := errmsg.Errors[error](ctx, err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(url, ""))
}
