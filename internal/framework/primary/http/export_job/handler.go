package handler

import (
	portsCore "codebase-app/internal/ports/core"

	"github.com/gofiber/fiber/v2"
)

type exportJobHandler struct {
	core portsCore.ExportJobCore
}

type ExportJobHandlerConfig struct {
	Core portsCore.ExportJobCore
}

func NewExportJobHandler(cfg ExportJobHandlerConfig) *exportJobHandler {
	return &exportJobHandler{core: cfg.Core}
}

func (h *exportJobHandler) Register(router fiber.Router) {
	router.Get("/", h.getExportJobs)
	router.Get("/:id", h.getExportJob)
	router.Post("/", h.createExportJob)
}
