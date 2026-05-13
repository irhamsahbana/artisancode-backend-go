package handler

import (
	portsCore "codebase-app/internal/ports/core"

	"github.com/gofiber/fiber/v3"
)

type workLocationHandler struct {
	core portsCore.WorkLocationCore
}

type Config struct {
	Core portsCore.WorkLocationCore
}

func NewWorkLocationHandler(cfg Config) *workLocationHandler {
	return &workLocationHandler{core: cfg.Core}
}

func (h *workLocationHandler) Register(router fiber.Router) {
	router.Get("/", h.getWorkLocations)
	router.Get("/:id", h.getWorkLocation)
	router.Post("/", h.createWorkLocation)
	router.Put("/:id", h.updateWorkLocation)
	router.Delete("/:id", h.deleteWorkLocation)
}
