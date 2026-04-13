package handler

import (
	portsCore "codebase-app/internal/ports/core"

	"github.com/gofiber/fiber/v2"
)

type meHandler struct {
	core portsCore.MeCore
}

type MeHandlerConfig struct {
	Core portsCore.MeCore
}

func NewMeHandler(cfg MeHandlerConfig) *meHandler {
	return &meHandler{
		core: cfg.Core,
	}
}

func (h *meHandler) Register(router fiber.Router) {
	router.Get("/", h.getMe)
	router.Get("/employee", h.getMyEmployee)
	router.Get("/shift-today", h.getMyShiftToday)
}
