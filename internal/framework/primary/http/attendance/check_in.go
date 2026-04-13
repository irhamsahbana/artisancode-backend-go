package handler

import (
	"codebase-app/internal/adapter"
	"codebase-app/internal/entity/mapper"
	"codebase-app/internal/entity/restentity"
	"codebase-app/internal/infrastructure/tracing"
	"codebase-app/pkg/errmsg"
	"codebase-app/pkg/response"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (h *attendanceHandler) checkIn(c *fiber.Ctx) error {
	tracedCtx, span := tracing.StartSpan(c.UserContext(), "internal:framework:primary:http:attendance:check_in:checkIn")
	defer span.End()
	c.SetUserContext(tracedCtx)

	return h.handleAttendanceAction(c, "check_in")
}

func (h *attendanceHandler) handleAttendanceAction(c *fiber.Ctx, actionType string) error {
	tracedCtx, span := tracing.StartSpan(c.UserContext(), "internal:framework:primary:http:attendance:check_in:handleAttendanceAction")
	defer span.End()
	c.SetUserContext(tracedCtx)

	var (
		ctx = c.UserContext()
		req = new(restentity.CheckAttendanceReq)
		v   = adapter.Adapters.Validator
	)

	if err := c.BodyParser(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Invalid request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error(err))
	}

	if err := v.Validate(req); err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("Failed to validate request body")
		code, errors := errmsg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errors))
	}

	item, err := map[string]func() (*restentity.CheckAttendanceResp, error){
		"check_in": func() (*restentity.CheckAttendanceResp, error) {
			resp, err := h.core.CheckIn(ctx, mapper.AttendanceActionFromRest(ctx, *req, actionType))
			if err != nil {
				return nil, err
			}
			return &restentity.CheckAttendanceResp{ID: resp.ID}, nil
		},
		"check_out": func() (*restentity.CheckAttendanceResp, error) {
			resp, err := h.core.CheckOut(ctx, mapper.AttendanceActionFromRest(ctx, *req, actionType))
			if err != nil {
				return nil, err
			}
			return &restentity.CheckAttendanceResp{ID: resp.ID}, nil
		},
	}[actionType]()
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Failed to create attendance log")
		code, errors := errmsg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errors))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(item, "Attendance recorded successfully"))
}
