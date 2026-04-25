package mapper

import (
	"context"

	"codebase-app/internal/entity/common"
	"codebase-app/internal/entity/coreentity"
	"codebase-app/internal/entity/restentity"

	"github.com/shopspring/decimal"
)

func InternalProductFromCoreToRest(item coreentity.InternalProduct) restentity.InternalProduct {
	return restentity.InternalProduct{
		ID:          item.ID,
		Code:        item.Code,
		Name:        item.Name,
		Description: item.Description,
		Status:      item.Status,
		Metadata:    item.Metadata,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}
}

func InternalProductFromRestCreateToCore(ctx context.Context, req restentity.CreateInternalProductReq) coreentity.InternalProduct {
	uc := common.GetUserContext(ctx)
	return coreentity.InternalProduct{
		UserCtx:     uc,
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		Status:      req.Status,
		Metadata:    req.Metadata,
	}
}

func InternalProductFromRestUpdateToCore(ctx context.Context, req restentity.UpdateInternalProductReq) coreentity.InternalProduct {
	uc := common.GetUserContext(ctx)
	return coreentity.InternalProduct{
		UserCtx:     uc,
		ID:          req.ID,
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		Status:      req.Status,
		Metadata:    req.Metadata,
	}
}

func InternalProductPricingFromCoreToRest(item coreentity.InternalProductPricing) restentity.InternalProductPricing {
	return restentity.InternalProductPricing{
		ID:                item.ID,
		InternalProductID: item.InternalProductID,
		Code:              item.Code,
		Name:              item.Name,
		Description:       item.Description,
		Status:            item.Status,
		Metadata:          item.Metadata,
		CreatedAt:         item.CreatedAt,
		UpdatedAt:         item.UpdatedAt,
	}
}

func InternalProductPricingFromRestCreateToCore(ctx context.Context, req restentity.CreateInternalProductPricingReq) coreentity.InternalProductPricing {
	uc := common.GetUserContext(ctx)
	return coreentity.InternalProductPricing{
		UserCtx:           uc,
		InternalProductID: req.ID,
		Code:              req.Code,
		Name:              req.Name,
		Description:       req.Description,
		Status:            req.Status,
		Metadata:          req.Metadata,
	}
}

func InternalProductPricingFromRestUpdateToCore(ctx context.Context, req restentity.UpdateInternalProductPricingReq) coreentity.InternalProductPricing {
	uc := common.GetUserContext(ctx)
	return coreentity.InternalProductPricing{
		UserCtx:           uc,
		ID:                req.PricingID,
		InternalProductID: req.ProductID,
		Code:              req.Code,
		Name:              req.Name,
		Description:       req.Description,
		Status:            req.Status,
		Metadata:          req.Metadata,
	}
}

func InternalProductPriceFromCoreToRest(item coreentity.InternalProductPrice) restentity.InternalProductPrice {
	return restentity.InternalProductPrice{
		ID:                       item.ID,
		InternalProductPricingID: item.InternalProductPricingID,
		CurrencyCode:             item.CurrencyCode,
		Amount:                   item.Amount.String(),
		StartedAt:                item.StartedAt,
		EndedAt:                  item.EndedAt,
		Metadata:                 item.Metadata,
		CreatedAt:                item.CreatedAt,
		UpdatedAt:                item.UpdatedAt,
	}
}

func InternalProductPriceFromRestCreateToCore(ctx context.Context, req restentity.CreateInternalProductPriceReq) coreentity.InternalProductPrice {
	uc := common.GetUserContext(ctx)
	amount, _ := decimal.NewFromString(req.Amount)
	return coreentity.InternalProductPrice{
		UserCtx:                  uc,
		InternalProductPricingID: req.PricingID,
		CurrencyCode:             req.CurrencyCode,
		Amount:                   amount,
		StartedAt:                req.StartedAt,
		EndedAt:                  req.EndedAt,
		Metadata:                 req.Metadata,
	}
}

func InternalProductPriceFromRestUpdateToCore(ctx context.Context, req restentity.UpdateInternalProductPriceReq) coreentity.InternalProductPrice {
	uc := common.GetUserContext(ctx)
	amount, _ := decimal.NewFromString(req.Amount)
	return coreentity.InternalProductPrice{
		UserCtx:                  uc,
		ID:                       req.PriceID,
		InternalProductPricingID: req.PricingID,
		CurrencyCode:             req.CurrencyCode,
		Amount:                   amount,
		StartedAt:                req.StartedAt,
		EndedAt:                  req.EndedAt,
		Metadata:                 req.Metadata,
	}
}
