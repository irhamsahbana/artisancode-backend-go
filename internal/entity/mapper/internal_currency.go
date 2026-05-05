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
	}
}
