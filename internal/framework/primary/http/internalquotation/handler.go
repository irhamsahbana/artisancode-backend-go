package handler

import (
	portsCore "codebase-app/internal/ports/core"

	"github.com/gofiber/fiber/v3"
)

type internalQuotationHandler struct {
	core portsCore.InternalQuotationCore
}

type Config struct {
	Core portsCore.InternalQuotationCore
}

func NewInternalQuotationHandler(cfg Config) *internalQuotationHandler {
	return &internalQuotationHandler{core: cfg.Core}
}

func (h *internalQuotationHandler) Register(router fiber.Router) {
	router.Get("/quotations", h.getQuotations)
	router.Post("/quotations", h.createQuotation)
	router.Get("/quotations/:id", h.getQuotation)
	router.Post("/quotations/:id/actions", h.executeAction)
}
