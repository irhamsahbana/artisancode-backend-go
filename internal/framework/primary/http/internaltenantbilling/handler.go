package handler

import (
	corePorts "codebase-app/internal/ports/core"

	"github.com/gofiber/fiber/v3"
)

type internalTenantBillingHandler struct {
	core corePorts.InternalTenantBillingCore
}

type Config struct {
	Core corePorts.InternalTenantBillingCore
}

func NewInternalTenantBillingHandler(cfg Config) *internalTenantBillingHandler {
	return &internalTenantBillingHandler{core: cfg.Core}
}

func (h *internalTenantBillingHandler) Register(router fiber.Router) {
	router.Get("/plans", h.getPlans)
	router.Get("/subscription", h.getSubscription)
	router.Post("/subscription/actions", h.executeSubscriptionAction)
	router.Get("/entitlements", h.getEntitlements)
	router.Get("/invoices/:id", h.getInvoice)
	router.Get("/invoices/:id/payment-attempts", h.getPaymentAttempts)
	router.Post("/invoices/:id/actions", h.executeInvoiceAction)
	router.Get("/invoices", h.getInvoices)
	router.Post("/checkouts", h.createCheckout)
	router.Post("/payment-attempts/:id/actions", h.executePaymentAttemptAction)
	router.Post("/add-ons/actions", h.executeAddOnsAction)
}
