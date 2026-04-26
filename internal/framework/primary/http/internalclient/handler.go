package handler

import (
	corePorts "codebase-app/internal/ports/core"

	"github.com/gofiber/fiber/v2"
)

type internalClientHandler struct {
	core corePorts.InternalClientCore
}

type Config struct {
	Core corePorts.InternalClientCore
}

func NewInternalClientHandler(cfg Config) *internalClientHandler {
	return &internalClientHandler{core: cfg.Core}
}

func (h *internalClientHandler) Register(router fiber.Router) {
	router.Get("/", h.getInternalClients)
	router.Get("/:id/owner-permissions", h.getInternalClientOwnerPermissions)
	router.Put("/:id/owner-permissions", h.updateInternalClientOwnerPermissions)
}
