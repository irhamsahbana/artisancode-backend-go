package handler

import (
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/entity/restentity"
)

func mapperInternalCurrencyListFilter(req restentity.GetInternalCurrenciesReq) coreentity.InternalCurrencyListFilter {
	return coreentity.InternalCurrencyListFilter{
		Q:        req.Q,
		IsActive: req.IsActive,
		Page:     req.Page,
		Paginate: req.Paginate,
	}
}

func mapperInternalProviderCurrencyListFilter(
	req restentity.GetInternalPaymentProviderCurrenciesReq,
) coreentity.InternalPaymentProviderCurrencyListFilter {
	return coreentity.InternalPaymentProviderCurrencyListFilter{
		Provider: req.Provider,
		Page:     req.Page,
		Paginate: req.Paginate,
	}
}
