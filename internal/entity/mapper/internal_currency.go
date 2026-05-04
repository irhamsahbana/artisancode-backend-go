package mapper

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/entity/restentity"
)

func InternalCurrencyFromCoreToRest(item coreentity.InternalCurrency) restentity.InternalCurrency {
	return restentity.InternalCurrency{
		Code:          item.Code,
		Name:          item.Name,
		Symbol:        item.Symbol,
		DecimalPlaces: item.DecimalPlaces,
		IsActive:      item.IsActive,
		IsDefault:     item.IsDefault,
		SortOrder:     item.SortOrder,
		Metadata:      item.Metadata,
		CreatedAt:     item.CreatedAt,
		UpdatedAt:     item.UpdatedAt,
	}
}

func InternalCurrencyFromRestCreateToCore(
	ctx context.Context,
	req restentity.CreateInternalCurrencyReq,
) coreentity.InternalCurrency {
	return coreentity.InternalCurrency{
		UserCtx:       common.GetUserContext(ctx),
		Code:          req.Code,
		Name:          req.Name,
		Symbol:        req.Symbol,
		DecimalPlaces: req.DecimalPlaces,
		IsActive:      req.IsActive,
		IsDefault:     req.IsDefault,
		SortOrder:     req.SortOrder,
		Metadata:      req.Metadata,
	}
}

func InternalCurrencyFromRestUpdateToCore(
	ctx context.Context,
	req restentity.UpdateInternalCurrencyReq,
) coreentity.InternalCurrency {
	return coreentity.InternalCurrency{
		UserCtx:       common.GetUserContext(ctx),
		Code:          req.Code,
		Name:          req.Name,
		Symbol:        req.Symbol,
		DecimalPlaces: req.DecimalPlaces,
		IsActive:      req.IsActive,
		IsDefault:     req.IsDefault,
		SortOrder:     req.SortOrder,
		Metadata:      req.Metadata,
	}
}

func InternalPaymentProviderCurrencyFromCoreToRest(
	item coreentity.InternalPaymentProviderCurrency,
) restentity.InternalPaymentProviderCurrency {
	return restentity.InternalPaymentProviderCurrency{
		Provider:     item.Provider,
		CurrencyCode: item.CurrencyCode,
		IsActive:     item.IsActive,
		MinAmount:    item.MinAmount,
		MaxAmount:    item.MaxAmount,
		Metadata:     item.Metadata,
		CreatedAt:    item.CreatedAt,
		UpdatedAt:    item.UpdatedAt,
	}
}

func InternalPaymentProviderCurrencyFromRestUpsertToCore(
	ctx context.Context,
	req restentity.UpsertInternalPaymentProviderCurrencyReq,
) coreentity.InternalPaymentProviderCurrency {
	return coreentity.InternalPaymentProviderCurrency{
		UserCtx:      common.GetUserContext(ctx),
		Provider:     req.Provider,
		CurrencyCode: req.CurrencyCode,
		IsActive:     req.IsActive,
		MinAmount:    req.MinAmount,
		MaxAmount:    req.MaxAmount,
		Metadata:     req.Metadata,
	}
}
