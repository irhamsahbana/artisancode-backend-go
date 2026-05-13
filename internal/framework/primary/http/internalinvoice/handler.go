package handler

import (
	portsCore "codebase-app/internal/ports/core"

	"github.com/gofiber/fiber/v3"
)

type internalInvoiceHandler struct {
	core portsCore.InternalInvoiceCore
}

type Config struct {
	Core portsCore.InternalInvoiceCore
}

func NewInternalInvoiceHandler(cfg Config) *internalInvoiceHandler {
	return &internalInvoiceHandler{core: cfg.Core}
}

func (h *internalInvoiceHandler) Register(router fiber.Router) {
	router.Get("/invoices", h.getInvoices)
	router.Get("/invoices/:id", h.getInvoice)
	router.Post("/invoices/:id/actions", h.executeInvoiceAction)
	router.Get("/invoices/:id/payment-attempts", h.getPaymentAttempts)
	router.Post("/payment-attempts/:id/actions", h.executePaymentAttemptAction)
	router.Post("/payment-receipts", h.createPaymentReceipt)
}
