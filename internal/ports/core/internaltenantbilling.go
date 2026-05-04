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
}
