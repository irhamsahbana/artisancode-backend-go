package handler

import (
	orgUnitPorts "codebase-app/internal/ports/core"

	"github.com/gofiber/fiber/v2"
)

type orgUnitHandler struct {
	core orgUnitPorts.OrgUnitCore
}

type OrgUnitHandlerConfig struct {
	Core orgUnitPorts.OrgUnitCore
}

func NewOrgUnitHandler(cfg OrgUnitHandlerConfig) *orgUnitHandler {
	return &orgUnitHandler{core: cfg.Core}
}

func (h *orgUnitHandler) Register(router fiber.Router) {
	router.Get("/", h.getOrgUnits)
	router.Get("/tree/:companyId", h.getOrgUnitTree)
	router.Get("/:id", h.getOrgUnit)
	router.Post("/", h.createOrgUnit)
	router.Put("/:id", h.updateOrgUnit)
	router.Delete("/:id", h.deleteOrgUnit)
}
