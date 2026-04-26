package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
)

type InternalQuotationCore interface {
	GetQuotations(ctx context.Context, filter coreentity.InternalCommerceListFilter) ([]coreentity.InternalQuotation, int, error)
	CreateQuotation(ctx context.Context, input coreentity.CreateInternalQuotationInput) (*coreentity.InternalQuotation, error)
	GetQuotation(ctx context.Context, id string) (*coreentity.InternalQuotation, error)
	ExecuteQuotationAction(ctx context.Context, input coreentity.InternalQuotationActionInput) (*coreentity.InternalCommerceBundle, error)
}

type InternalOrderCore interface {
	GetOrders(ctx context.Context, filter coreentity.InternalCommerceListFilter) ([]coreentity.InternalOrder, int, error)
	CreateOrder(ctx context.Context, input coreentity.CreateInternalOrderInput) (*coreentity.InternalCommerceBundle, error)
	GetOrder(ctx context.Context, id string) (*coreentity.InternalCommerceBundle, error)
}

type InternalInvoiceCore interface {
	GetInvoices(ctx context.Context, filter coreentity.InternalCommerceListFilter) ([]coreentity.InternalInvoice, int, error)
	GetInvoice(ctx context.Context, id string) (*coreentity.InternalInvoiceDetail, error)
	ExecuteInvoiceAction(ctx context.Context, input coreentity.InternalInvoiceActionInput) (*coreentity.InternalPaymentAttemptResult, error)
	GetPaymentAttempts(ctx context.Context, invoiceID string) ([]coreentity.InternalPaymentAttempt, error)
	ExecutePaymentAttemptAction(ctx context.Context, input coreentity.InternalPaymentAttemptActionInput) (*coreentity.InternalPaymentAttemptResult, error)
	CreatePaymentReceipt(ctx context.Context, receipt coreentity.InternalPaymentReceipt) (*coreentity.InternalPaymentReceiptResult, error)
}
