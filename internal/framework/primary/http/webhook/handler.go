package handler

import (
	corePorts "codebase-app/internal/ports/core"

	"github.com/gofiber/fiber/v2"
)

type webhookHandler struct {
	core corePorts.WebhookCore
}

type WebhookHandlerConfig struct {
	Core corePorts.WebhookCore
}

func NewWebhookHandler(cfg WebhookHandlerConfig) *webhookHandler {
	return &webhookHandler{core: cfg.Core}
}

func (h *webhookHandler) Register(router fiber.Router) {
	router.Post("/doku", h.handleDOKUWebhook)
}
