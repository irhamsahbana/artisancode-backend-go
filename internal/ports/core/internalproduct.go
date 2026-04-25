package core

import (
	"context"

	"codebase-app/internal/entity/coreentity"
)

type InternalProductCore interface {
	GetInternalProducts(ctx context.Context, filter coreentity.InternalProductListFilter) ([]coreentity.InternalProduct, int, error)
	GetInternalProduct(ctx context.Context, filter coreentity.InternalProductFilter) (*coreentity.InternalProduct, error)
	CreateInternalProduct(ctx context.Context, data coreentity.InternalProduct) (*coreentity.InternalProduct, error)
	UpdateInternalProduct(ctx context.Context, data coreentity.InternalProduct) error
	DeleteInternalProduct(ctx context.Context, filter coreentity.InternalProductDeleteFilter) error

	GetInternalProductPricings(ctx context.Context, filter coreentity.InternalProductPricingListFilter) ([]coreentity.InternalProductPricing, int, error)
	GetInternalProductPricing(ctx context.Context, filter coreentity.InternalProductPricingFilter) (*coreentity.InternalProductPricing, error)
	CreateInternalProductPricing(ctx context.Context, data coreentity.InternalProductPricing) (*coreentity.InternalProductPricing, error)
	UpdateInternalProductPricing(ctx context.Context, data coreentity.InternalProductPricing) error
	DeleteInternalProductPricing(ctx context.Context, filter coreentity.InternalProductPricingDeleteFilter) error

	GetInternalProductPrices(ctx context.Context, filter coreentity.InternalProductPriceListFilter) ([]coreentity.InternalProductPrice, int, error)
	CreateInternalProductPrice(ctx context.Context, data coreentity.InternalProductPrice) (*coreentity.InternalProductPrice, error)
	UpdateInternalProductPrice(ctx context.Context, data coreentity.InternalProductPrice) error
	DeleteInternalProductPrice(ctx context.Context, filter coreentity.InternalProductPriceDeleteFilter) error
}
