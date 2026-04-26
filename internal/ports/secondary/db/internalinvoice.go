package repository

import (
	"context"

	"codebase-app/internal/entity/coreentity"
)

type InternalInvoiceRepository interface {
	GetInvoices(ctx context.Context, filter coreentity.InternalCommerceListFilter) ([]coreentity.InternalInvoice, int, error)
	CreateInvoice(ctx context.Context, data coreentity.InternalInvoice) (*coreentity.InternalInvoice, error)
	GetInvoice(ctx context.Context, tenantID, id string) (*coreentity.InternalInvoice, error)
	GetInvoiceByOrderID(ctx context.Context, tenantID, orderID string) (*coreentity.InternalInvoice, error)
	MarkInvoicePaid(ctx context.Context, tenantID, id string, amountPaid string) (*coreentity.InternalInvoice, error)
	CreatePaymentAttempt(ctx context.Context, data coreentity.InternalPaymentAttempt) (*coreentity.InternalPaymentAttempt, error)
	UpdatePaymentAttemptGateway(ctx context.Context, data coreentity.InternalPaymentAttempt) (*coreentity.InternalPaymentAttempt, error)
	GetPaymentAttempt(ctx context.Context, tenantID, id string) (*coreentity.InternalPaymentAttempt, error)
	GetLatestPaymentAttempt(ctx context.Context, tenantID, invoiceID string) (*coreentity.InternalPaymentAttempt, error)
	GetPaymentAttempts(ctx context.Context, tenantID, invoiceID string) ([]coreentity.InternalPaymentAttempt, error)
	CreatePaymentReceipt(ctx context.Context, data coreentity.InternalPaymentReceipt) (*coreentity.InternalPaymentReceipt, error)
}
