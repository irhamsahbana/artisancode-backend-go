package repository

import (
	"context"

	"codebase-app/internal/entity/coreentity"
)

type InternalTenantBillingRepository interface {
	EnsureBillingAccount(ctx context.Context, tenantID string) (*coreentity.InternalBillingAccount, error)
	GetInvoices(
		ctx context.Context,
		filter coreentity.TenantBillingInvoiceListFilter,
	) ([]coreentity.TenantBillingInvoice, error)
	CreateInvoice(ctx context.Context, data coreentity.InternalTenantInvoice) (*coreentity.InternalTenantInvoice, error)
	CreatePaymentAttempt(
		ctx context.Context,
		data coreentity.InternalTenantPaymentAttempt,
	) (*coreentity.InternalTenantPaymentAttempt, error)
	UpdatePaymentAttemptGateway(
		ctx context.Context,
		data coreentity.InternalTenantPaymentAttempt,
	) (*coreentity.InternalTenantPaymentAttempt, error)
}
