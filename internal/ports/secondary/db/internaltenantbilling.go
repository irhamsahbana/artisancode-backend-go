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

	GetSubscription(ctx context.Context, tenantID string) (*coreentity.InternalTenantSubscription, error)
	GetLatestEntitlementSnapshot(ctx context.Context, tenantID string) (*coreentity.InternalEntitlementSnapshot, error)
	GetActiveSubscription(ctx context.Context, tenantID string) (*coreentity.InternalTenantSubscription, error)
	UpsertSubscription(ctx context.Context, data coreentity.InternalTenantSubscription) (*coreentity.InternalTenantSubscription, error)
	CreateSubscriptionChange(ctx context.Context, data coreentity.InternalTenantSubscriptionChange) error
	CreateEntitlementSnapshot(ctx context.Context, data coreentity.InternalEntitlementSnapshot) error
	CreateLedgerEntry(ctx context.Context, data coreentity.InternalBillingLedgerEntry) error
	UpdatePaymentAttemptStatus(ctx context.Context, id string, status string, metadata map[string]any) error
	UpdateInvoiceStatus(ctx context.Context, id string, status string) error
	GetInvoiceByNumber(ctx context.Context, number string) (*coreentity.InternalTenantInvoice, error)
	GetPaymentAttemptByProviderRef(ctx context.Context, providerReference string) (*coreentity.InternalTenantPaymentAttempt, error)
}
