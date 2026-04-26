package handler

import (
	portsCore "codebase-app/internal/ports/core"

	"github.com/gofiber/fiber/v2"
)

type rbacHandler struct {
	core portsCore.RbacCore
}

// RbacHandlerConfig is the config struct for the RBAC handler
type Config struct {
	Core portsCore.RbacCore
}

func NewRbacHandler(cfg Config) *rbacHandler {
	return &rbacHandler{core: cfg.Core}
}

func (h *rbacHandler) Register(router fiber.Router) {
	router.Get("/roles", h.getRoles)
	router.Get("/roles/:id", h.getRole)
	router.Post("/roles", h.createRole)
	router.Put("/roles/:id", h.updateRole)
	router.Delete("/roles/:id", h.deleteRole)
	router.Get("/permissions", h.getPermissions)
}
