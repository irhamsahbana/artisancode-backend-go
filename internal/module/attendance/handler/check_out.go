package handler

import "github.com/gofiber/fiber/v2"

func (h *attendanceHandler) checkOut(c *fiber.Ctx) error {
	return h.handleAttendanceAction(c, "check_out")
}
