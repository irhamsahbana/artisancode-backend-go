package restentity

import "codebase-app/pkg/types"

type InternalCurrency struct {
	Code          string         `json:"code"`
	Name          string         `json:"name"`
	Symbol        string         `json:"symbol"`
	DecimalPlaces int            `json:"decimal_places"`
	IsActive      bool           `json:"is_active"`
	IsDefault     bool           `json:"is_default"`
	SortOrder     int            `json:"sort_order"`
	Metadata      map[string]any `json:"metadata"`
	CreatedAt     string         `json:"created_at"`
	UpdatedAt     string         `json:"updated_at"`
}

type GetInternalCurrenciesReq struct {
	Q        string `query:"q" validate:"omitempty"`
	IsActive *bool  `query:"is_active"`
	types.MetaQuery
}

func (r *GetInternalCurrenciesReq) SetDefault() {
	r.MetaQuery.SetDefault()
}

type GetInternalCurrenciesResp struct {
	Items []InternalCurrency `json:"items"`
	Meta  types.Meta         `json:"meta"`
}

type GetInternalCurrencyReq struct {
	Code string `params:"code" validate:"required,len=3"`
}

type CreateInternalCurrencyReq struct {
	Code          string         `json:"code" validate:"required,len=3"`
	Name          string         `json:"name" validate:"required,min=2,max=64"`
	Symbol        string         `json:"symbol" validate:"max=16"`
	DecimalPlaces int            `json:"decimal_places" validate:"min=0,max=6"`
	IsActive      bool           `json:"is_active"`
	IsDefault     bool           `json:"is_default"`
	SortOrder     int            `json:"sort_order"`
	Metadata      map[string]any `json:"metadata"`
}

type UpdateInternalCurrencyReq struct {
	Code          string         `params:"code" validate:"required,len=3"`
	Name          string         `json:"name" validate:"required,min=2,max=64"`
	Symbol        string         `json:"symbol" validate:"max=16"`
	DecimalPlaces int            `json:"decimal_places" validate:"min=0,max=6"`
	IsActive      bool           `json:"is_active"`
	IsDefault     bool           `json:"is_default"`
	SortOrder     int            `json:"sort_order"`
	Metadata      map[string]any `json:"metadata"`
}

type DeleteInternalCurrencyReq struct {
	Code string `params:"code" validate:"required,len=3"`
}

type InternalPaymentProviderCurrency struct {
	Provider     string         `json:"provider"`
	CurrencyCode string         `json:"currency_code"`
	IsActive     bool           `json:"is_active"`
	MinAmount    *string        `json:"min_amount"`
	MaxAmount    *string        `json:"max_amount"`
	Metadata     map[string]any `json:"metadata"`
	CreatedAt    string         `json:"created_at"`
	UpdatedAt    string         `json:"updated_at"`
}

type GetInternalPaymentProviderCurrenciesReq struct {
	Provider string `params:"provider" validate:"required,min=2,max=32"`
	types.MetaQuery
}

func (r *GetInternalPaymentProviderCurrenciesReq) SetDefault() {
	r.MetaQuery.SetDefault()
}

type GetInternalPaymentProviderCurrenciesResp struct {
	Items []InternalPaymentProviderCurrency `json:"items"`
	Meta  types.Meta                        `json:"meta"`
}

type UpsertInternalPaymentProviderCurrencyReq struct {
	Provider     string         `params:"provider" validate:"required,min=2,max=32"`
	CurrencyCode string         `params:"currency_code" validate:"required,len=3"`
	IsActive     bool           `json:"is_active"`
	MinAmount    *string        `json:"min_amount"`
	MaxAmount    *string        `json:"max_amount"`
	Metadata     map[string]any `json:"metadata"`
}

type DeleteInternalPaymentProviderCurrencyReq struct {
	Provider     string `params:"provider" validate:"required,min=2,max=32"`
	CurrencyCode string `params:"currency_code" validate:"required,len=3"`
}
