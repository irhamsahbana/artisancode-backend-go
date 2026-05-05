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
