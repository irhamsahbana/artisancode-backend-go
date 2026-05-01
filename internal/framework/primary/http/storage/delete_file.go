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

func (h *storageHandler) deleteFile(c *fiber.Ctx) error {
	tracedCtx, span := tracing.StartSpan(c.UserContext(), "internal:framework:primary:http:storage:delete_file:deleteFile")
	defer span.End()
	c.SetUserContext(tracedCtx)

	var (
		ctx      = c.UserContext()
		filename = c.Params("*")
		uc       = common.GetUserContext(ctx)
	)

	req := &coreentity.DeleteFileReq{
		TenantID: uc.TenantID,
		Filename: filename,
	}

	if err := h.core.DeleteFile(ctx, req); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to delete file")
		code, errs := errmsg.Errors[error](ctx, err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(nil, ""))
}
