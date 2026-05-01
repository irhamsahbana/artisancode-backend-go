package handler

import (
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (h *storageHandler) listFiles(c *fiber.Ctx) error {
	tracedCtx, span := tracing.StartSpan(c.UserContext(), "internal:framework:primary:http:storage:list_files:listFiles")
	defer span.End()
	c.SetUserContext(tracedCtx)

	ctx := c.UserContext()
	resp, err := h.core.ListFiles(ctx)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to list files")
		code, errs := errmsg.Errors[error](ctx, err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(resp, ""))
}
