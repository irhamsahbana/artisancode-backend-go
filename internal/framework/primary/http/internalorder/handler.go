package handler

import (
	portsCore "codebase-app/internal/ports/core"

	"github.com/gofiber/fiber/v3"
)

type internalOrderHandler struct {
	core portsCore.InternalOrderCore
}

type Config struct {
	Core portsCore.InternalOrderCore
}

func NewInternalOrderHandler(cfg Config) *internalOrderHandler {
	return &internalOrderHandler{core: cfg.Core}
}

func (h *internalOrderHandler) Register(router fiber.Router) {
	router.Get("/orders", h.getOrders)
	router.Post("/orders", h.createOrder)
	router.Get("/orders/:id", h.getOrder)
}
