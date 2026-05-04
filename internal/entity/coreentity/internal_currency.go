package coreentity

import "codebase-app/internal/entity/common"

type InternalCurrency struct {
	UserCtx common.UserContext

	Code          string
	Name          string
	Symbol        string
	DecimalPlaces int
	IsActive      bool
	IsDefault     bool
	SortOrder     int
	Metadata      map[string]any
	CreatedAt     string
	UpdatedAt     string
}

type InternalCurrencyListFilter struct {
	Q        string
	IsActive *bool
	Page     int
	Paginate int
}

type InternalCurrencyFilter struct {
	Code string
}

type InternalCurrencyDeleteFilter struct {
	Code string
}

type InternalPaymentProviderCurrency struct {
	UserCtx common.UserContext

	Provider     string
	CurrencyCode string
	IsActive     bool
	MinAmount    *string
	MaxAmount    *string
	Metadata     map[string]any
	CreatedAt    string
	UpdatedAt    string
}

type InternalPaymentProviderCurrencyListFilter struct {
	Provider string
	Page     int
	Paginate int
}

type InternalPaymentProviderCurrencyFilter struct {
	Provider     string
	CurrencyCode string
}
