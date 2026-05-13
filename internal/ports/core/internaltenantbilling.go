package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
)

type InternalTenantBillingCore interface {
	GetPlans(ctx context.Context) ([]coreentity.TenantBillingPlan, error)
	GetSubscription(ctx context.Context) (*coreentity.TenantBillingSubscriptionView, error)
	GetEntitlements(ctx context.Context) (*coreentity.InternalEntitlementSnapshot, error)
	GetInvoices(
		ctx context.Context,
		filter coreentity.TenantBillingInvoiceListFilter,
	) ([]coreentity.TenantBillingInvoice, error)
	CreateCheckout(
		ctx context.Context,
		input coreentity.TenantBillingCheckoutInput,
	) (*coreentity.TenantBillingCheckoutResult, error)
	GetInvoice(
		ctx context.Context,
		id string,
	) (*coreentity.TenantBillingInvoiceDetail, error)
	GetPaymentAttempts(
		ctx context.Context,
		invoiceID string,
	) ([]coreentity.TenantBillingPaymentAttemptView, error)
	ExecuteInvoiceAction(
		ctx context.Context,
		input coreentity.TenantBillingInvoiceActionInput,
	) (*coreentity.TenantBillingInvoiceActionResult, error)
	ExecutePaymentAttemptAction(
		ctx context.Context,
		input coreentity.TenantBillingPaymentAttemptActionInput,
	) (*coreentity.TenantBillingPaymentAttemptActionResult, error)
	ExecuteSubscriptionAction(
		ctx context.Context,
		input coreentity.TenantBillingSubscriptionActionInput,
	) (*coreentity.TenantBillingSubscriptionActionResult, error)
	ExecuteAddOnsAction(
		ctx context.Context,
		input coreentity.TenantBillingAddOnsActionInput,
	) (*coreentity.TenantBillingAddOnsActionResult, error)

	GetAllInvoices(
		ctx context.Context,
		filter coreentity.InternalBillingInvoiceListFilter,
	) ([]coreentity.InternalTenantInvoice, error)
	GetInvoiceByID(
		ctx context.Context,
		id string,
	) (*coreentity.InternalTenantInvoice, error)
	GetLedgerEntries(
		ctx context.Context,
		tenantID string,
		filter coreentity.InternalBillingLedgerListFilter,
	) ([]coreentity.InternalBillingLedgerEntry, error)
	GetReconciliationCases(
		ctx context.Context,
		filter coreentity.InternalBillingReconciliationCaseFilter,
	) ([]coreentity.InternalBillingReconciliationCase, error)
	CreateManualInvoice(
		ctx context.Context,
		input coreentity.InternalBillingManualInvoiceInput,
	) (*coreentity.InternalBillingManualInvoiceResult, error)
	CreatePaymentReceipt(
		ctx context.Context,
		input coreentity.InternalPaymentReceipt,
	) (*coreentity.InternalBillingPaymentReceiptActionResult, error)
	ExecutePaymentReceiptAction(
		ctx context.Context,
		input coreentity.InternalBillingPaymentReceiptActionInput,
	) (*coreentity.InternalBillingPaymentReceiptActionResult, error)

	CheckUsageLimit(
		ctx context.Context,
		resourceType string,
		currentCount int64,
	) error
}
