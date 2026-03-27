package handler

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/entity/mapper"
	"codebase-app/internal/entity/restentity"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (h *workShiftHandler) createWorkShift(c *fiber.Ctx) error {
	var (
		ctx = c.UserContext()
		req = new(restentity.CreateWorkShiftReq)
		v   = adapter.Adapters.Validator
	)

	if err := c.BodyParser(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := v.Validate(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Failed to validate request body")
		code, errors := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errors))
	}

	data := mapper.WorkShiftFromRestCreateToCore(ctx, *req)
	created, err := h.core.CreateWorkShift(ctx, data)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to create work shift")
		code, errors := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errors))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(restentity.CreateWorkShiftResp{ID: created.ID}, ""))
}