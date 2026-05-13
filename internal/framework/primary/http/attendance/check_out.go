package handler

import (
	"codebase-app/internal/infrastructure/tracing"
	"github.com/gofiber/fiber/v3"
)

func (h *attendanceHandler) checkOut(c fiber.Ctx) error {
	tracedCtx, span := tracing.StartSpan(c.Context(), "internal:framework:primary:http:attendance:check_out:checkOut")
	defer span.End()
	c.SetContext(tracedCtx)

	return h.handleAttendanceAction(c, "check_out")
}
