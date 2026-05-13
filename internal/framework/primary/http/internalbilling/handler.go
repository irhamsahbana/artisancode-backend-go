package handler

import (
	"codebase-app/internal/ports/core"

	"github.com/gofiber/fiber/v3"
)

type internalBillingHandler struct {
	core core.InternalTenantBillingCore
}

func NewInternalBillingHandler(core core.InternalTenantBillingCore) *internalBillingHandler {
	return &internalBillingHandler{core: core}
}

func (h *internalBillingHandler) Register(router fiber.Router) {
	group := router.Group("/internal-billing")
	group.Get("/invoices", h.getAllInvoices)
	group.Get("/invoices/:id", h.getInvoiceByID)
	group.Get("/ledger/:tenantId", h.getLedgerEntries)
	group.Get("/reconciliation-cases", h.getReconciliationCases)
	group.Post("/manual-invoices", h.createManualInvoice)
	group.Post("/payment-receipts", h.createPaymentReceipt)
	group.Post("/payment-receipts/:id/actions", h.executePaymentReceiptAction)
}
