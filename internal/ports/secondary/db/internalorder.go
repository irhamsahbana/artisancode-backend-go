package repository

import (
	"context"

	"codebase-app/internal/entity/coreentity"
)

type InternalOrderRepository interface {
	GetOrders(ctx context.Context, filter coreentity.InternalCommerceListFilter) ([]coreentity.InternalOrder, int, error)
	CreateOrder(ctx context.Context, data coreentity.InternalOrder) (*coreentity.InternalOrder, error)
	GetOrder(ctx context.Context, tenantID, id string) (*coreentity.InternalOrder, error)
	GetOrderByInvoiceID(ctx context.Context, tenantID, invoiceID string) (*coreentity.InternalOrder, error)
	MarkOrderPaid(ctx context.Context, tenantID, id string) (*coreentity.InternalOrder, error)
	GetPricingSnapshot(ctx context.Context, productID, pricingID, currencyCode string) (map[string]any, string, error)
}
