package coreentity

import "codebase-app/internal/entity/common"

const (
	InternalProductStatusDraft    = "draft"
	InternalProductStatusActive   = "active"
	InternalProductStatusInactive = "inactive"
	InternalProductStatusArchived = "archived"
)

type InternalProduct struct {
	UserCtx common.UserContext

	ID          string
	Code        string
	Name        string
	Description string
	Status      string
	Metadata    map[string]any
	CreatedAt   string
	UpdatedAt   string
}

type InternalProductListFilter struct {
	Q        string
	Page     int
	Paginate int
}

type InternalProductFilter struct {
	ID string
}

type InternalProductDeleteFilter struct {
	ID string
}

type InternalProductPricing struct {
	UserCtx common.UserContext

	ID                string
	InternalProductID string
	Code              string
	Name              string
	Description       string
	Status            string
	Metadata          map[string]any
	CreatedAt         string
	UpdatedAt         string
}

type InternalProductPricingListFilter struct {
	InternalProductID string
	Q                 string
	Page              int
	Paginate          int
}

type InternalProductPricingFilter struct {
	ID string
}

type InternalProductPricingDeleteFilter struct {
	ID string
}

type InternalProductPrice struct {
	UserCtx common.UserContext

	ID                       string
	InternalProductPricingID string
	CurrencyCode             string
	Amount                   string
	StartedAt                string
	EndedAt                  *string
	Metadata                 map[string]any
	CreatedAt                string
	UpdatedAt                string
}

type InternalProductPriceListFilter struct {
	InternalProductPricingID string
	CurrencyCode             string
	Page                     int
	Paginate                 int
}

type InternalProductPriceDeleteFilter struct {
	ID string
}

type InternalProductPriceOverlapFilter struct {
	InternalProductPricingID string
	CurrencyCode             string
	StartedAt                string
	EndedAt                  *string
	ExcludeID                string
}
