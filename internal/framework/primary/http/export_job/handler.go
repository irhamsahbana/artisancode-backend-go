package handler

import (
	portsCore "codebase-app/internal/ports/core"

	"github.com/gofiber/fiber/v3"
)

type exportJobHandler struct {
	core portsCore.ExportJobCore
}

type Config struct {
	Core portsCore.ExportJobCore
}

func NewExportJobHandler(cfg Config) *exportJobHandler {
	return &exportJobHandler{core: cfg.Core}
}

func (h *exportJobHandler) Register(router fiber.Router) {
	router.Get("/", h.getExportJobs)
	router.Get("/:id", h.getExportJob)
	router.Post("/", h.createExportJob)
}
