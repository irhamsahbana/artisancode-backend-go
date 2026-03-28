package handler

import (
	portsCore "codebase-app/internal/ports/core"

	"github.com/gofiber/fiber/v2"
)

type attendanceHandler struct {
	core portsCore.AttendanceCore
}

type AttendanceHandlerConfig struct {
	Core portsCore.AttendanceCore
}

func NewAttendanceHandler(cfg AttendanceHandlerConfig) *attendanceHandler {
	return &attendanceHandler{core: cfg.Core}
}

func (h *attendanceHandler) Register(router fiber.Router) {
	router.Get("/", h.getAttendanceLogs)
	router.Get("/:id", h.getAttendanceLog)
	router.Post("/check-in", h.checkIn)
	router.Post("/check-out", h.checkOut)
}

func (h *attendanceHandler) RegisterSummary(router fiber.Router) {
	router.Get("/today", h.getAttendanceSummaryToday)
}

func (h *attendanceHandler) RegisterPolicy(router fiber.Router) {
	router.Get("/", h.getAttendancePolicy)
}
