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
	GetInvoice(
		ctx context.Context,
		tenantID string,
		id string,
	) (*coreentity.TenantBillingInvoiceDetail, error)
	GetPaymentAttemptsByInvoice(
		ctx context.Context,
		tenantID string,
		invoiceID string,
	) ([]coreentity.TenantBillingPaymentAttemptView, error)
	GetInvoiceRaw(
		ctx context.Context,
		tenantID string,
		id string,
	) (*coreentity.InternalTenantInvoice, error)
	GetPaymentAttemptByID(
		ctx context.Context,
		tenantID string,
		id string,
	) (*coreentity.InternalTenantPaymentAttempt, error)
	CancelInvoice(
		ctx context.Context,
		tenantID string,
		id string,
	) error
	CancelPaymentAttempt(
		ctx context.Context,
		tenantID string,
		id string,
	) error
	SetSubscriptionStatus(
		ctx context.Context,
		tenantID string,
		subscriptionID string,
		status string,
	) error
	GetAddOnsByIDs(
		ctx context.Context,
		addOnIDs []string,
	) ([]coreentity.TenantBillingAddOn, error)

	GetAllInvoices(ctx context.Context, filter coreentity.InternalBillingInvoiceListFilter) ([]coreentity.InternalTenantInvoice, error)
	GetInvoiceByID(ctx context.Context, id string) (*coreentity.InternalTenantInvoice, error)
	GetLedgerEntries(ctx context.Context, tenantID string, filter coreentity.InternalBillingLedgerListFilter) ([]coreentity.InternalBillingLedgerEntry, error)
	GetReconciliationCases(ctx context.Context, filter coreentity.InternalBillingReconciliationCaseFilter) ([]coreentity.InternalBillingReconciliationCase, error)
	CreateInternalInvoice(ctx context.Context, data coreentity.InternalTenantInvoice) (*coreentity.InternalTenantInvoice, error)
	CreatePaymentReceipt(ctx context.Context, data coreentity.InternalPaymentReceipt) (*coreentity.InternalPaymentReceipt, error)
	GetPaymentReceipt(ctx context.Context, id string) (*coreentity.InternalPaymentReceipt, error)
	UpdatePaymentReceipt(ctx context.Context, data coreentity.InternalPaymentReceipt) (*coreentity.InternalPaymentReceipt, error)

	GetSubscriptionsPastPeriodEnd(ctx context.Context, limit int) ([]coreentity.InternalTenantSubscription, error)
	GetSubscriptionsInGracePastDue(ctx context.Context, limit int) ([]coreentity.InternalTenantSubscription, error)
	GetPricingInfo(ctx context.Context, pricingID string) (*coreentity.PricingInfo, error)
}
