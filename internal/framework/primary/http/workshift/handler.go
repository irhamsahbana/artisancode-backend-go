package handler

import (
	portsCore "codebase-app/internal/ports/core"

	"github.com/gofiber/fiber/v2"
)

type workShiftHandler struct {
	core portsCore.WorkShiftCore
}

type Config struct {
	Core portsCore.WorkShiftCore
}

func NewWorkShiftHandler(cfg Config) *workShiftHandler {
	return &workShiftHandler{core: cfg.Core}
}

func (h *workShiftHandler) Register(router fiber.Router) {
	router.Get("/", h.getWorkShifts)
	router.Get("/:id", h.getWorkShift)
	router.Post("/", h.createWorkShift)
	router.Put("/:id", h.updateWorkShift)
	router.Delete("/:id", h.deleteWorkShift)
}
