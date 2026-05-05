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
