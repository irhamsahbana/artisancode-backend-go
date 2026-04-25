package restentity

import "codebase-app/pkg/types"

type InternalProduct struct {
	ID          string         `json:"id"`
	Code        string         `json:"code"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Status      string         `json:"status"`
	Metadata    map[string]any `json:"metadata"`
	CreatedAt   string         `json:"created_at"`
	UpdatedAt   string         `json:"updated_at"`
}

type GetInternalProductsReq struct {
	Q string `query:"q" validate:"omitempty,min=2"`
	types.MetaQuery
}

func (r *GetInternalProductsReq) SetDefault() {
	r.MetaQuery.SetDefault()
}

type GetInternalProductsResp struct {
	Items []InternalProduct `json:"items"`
	Meta  types.Meta        `json:"meta"`
}

type GetInternalProductReq struct {
	ID string `params:"id" validate:"required"`
}

type GetInternalProductResp struct {
	InternalProduct
}

type CreateInternalProductReq struct {
	Code        string         `json:"code" validate:"required,min=2,max=64"`
	Name        string         `json:"name" validate:"required,min=2,max=255"`
	Description string         `json:"description"`
	Status      string         `json:"status" validate:"required,oneof=draft active inactive archived"`
	Metadata    map[string]any `json:"metadata"`
}

type CreateInternalProductResp struct {
	ID string `json:"id"`
}

type UpdateInternalProductReq struct {
	ID          string         `params:"id" validate:"required"`
	Code        string         `json:"code" validate:"required,min=2,max=64"`
	Name        string         `json:"name" validate:"required,min=2,max=255"`
	Description string         `json:"description"`
	Status      string         `json:"status" validate:"required,oneof=draft active inactive archived"`
	Metadata    map[string]any `json:"metadata"`
}

type DeleteInternalProductReq struct {
	ID string `params:"id" validate:"required"`
}

type InternalProductPricing struct {
	ID                string         `json:"id"`
	InternalProductID string         `json:"internal_product_id"`
	Code              string         `json:"code"`
	Name              string         `json:"name"`
	Description       string         `json:"description"`
	Status            string         `json:"status"`
	Metadata          map[string]any `json:"metadata"`
	CreatedAt         string         `json:"created_at"`
	UpdatedAt         string         `json:"updated_at"`
}

type GetInternalProductPricingsReq struct {
	ID string `params:"id" validate:"required"`
	Q  string `query:"q" validate:"omitempty,min=2"`
	types.MetaQuery
}

func (r *GetInternalProductPricingsReq) SetDefault() {
	r.MetaQuery.SetDefault()
}

type GetInternalProductPricingsResp struct {
	Items []InternalProductPricing `json:"items"`
	Meta  types.Meta               `json:"meta"`
}

type GetInternalProductPricingReq struct {
	PricingID string `params:"pricing_id" validate:"required"`
}

type GetInternalProductPricingResp struct {
	InternalProductPricing
}

type CreateInternalProductPricingReq struct {
	ID          string         `params:"id" validate:"required"`
	Code        string         `json:"code" validate:"required,min=2,max=64"`
	Name        string         `json:"name" validate:"required,min=2,max=255"`
	Description string         `json:"description"`
	Status      string         `json:"status" validate:"required,oneof=draft active inactive archived"`
	Metadata    map[string]any `json:"metadata"`
}

type CreateInternalProductPricingResp struct {
	ID string `json:"id"`
}

type UpdateInternalProductPricingReq struct {
	PricingID   string         `params:"pricing_id" validate:"required"`
	ProductID   string         `json:"internal_product_id" validate:"required"`
	Code        string         `json:"code" validate:"required,min=2,max=64"`
	Name        string         `json:"name" validate:"required,min=2,max=255"`
	Description string         `json:"description"`
	Status      string         `json:"status" validate:"required,oneof=draft active inactive archived"`
	Metadata    map[string]any `json:"metadata"`
}

type DeleteInternalProductPricingReq struct {
	PricingID string `params:"pricing_id" validate:"required"`
}

type InternalProductPrice struct {
	ID                       string         `json:"id"`
	InternalProductPricingID string         `json:"internal_product_pricing_id"`
	CurrencyCode             string         `json:"currency_code"`
	Amount                   string         `json:"amount"`
	StartedAt                string         `json:"started_at"`
	EndedAt                  *string        `json:"ended_at"`
	Metadata                 map[string]any `json:"metadata"`
	CreatedAt                string         `json:"created_at"`
	UpdatedAt                string         `json:"updated_at"`
}

type GetInternalProductPricesReq struct {
	PricingID    string `params:"pricing_id" validate:"required"`
	CurrencyCode string `query:"currency_code" validate:"omitempty,len=3"`
	types.MetaQuery
}

func (r *GetInternalProductPricesReq) SetDefault() {
	r.MetaQuery.SetDefault()
}

type GetInternalProductPricesResp struct {
	Items []InternalProductPrice `json:"items"`
	Meta  types.Meta             `json:"meta"`
}

type CreateInternalProductPriceReq struct {
	PricingID    string         `params:"pricing_id" validate:"required"`
	CurrencyCode string         `json:"currency_code" validate:"required,len=3"`
	Amount       string         `json:"amount" validate:"required"`
	StartedAt    string         `json:"started_at" validate:"required"`
	EndedAt      *string        `json:"ended_at"`
	Metadata     map[string]any `json:"metadata"`
}

type CreateInternalProductPriceResp struct {
	ID string `json:"id"`
}

type UpdateInternalProductPriceReq struct {
	PriceID      string         `params:"price_id" validate:"required"`
	PricingID    string         `json:"internal_product_pricing_id" validate:"required"`
	CurrencyCode string         `json:"currency_code" validate:"required,len=3"`
	Amount       string         `json:"amount" validate:"required"`
	StartedAt    string         `json:"started_at" validate:"required"`
	EndedAt      *string        `json:"ended_at"`
	Metadata     map[string]any `json:"metadata"`
}

type DeleteInternalProductPriceReq struct {
	PriceID string `params:"price_id" validate:"required"`
}
