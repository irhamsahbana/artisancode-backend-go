package handler

import (
	corePorts "codebase-app/internal/ports/core"

	"github.com/gofiber/fiber/v2"
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
	router.Get("/entitlements", h.getEntitlements)
	router.Get("/invoices", h.getInvoices)
	router.Post("/checkouts", h.createCheckout)
}
