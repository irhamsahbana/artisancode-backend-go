package handler

import (
	portsCore "codebase-app/internal/ports/core"

	"github.com/gofiber/fiber/v2"
)

type jobPositionHandler struct {
	core portsCore.JobPositionCore
}

type Config struct {
	Core portsCore.JobPositionCore
}

func NewJobPositionHandler(cfg Config) *jobPositionHandler {
	return &jobPositionHandler{core: cfg.Core}
}

func (h *jobPositionHandler) Register(router fiber.Router) {
	router.Get("/", h.getJobPositions)
	router.Get("/:id", h.getJobPosition)
	router.Post("/", h.createJobPosition)
	router.Put("/:id", h.updateJobPosition)
	router.Delete("/:id", h.deleteJobPosition)
}
