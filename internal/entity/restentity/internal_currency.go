package restentity

import "codebase-app/pkg/types"

type InternalCurrency struct {
	Code          string `json:"code"`
	Name          string `json:"name"`
	Symbol        string `json:"symbol"`
	DecimalPlaces int    `json:"decimal_places"`
	IsActive      bool   `json:"is_active"`
	IsDefault     bool   `json:"is_default"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
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
	Code          string `json:"code" validate:"required,len=3"`
	Name          string `json:"name" validate:"required,min=2,max=64"`
	Symbol        string `json:"symbol" validate:"max=16"`
	DecimalPlaces int    `json:"decimal_places" validate:"min=0,max=6"`
	IsActive      bool   `json:"is_active"`
	IsDefault     bool   `json:"is_default"`
}

type UpdateInternalCurrencyReq struct {
	Code          string `params:"code" validate:"required,len=3"`
	Name          string `json:"name" validate:"required,min=2,max=64"`
	Symbol        string `json:"symbol" validate:"max=16"`
	DecimalPlaces int    `json:"decimal_places" validate:"min=0,max=6"`
	IsActive      bool   `json:"is_active"`
	IsDefault     bool   `json:"is_default"`
}

type DeleteInternalCurrencyReq struct {
	Code string `params:"code" validate:"required,len=3"`
}
