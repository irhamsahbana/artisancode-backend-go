package repository

import (
	"context"

	"codebase-app/internal/entity/coreentity"
)

type InternalProductRepository interface {
	GetInternalProducts(ctx context.Context, filter coreentity.InternalProductListFilter) ([]coreentity.InternalProduct, int, error)
	GetInternalProduct(ctx context.Context, filter coreentity.InternalProductFilter) (*coreentity.InternalProduct, error)
	CreateInternalProduct(ctx context.Context, data coreentity.InternalProduct) (*coreentity.InternalProduct, error)
	UpdateInternalProduct(ctx context.Context, data coreentity.InternalProduct) error
	DeleteInternalProduct(ctx context.Context, filter coreentity.InternalProductDeleteFilter) error
	ExistsInternalProductByCode(ctx context.Context, code, excludeID string) (bool, error)

	GetInternalProductPricings(ctx context.Context, filter coreentity.InternalProductPricingListFilter) ([]coreentity.InternalProductPricing, int, error)
	GetInternalProductPricing(ctx context.Context, filter coreentity.InternalProductPricingFilter) (*coreentity.InternalProductPricing, error)
	CreateInternalProductPricing(ctx context.Context, data coreentity.InternalProductPricing) (*coreentity.InternalProductPricing, error)
	UpdateInternalProductPricing(ctx context.Context, data coreentity.InternalProductPricing) error
	DeleteInternalProductPricing(ctx context.Context, filter coreentity.InternalProductPricingDeleteFilter) error
	ExistsInternalProductPricingByCode(ctx context.Context, internalProductID, code, excludeID string) (bool, error)

	GetInternalProductPrices(ctx context.Context, filter coreentity.InternalProductPriceListFilter) ([]coreentity.InternalProductPrice, int, error)
	GetInternalProductPrice(ctx context.Context, filter coreentity.InternalProductPriceFilter) (*coreentity.InternalProductPrice, error)
	CreateInternalProductPrice(ctx context.Context, data coreentity.InternalProductPrice) (*coreentity.InternalProductPrice, error)
	UpdateInternalProductPrice(ctx context.Context, data coreentity.InternalProductPrice) error
	DeleteInternalProductPrice(ctx context.Context, filter coreentity.InternalProductPriceDeleteFilter) error
	ExistsOverlappingInternalProductPrice(ctx context.Context, filter coreentity.InternalProductPriceOverlapFilter) (bool, error)
}
