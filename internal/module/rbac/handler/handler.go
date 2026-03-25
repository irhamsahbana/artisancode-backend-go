package handler

import (
	corePorts "codebase-app/internal/ports/core"

	"github.com/gofiber/fiber/v2"
)

type rbacHandler struct {
	core corePorts.RbacCore
}

func NewRbacHandler(core corePorts.RbacCore) *rbacHandler {
	return &rbacHandler{
		core: core,
	}
}

func (h *rbacHandler) Register(router fiber.Router) {
}
