package handler

import (
	corePorts "codebase-app/internal/ports/core"

	"github.com/gofiber/fiber/v3"
)

type webhookHandler struct {
	core corePorts.WebhookCore
}

type Config struct {
	Core corePorts.WebhookCore
}

func NewWebhookHandler(cfg Config) *webhookHandler {
	return &webhookHandler{core: cfg.Core}
}

func (h *webhookHandler) Register(router fiber.Router) {
	router.Post("/doku", h.handleDOKUWebhook)
}
