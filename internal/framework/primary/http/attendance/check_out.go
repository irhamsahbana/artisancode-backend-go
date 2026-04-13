package handler

import (
	"codebase-app/internal/infrastructure/tracing"
	"github.com/gofiber/fiber/v2"
)

func (h *attendanceHandler) checkOut(c *fiber.Ctx) error {
	tracedCtx, span := tracing.StartSpan(c.UserContext(), "internal:framework:primary:http:attendance:check_out:checkOut")
	defer span.End()
	c.SetUserContext(tracedCtx)

	return h.handleAttendanceAction(c, "check_out")
}
